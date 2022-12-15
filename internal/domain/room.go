package domain

import (
	"time"

	"github.com/google/uuid"
)

// Room 対戦部屋の情報
type Room struct {
	ID        uuid.UUID
	Players   []*Player
	StartedAt time.Time
}

// NewRoom 新しい対戦部屋を作成する
func NewRoom(id uuid.UUID) *Room {
	return &Room{
		ID:        id,
		Players:   make([]*Player, 0),
		StartedAt: time.Time{},
	}
}

func (r *Room) FindPlayer(playerID uuid.UUID) (*Player, bool) {
	for _, p := range r.Players {
		if p.ID == playerID {
			return p, true
		}
	}

	return nil, false
}
