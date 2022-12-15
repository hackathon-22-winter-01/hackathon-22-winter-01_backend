package oapi

import "github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"

func PlayerFromDomain(dp *domain.Player) Player {
	p := Player{
		PlayerId: dp.ID,
		Life:     3,
	}

	for _, le := range dp.LifeEvents {
		if le.Type == domain.LifeEventDecrement {
			p.Life--
		}
	}

	return p
}
