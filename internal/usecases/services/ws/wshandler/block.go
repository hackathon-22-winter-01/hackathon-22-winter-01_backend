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
			b.CardType.ToDomain(),
			jst.Now(),
			domain.BlockEventTypeCanceled,
			target.ID,
			b.RailIndex,
		))

		res, err = oapi.NewWsResponseBlockCanceled(now, h.playerID, b.RailIndex, b.CardType)
		if err != nil {
			return err
		}

		if err := h.sender.Broadcast(h.room.ID, res); err != nil {
			return err
		}

	case oapi.BlockEventTypeCrashed:
		target.BlockEvents = append(target.BlockEvents, domain.NewBlockEvent(
			uuid.New(),
			b.CardType.ToDomain(),
			jst.Now(),
			domain.BlockEventTypeCrashed,
			target.ID,
			b.RailIndex,
		))

		cardType := getCardTypeFromRailID(target, b.RailIndex)

		_, attack, err := cardType.DelayAndAttack()
		if err != nil {
			return err
		}

		target.LifeEvents = append(target.LifeEvents, domain.NewLifeEvent(
			uuid.New(),
			b.CardType.ToDomain(),
			jst.Now(),
			domain.LifeEventTypeDamaged,
			attack,
		))

		res, err = oapi.NewWsResponseBlockCrashed(
			now,
			target.ID,
			b.RailIndex,
			b.CardType,
		)
		if err != nil {
			return err
		}

		if err := h.sender.Broadcast(h.room.ID, res); err != nil {
			return err
		}

		res, err = oapi.NewWsResponseLifeChanged(
			now,
			target.ID,
			b.CardType,
			domain.CalculateLife(target.LifeEvents),
		)
		if err != nil {
			return err
		}

		if err := h.sender.Broadcast(h.room.ID, res); err != nil {
			return err
		}

	default:
		return errors.New("invalid block type")
	}

	return nil
}

func getCardTypeFromRailID(p *domain.Player, railIndex int) domain.CardType {
	res := domain.CardTypeNone

	for _, e := range p.BlockEvents {
		if e.TargetRailIndex == railIndex {
			res = e.CardType
		}
	}

	return res
}
