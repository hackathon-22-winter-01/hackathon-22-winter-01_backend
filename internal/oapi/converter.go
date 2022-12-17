package oapi

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"go.uber.org/zap"
)

func PlayerFromDomain(dp *domain.Player) Player {
	rails := []int{0}

	if l := len(dp.BranchEvents); l > 0 {
		rails = make([]int, len(dp.BranchEvents[l-1].AfterRails))

		for i, r := range dp.BranchEvents[l-1].AfterRails {
			rails[i] = r.Index
		}
	}

	return Player{
		Id:       dp.ID,
		Life:     domain.CalculateLife(dp.LifeEvents),
		MainRail: consts.RailLimit / 2,
		Rails:    rails,
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

func RoomFromDomain(dr *domain.Room) Room {
	var PLayers = make([]Player, len(dr.Players))

	for i, p := range dr.Players {
		PLayers[i] = PlayerFromDomain(p)
	}

	return Room{
		Id:        dr.ID,
		Players:   PLayers,
		StartedAt: dr.StartedAt,
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
