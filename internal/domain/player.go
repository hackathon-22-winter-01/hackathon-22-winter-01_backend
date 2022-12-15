package domain

import "github.com/google/uuid"

// Player 対戦部屋内の各プレイヤーの情報
type Player struct {
	ID         uuid.UUID
	Name       string
	Main       *Rail
	Cards      []*Card
	Events     []*CardEvent
	LifeEvents []*LifeEvent
}

// NewPlayer 新しいプレイヤーを作成する
func NewPlayer(id uuid.UUID, name string) *Player {
	return &Player{
		ID:         id,
		Name:       name,
		Main:       NewRail(),
		Cards:      make([]*Card, 0),
		Events:     make([]*CardEvent, 0),
		LifeEvents: make([]*LifeEvent, 0),
	}
}
