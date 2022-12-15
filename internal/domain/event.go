package domain

import (
	"time"

	"github.com/google/uuid"
)

// CardEvent プレイヤーがカードを使用して発生するイベントの情報
type CardEvent struct {
	ID          uuid.UUID
	Type        CardEventType
	Attacker    *Player
	Target      *Player
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
