package repoimpl

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/sync"
)

type playerRepository struct {
	playerMap sync.Map[uuid.UUID, *domain.Player]
}

func NewPlayerRepository() *playerRepository {
	return &playerRepository{
		playerMap: sync.Map[uuid.UUID, *domain.Player]{},
	}
}

func (r *playerRepository) FindPlayer(playerID uuid.UUID) (*domain.Player, error) {
	player, ok := r.playerMap.Load(playerID)
	if !ok {
		return nil, errors.New("プレイヤーが存在しません")
	}

	return player, nil
}

func (r *playerRepository) CreatePlayer(name string) (*domain.Player, error) {
	playerID := uuid.New()

	player := domain.NewPlayer(playerID, name)

	r.playerMap.Store(playerID, player)

	return player, nil
}
