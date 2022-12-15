package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleCardEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyCardEvent()
	if err != nil {
		return err
	}

	room, err := h.roomRepo.FindRoom(repository.CommonRoomID)
	if err != nil {
		return err
	}

	target, ok := room.FindPlayer(b.TargetId)
	if !ok {
		return errors.New("player not found")
	}

	var (
		newRail     = domain.NewRail()
		beforeRails = []*domain.Rail{newRail}
		afterRails  = []*domain.Rail{}
		res         *oapi.WsResponse
	)

	switch b.Type {
	case oapi.CardTypeCreateRail:
		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = append(beforeRails, newRail)
		}

		res, err = oapi.NewWsResponseRailCreated(jst.Now(), uuid.New(), target.Main.ID, h.playerID, b.TargetId)
		if err != nil {
			return err
		}

	case oapi.CardTypeCreateBlock:
		res, err = oapi.NewWsResponseBlockCreated(jst.Now(), h.playerID, b.TargetId)
		if err != nil {
			return err
		}

	default:
		return errors.New("invalid card type")
	}

	h.sender.Bloadcast(room.ID, res)

	target.Events = append(target.Events, domain.NewCardEvent(
		uuid.New(),
		domain.RailCreated,
		h.playerID,
		target.ID,
		beforeRails,
		afterRails,
	))

	return nil
}
