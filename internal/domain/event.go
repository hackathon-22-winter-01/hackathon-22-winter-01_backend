package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
)

// commonEvent 共通のイベント情報
type commonEvent struct {
	ID        uuid.UUID
	CardType  CardType
	CreatedAt time.Time
}

type JustCardEvent struct {
	commonEvent
}

func NewJustCardEvent(id uuid.UUID, cardType CardType, createdAt time.Time) *JustCardEvent {
	return &JustCardEvent{
		commonEvent{
			ID:        id,
			CardType:  cardType,
			CreatedAt: createdAt,
		},
	}
}

// BlockEvent 妨害イベントの情報
type BlockEvent struct {
	commonEvent
	Type         BlockEventType
	AttackerID   uuid.UUID
	TargetID     uuid.UUID
	TargetRailID uuid.UUID
}

// BlockEventType 妨害イベントの種類
type BlockEventType uint8

const (
	BlockEventTypeCreated BlockEventType = iota
	BlockEventTypeCanceled
)

func NewBlockEvent(id uuid.UUID, cardType CardType, createdAt time.Time, typ BlockEventType, attackerID, targetID, targetRailID uuid.UUID) *BlockEvent {
	return &BlockEvent{
		commonEvent: commonEvent{
			ID:        id,
			CardType:  cardType,
			CreatedAt: createdAt,
		},
		Type:         typ,
		AttackerID:   attackerID,
		TargetID:     targetID,
		TargetRailID: targetRailID,
	}
}

// LifeEvent プレイヤーのライフに関する情報
type LifeEvent struct {
	commonEvent
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

func NewLifeEvent(id uuid.UUID, cardType CardType, createdAt time.Time, typ LifeEventType, diff float32) *LifeEvent {
	return &LifeEvent{
		commonEvent: commonEvent{
			ID:        id,
			CardType:  cardType,
			CreatedAt: createdAt,
		},
		Type: typ,
		Diff: diff,
	}
}

func CalculateLife(events []*LifeEvent) float32 {
	life := consts.MaxLife

	for _, e := range events {
		switch e.Type {
		case LifeEventTypeDamaged:
			life -= e.Diff
		case LifeEventTypeHealed:
			life += e.Diff
		}
	}

	return life
}

// BranchEvent プレイヤーのレールの分岐/結合イベントの情報
type BranchEvent struct {
	commonEvent
	Type       BranchEventType
	AttackerID uuid.UUID
	TargetID   uuid.UUID
	AfterRails []*Rail
}

// BranchEventType イベントの種類
type BranchEventType uint8

const (
	BranchEventBlockCreated BranchEventType = iota
	BranchEventCardReset
	BranchEventCardUsed
	BranchEventLifeChanged
	BranchEventCreated
	BranchEventMerged
)

func NewBranchEvent(id uuid.UUID, cardType CardType, createdAt time.Time, typ BranchEventType, attackerID, targetID uuid.UUID, afterRails []*Rail) *BranchEvent {
	return &BranchEvent{
		commonEvent: commonEvent{
			ID:        id,
			CardType:  cardType,
			CreatedAt: createdAt,
		},
		Type:       typ,
		AttackerID: attackerID,
		TargetID:   targetID,
		AfterRails: afterRails,
	}
}
