package domain

import "github.com/google/uuid"

// Player 対戦部屋内の各プレイヤーの情報
type Player struct {
	ID             uuid.UUID
	Name           string
	Main           *Rail
	Cards          []*Card // TODO: remove
	JustCardEvents []*JustCardEvent
	BlockEvents    []*BlockEvent
	BranchEvents   []*BranchEvent
	LifeEvents     []*LifeEvent
	IsAlive        bool
}

// NewPlayer 新しいプレイヤーを作成する
func NewPlayer(id uuid.UUID, name string) *Player {
	return &Player{
		ID:           id,
		Name:         name,
		Main:         NewRail(),
		Cards:        make([]*Card, 0),
		BlockEvents:  make([]*BlockEvent, 0),
		BranchEvents: make([]*BranchEvent, 0),
		LifeEvents:   make([]*LifeEvent, 0),
		IsAlive:      true,
	}
}
