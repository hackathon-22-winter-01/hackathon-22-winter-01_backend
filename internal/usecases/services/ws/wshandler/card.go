package wshandler

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleCardEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyCardEvent()
	if err != nil {
		return err
	}

	target, ok := h.room.FindPlayer(b.TargetId)
	if !ok {
		return errors.New("player not found")
	}

	var (
		beforeRails = []*domain.Rail{target.Main}
		afterRails  = []*domain.Rail{target.Main}
		res         *oapi.WsResponse
	)

	switch b.Type {
	case oapi.CardTypeYolo:
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
		}

		var (
			parentID = target.Main.ID
			childID  uuid.UUID
		)

		// シャッフルする
		rails := []*domain.Rail{}
		copy(rails, beforeRails)

		rand.Shuffle(len(rails), func(i, j int) { rails[i], rails[j] = rails[j], rails[i] })

		// 一番最後のブロックの親を探す
		for _, rail := range rails {
			if rail.ID == parentID {
				continue
			}

			if !rail.HasBlock {
				childID = rail.ID
			}
		}

		if childID != uuid.Nil && parentID != uuid.Nil {
			// 親が見つかったら、親の子を消す
			for _, rail := range rails {
				if rail.ID != target.Main.ID {
					afterRails = append(afterRails, rail)
				}
			}

			// マージしたことを通知
			res, err = oapi.NewWsResponseRailMerged(jst.Now(), childID, parentID, b.TargetId)
			if err != nil {
				return err
			}
		} else {
			// 親が見つからなかったら、何もしない
			afterRails = beforeRails
			res = oapi.WsResponseFromType(oapi.WsResponseTypeNoop, jst.Now())
		}

	case oapi.CardTypeGalaxyBrain:
		return nil

	case oapi.CardTypeOpenSourcerer:
		return nil

	case oapi.CardTypeRefactoring:
		return nil

	case oapi.CardTypePairExtraordinaire:
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = lastEvent.AfterRails
		}

		res, err = oapi.NewWsResponseBlockCreated(jst.Now(), h.playerID, b.TargetId, 2, 30)
		if err != nil {
			return err
		}

	case oapi.CardTypeLgtm:
		return nil

	case oapi.CardTypePullShark:
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = append(beforeRails, domain.NewRail())
		}

		res, err = oapi.NewWsResponseRailCreated(jst.Now(), uuid.New(), target.Main.ID, h.playerID, b.TargetId)
		if err != nil {
			return err
		}

	case oapi.CardTypeStarstruck:
		return nil

	default:
		return errors.New("invalid card type")
	}

	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	target.Events = append(target.Events, domain.NewRailEvent(
		uuid.New(),
		domain.RailEventCreated,
		h.playerID,
		target.ID,
		beforeRails,
		afterRails,
	))

	return nil
}
