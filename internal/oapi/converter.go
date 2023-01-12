package oapi

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"go.uber.org/zap"
)

func PlayerFromDomain(dp *domain.Player) Player {
	return Player{
		Id:   dp.ID,
		Name: dp.Name,
		Life: domain.CalculateLife(dp.LifeEvents),
	}
}

func CardFromDomain(dc *domain.Card) Card {
	var typ CardType

	switch dc.Type {
	case domain.CardTypePullShark:
		typ = CardTypePullShark
	case domain.CardTypePairExtraordinaire:
		typ = CardTypePairExtraordinaire
	default:
		typ = CardTypePullShark

		log.L().Error("unknown card type", zap.String("type", string(dc.Type)))
	}

	return Card{
		Id:   dc.ID,
		Type: typ,
	}
}

func (t CardType) ToDomain() domain.CardType {
	m := map[CardType]domain.CardType{
		CardTypeYolo:               domain.CardTypeYolo,
		CardTypeGalaxyBrain:        domain.CardTypeGalaxyBrain,
		CardTypeOpenSourcerer:      domain.CardTypeOpenSourcerer,
		CardTypeRefactoring:        domain.CardTypeRefactoring,
		CardTypePairExtraordinaire: domain.CardTypePairExtraordinaire,
		CardTypeLgtm:               domain.CardTypeLgtm,
		CardTypePullShark:          domain.CardTypePullShark,
		CardTypeStarstruck:         domain.CardTypeStarstruck,
	}

	if t, ok := m[t]; ok {
		return t
	}

	return domain.CardTypeNone
}
