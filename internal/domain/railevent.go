package domain

import (
	"github.com/google/uuid"
)

// RailEvent プレイヤーのレール関するイベントの情報
type RailEvent struct {
	ID          uuid.UUID
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
		ID:          id,
		Type:        typ,
		AttackerID:  attackerID,
		TargetID:    targetID,
		BeforeRails: beforeRails,
		AfterRails:  afterRails,
	}
}
