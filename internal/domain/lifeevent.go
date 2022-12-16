package domain

import (
	"time"

	"github.com/google/uuid"
)

// LifeEvent プレイヤーのライフに関する情報
type LifeEvent struct {
	ID        uuid.UUID
	Type      LifeEventType
	Diff      float32
	CreatedAt time.Time
}

// LifeEventType ライフイベントの種類
type LifeEventType uint8

const (
	LifeEventTypeDamaged LifeEventType = iota
	LifeEventTypeHealed
)

func NewLifeEvent(id uuid.UUID, typ LifeEventType, diff float32, createdAt time.Time) *LifeEvent {
	return &LifeEvent{
		ID:        id,
		Type:      typ,
		Diff:      diff,
		CreatedAt: createdAt,
	}
}
