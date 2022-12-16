package wshandler

import (
	"errors"

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
		return nil

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
