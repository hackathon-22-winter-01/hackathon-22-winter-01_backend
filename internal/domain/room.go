package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/sync"
)

// Room 対戦部屋の情報
type Room struct {
	ID        uuid.UUID
	Players   sync.Slice[*Player]
	StartedAt time.Time
}

// NewRoom 新しい対戦部屋を作成する
func NewRoom(id uuid.UUID) *Room {
	return &Room{
		ID:        id,
		Players:   sync.NewSlice[*Player](),
		StartedAt: time.Time{},
	}
}

func (r *Room) FindPlayer(playerID uuid.UUID) (*Player, bool) {
	for _, p := range r.Players.Clone() {
		if p.ID == playerID {
			return p, true
		}
	}

	return nil, false
}
