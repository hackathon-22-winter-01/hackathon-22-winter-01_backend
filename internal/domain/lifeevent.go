package domain

import (
	"time"

	"github.com/google/uuid"
)

// LifeEvent プレイヤーのライフに関する情報
type LifeEvent struct {
	ID        uuid.UUID
	Type      LifeEventType
	CreatedAt time.Time
}

// LifeEventType ライフイベントの種類
type LifeEventType uint8

const (
	LifeEventDecrement LifeEventType = iota
)

func NewLifeEvent(id uuid.UUID, typ LifeEventType, createdAt time.Time) *LifeEvent {
	return &LifeEvent{
		ID:        id,
		Type:      typ,
		CreatedAt: createdAt,
	}
}
