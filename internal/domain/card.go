package domain

import "github.com/google/uuid"

// Card プレイヤーが使用するカードの情報
type Card struct {
	ID   uuid.UUID
	Type CardType
}

// CardType カードの種類
type CardType uint8

const (
	CardTypePullShark CardType = iota
	CardTypePairExtraordinaire
)

func NewCard(id uuid.UUID, typ CardType) *Card {
	return &Card{
		ID:   id,
		Type: typ,
	}
}
