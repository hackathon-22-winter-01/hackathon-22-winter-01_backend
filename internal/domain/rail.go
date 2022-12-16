package domain

import "github.com/google/uuid"

// Rail 路線の情報
type Rail struct {
	ID       uuid.UUID
	HasBlock bool
}

// NewRail 新しい路線を作成する
func NewRail() *Rail {
	return &Rail{
		ID:       uuid.New(),
		HasBlock: false,
	}
}
