package domain

import (
	"time"

	"github.com/google/uuid"
)

// CardEvent プレイヤーがカードを使用して発生するイベントの情報
type CardEvent struct {
	ID          uuid.UUID
	Type        CardEventType
	AttackerID  uuid.UUID
	TargetID    uuid.UUID
	BeforeRails []*Rail
	AfterRails  []*Rail
}

// CardEventType イベントの種類
type CardEventType uint8

const (
	BlockCreated CardEventType = iota
	CardReset
	CardUsed
	LifeChanged
	RailCreated
	RailMerged
)

func NewCardEvent(id uuid.UUID, typ CardEventType, attackerID, targetID uuid.UUID, beforeRails, afterRails []*Rail) *CardEvent {
	return &CardEvent{
		ID:          id,
		Type:        typ,
		AttackerID:  attackerID,
		TargetID:    targetID,
		BeforeRails: beforeRails,
		AfterRails:  afterRails,
	}
}

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
