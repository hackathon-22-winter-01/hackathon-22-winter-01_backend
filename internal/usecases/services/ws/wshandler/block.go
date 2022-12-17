package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleBlockEvent(reqbody oapi.WsRequest_Body) error {
	b, err := reqbody.AsWsRequestBodyBlockEvent()
	if err != nil {
		return err
	}

	target, ok := h.room.FindPlayer(h.playerID)
	if !ok {
		return errors.New("player not found")
	}

	var (
		res *oapi.WsResponse
		now = jst.Now()
	)

	switch b.Type {
	case oapi.BlockEventTypeCanceled:
		target.BlockEvents = append(target.BlockEvents, domain.NewBlockEvent(
			uuid.New(),
			domain.CardTypeNone,
			jst.Now(),
			domain.BlockEventTypeCanceled,
			target.ID,
			target.ID,
			b.RailId,
		))

		res, err = oapi.NewWsResponseBlockCanceled(now, b.RailId)
		if err != nil {
			return err
		}

	case oapi.BlockEventTypeCrashed:
		target.BlockEvents = append(target.BlockEvents, domain.NewBlockEvent(
			uuid.New(),
			domain.CardTypeNone,
			jst.Now(),
			domain.BlockEventTypeCrashed,
			target.ID,
			target.ID,
			b.RailId,
		))

		cardType := getCardTypeFromRailID(target, b.RailId)
		_, attack, _ := cardType.DelayAndAttack()

		target.LifeEvents = append(target.LifeEvents, domain.NewLifeEvent(
			uuid.New(),
			domain.CardTypeNone,
			jst.Now(),
			domain.LifeEventTypeDamaged,
			attack,
		))

		res, err = oapi.NewWsResponseBlockCrashed(now, domain.CalculateLife(target.LifeEvents), target.ID, b.RailId)
		if err != nil {
			return err
		}

	default:
		return errors.New("invalid block type")
	}

	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	return nil
}

func getCardTypeFromRailID(p *domain.Player, railID uuid.UUID) domain.CardType {
	res := domain.CardTypeNone

	for _, e := range p.BlockEvents {
		if e.TargetRailID == railID {
			res = e.CardType
		}
	}

	return res
}
