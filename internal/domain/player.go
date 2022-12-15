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
