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
