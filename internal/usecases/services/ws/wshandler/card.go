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
			childID uuid.UUID
		)

		rails := []*domain.Rail{}
		copy(rails, beforeRails)

		rand.Shuffle(len(rails), func(i, j int) { rails[i], rails[j] = rails[j], rails[i] })

		for _, rail := range rails {
			if rail.ID != target.Main.ID && !rail.HasBlock {
				childID = rail.ID
				break
			}
		}

		if childID != uuid.Nil {
			for _, rail := range rails {
				if rail.ID != childID {
					afterRails = append(afterRails, rail)
				}
			}

			res, err = oapi.NewWsResponseRailMerged(jst.Now(), childID, target.Main.ID, b.TargetId)
			if err != nil {
				return err
			}
		} else {
			afterRails = beforeRails
			res = oapi.WsResponseFromType(oapi.WsResponseTypeNoop, jst.Now())
		}

	case oapi.CardTypeGalaxyBrain:
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = lastEvent.AfterRails
		}

		res = oapi.WsResponseFromType(oapi.WsResponseTypeNoop, jst.Now())
		return nil

	case oapi.CardTypeOpenSourcerer:
		return nil

	case oapi.CardTypeRefactoring:
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = lastEvent.AfterRails
		}

		if h.playerID != b.TargetId {
			return errors.New("targetID is different from playerID")
		}

		res, err = oapi.NewWsResponseBlockCreated(jst.Now(), h.playerID, b.TargetId, 1, 5)
		if err != nil {
			return err
		}

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
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = lastEvent.AfterRails
		}

		res, err = oapi.NewWsResponseBlockCreated(jst.Now(), h.playerID, b.TargetId, 3, 20)
		if err != nil {
			return err
		}

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
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = lastEvent.AfterRails
		}

		res, err = oapi.NewWsResponseBlockCreated(jst.Now(), h.playerID, b.TargetId, 5, 50)
		if err != nil {
			return err
		}

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
