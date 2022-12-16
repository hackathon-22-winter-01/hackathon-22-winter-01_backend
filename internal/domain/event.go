package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

// commonEvent 共通のイベント情報
type commonEvent struct {
	ID        uuid.UUID
	CardType  CardType
	CreatedAt time.Time
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

func NewLifeEvent(id uuid.UUID, typ LifeEventType, diff float32, createdAt time.Time) *LifeEvent {
	return &LifeEvent{
		commonEvent: commonEvent{
			ID: id,
			// TODO: CardType
			CreatedAt: jst.Now(),
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

// RailEvent プレイヤーのレール関するイベントの情報
type RailEvent struct {
	commonEvent
	Type        RailEventType
	AttackerID  uuid.UUID
	TargetID    uuid.UUID
	BeforeRails []*Rail
	AfterRails  []*Rail
}

// RailEventType イベントの種類
type RailEventType uint8

const (
	RailEventBlockCreated RailEventType = iota
	RailEventCardReset
	RailEventCardUsed
	RailEventLifeChanged
	RailEventCreated
	RailEventMerged
)

func NewRailEvent(id uuid.UUID, typ RailEventType, attackerID, targetID uuid.UUID, beforeRails, afterRails []*Rail) *RailEvent {
	return &RailEvent{
		commonEvent: commonEvent{
			ID: id,
			// TODO: CardType
			CreatedAt: jst.Now(),
		},
		Type:        typ,
		AttackerID:  attackerID,
		TargetID:    targetID,
		BeforeRails: beforeRails,
		AfterRails:  afterRails,
	}
}
