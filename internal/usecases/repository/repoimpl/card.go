package repoimpl

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/config"
)

type cardRepository struct{}

func NewCardRepository() repository.CardRepository {
	return &cardRepository{}
}

func (r *cardRepository) DrawCards(roomID uuid.UUID, playerID uuid.UUID, num int) ([]*domain.Card, error) {
	// テスト時は同じカードを返すようにする
	// TODO: 将来的にはモックを立てる
	if !*config.Production {
		allCards := make([]*domain.Card, 0, num)

		for i := 0; i < num; i++ {
			var typ domain.CardType

			switch i % 2 {
			case 0:
				typ = domain.CardTypeCreateRail
			case 1:
				typ = domain.CardTypeCreateBlock
			}

			allCards = append(allCards, domain.NewCard(uuid.New(), typ))
		}

		return allCards, nil
	}

	// TODO: ロジックを実装する
	panic("not implemented")
}
