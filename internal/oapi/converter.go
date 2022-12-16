package oapi

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"go.uber.org/zap"
)

func PlayerFromDomain(dp *domain.Player) Player {
	var (
		life  = consts.MaxLife
		rails = []Rail{
			{Id: dp.Main.ID},
		}
	)

	for _, le := range dp.LifeEvents {
		if le.Type == domain.LifeEventDecrement {
			life--
		}
	}

	if eventLen := len(dp.Events); eventLen > 0 {
		rails = make([]Rail, len(dp.Events[eventLen-1].AfterRails))

		for i, r := range dp.Events[eventLen-1].AfterRails {
			rails[i] = Rail{Id: r.ID}
		}
	}

	return Player{
		PlayerId: dp.ID,
		Life:     life,
		MainRail: Rail{Id: dp.Main.ID},
		Rails:    rails,
	}
}

func CardFromDomain(dc *domain.Card) Card {
	var typ CardType

	switch dc.Type {
	case domain.CardTypeCreateRail:
		typ = CardTypeCreateRail
	case domain.CardTypeCreateBlock:
		typ = CardTypeCreateBlock
	default:
		typ = CardTypeCreateRail

		log.L().Error("unknown card type", zap.String("type", string(dc.Type)))
	}

	return Card{
		Id:   dc.ID,
		Type: typ,
	}
}
