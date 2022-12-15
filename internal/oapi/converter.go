package oapi

import "github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"

func PlayerFromDomain(dp *domain.Player) Player {
	var (
		life  = 3
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
