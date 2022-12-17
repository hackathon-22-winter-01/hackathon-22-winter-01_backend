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
	Type            BlockEventType
	AttackerID      uuid.UUID
	TargetRailIndex int
}

// BlockEventType 妨害イベントの種類
type BlockEventType uint8

const (
	BlockEventTypeCreated BlockEventType = iota
	BlockEventTypeCanceled
	BlockEventTypeCrashed
)

func NewBlockEvent(id uuid.UUID, cardType CardType, createdAt time.Time, typ BlockEventType, attackerID uuid.UUID, targetRailIndex int) *BlockEvent {
	return &BlockEvent{
		commonEvent: commonEvent{
			ID:        id,
			CardType:  cardType,
			CreatedAt: createdAt,
		},
		Type:            typ,
		AttackerID:      attackerID,
		TargetRailIndex: targetRailIndex,
	}
}

// CalcBlockedRails ブロックされているレールのインデックスを取得する
// blockedEvents[i] == true ならば i 番目のレールがブロックされている
func CalcBlockedRails(events []*BlockEvent) [consts.RailLimit]bool {
	blockedRails := [consts.RailLimit]bool{}

	for _, e := range events {
		switch e.Type {
		case BlockEventTypeCreated:
			blockedRails[e.TargetRailIndex] = true
		case BlockEventTypeCanceled:
			blockedRails[e.TargetRailIndex] = false
		case BlockEventTypeCrashed:
			blockedRails[e.TargetRailIndex] = false
		}
	}

	return blockedRails
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
	Type            BranchEventType
	AttackerID      uuid.UUID
	ChildRailIndex  int
	ParentRailIndex int
}

// BranchEventType イベントの種類
type BranchEventType uint8

const (
	BranchEventCreated BranchEventType = iota
	BranchEventMerged
)

func NewBranchEvent(id uuid.UUID, cardType CardType, createdAt time.Time, typ BranchEventType, attackerID uuid.UUID, childRailIndex int, parentRailIndex int) *BranchEvent {
	return &BranchEvent{
		commonEvent: commonEvent{
			ID:        id,
			CardType:  cardType,
			CreatedAt: createdAt,
		},
		Type:            typ,
		AttackerID:      attackerID,
		ChildRailIndex:  childRailIndex,
		ParentRailIndex: parentRailIndex,
	}
}

// CalcUsedRails 使用中のレールのインデックスを取得する
// unusedRails[i] == true ならばレールiは使用中
func CalcUsedRails(events []*BranchEvent) [consts.RailLimit]bool {
	usedRails := [consts.RailLimit]bool{}
	usedRails[consts.RailLimit/2] = true // mainレールは必ず使用中

	for _, e := range events {
		switch e.Type {
		case BranchEventCreated:
			usedRails[e.ChildRailIndex] = true
		case BranchEventMerged:
			usedRails[e.ChildRailIndex] = false
		}
	}

	return usedRails
}

// GetParentRailIndex
// newRailIndexからmain方向に向かって一番近い使用中のレールを親として分岐させる
func GetParentRailIndex(childRailIndex int, usedRails [consts.RailLimit]bool) int {
	parentRailIndex := consts.RailLimit / 2

	if childRailIndex < consts.RailLimit/2 {
		for i := childRailIndex + 1; i < consts.RailLimit/2; i++ {
			if usedRails[i] {
				parentRailIndex = i
				break
			}
		}
	} else {
		for i := childRailIndex - 1; i >= consts.RailLimit/2; i-- {
			if usedRails[i] {
				parentRailIndex = i
				break
			}
		}
	}

	return parentRailIndex
}
