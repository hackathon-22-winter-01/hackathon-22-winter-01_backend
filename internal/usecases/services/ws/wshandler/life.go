package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleLifeEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyLifeEvent()
	if err != nil {
		return err
	}

	target, ok := h.room.FindPlayer(h.playerID)
	if !ok {
		return errors.New("player not found")
	}

	var typ domain.LifeEventType

	switch b.Type {
	case oapi.LifeEventTypeDamaged:
		typ = domain.LifeEventTypeDamaged
	case oapi.LifeEventTypeHealed:
		typ = domain.LifeEventTypeHealed
	default:
		return errors.New("invalid life type")
	}

	now := jst.Now()
	target.LifeEvents = append(target.LifeEvents, domain.NewLifeEvent(
		uuid.New(),
		domain.CardTypeNone,
		now,
		typ,
		b.Diff,
	))

	life := domain.CalculateLife(target.LifeEvents)

	// TODO: 実装
	// if life <= 0 {
	// }

	res, err := oapi.NewWsResponseLifeChanged(now, h.playerID, life)
	if err != nil {
		return err
	}

	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	return nil
}
