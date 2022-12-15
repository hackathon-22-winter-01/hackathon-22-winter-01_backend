package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleLifeEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyLifeEvent()
	if err != nil {
		return err
	}

	room, err := h.roomRepo.FindRoom(repository.CommonRoomID)
	if err != nil {
		return err
	}

	target, ok := room.FindPlayer(h.playerID)
	if !ok {
		return errors.New("player not found")
	}

	now := jst.Now()
	target.LifeEvents = append(target.LifeEvents, domain.NewLifeEvent(
		uuid.New(),
		domain.LifeEventDecrement,
		now,
	))

	switch b.Type {
	case oapi.LifeEventTypeDecrement:
		life := consts.MaxLife

		for _, e := range target.LifeEvents {
			if e.Type == domain.LifeEventDecrement {
				life--
			}
		}

		// TODO: ライフが0になったらゲームオーバー

		res, err := oapi.NewWsResponseLifeChanged(now, h.playerID, life)
		if err != nil {
			return err
		}

		h.sender.Bloadcast(room.ID, res)

	default:
		return errors.New("invalid life type")
	}

	return nil
}
