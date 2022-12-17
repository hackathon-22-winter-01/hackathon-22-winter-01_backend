package domain

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
)

// Rail 路線の情報
type Rail struct {
	ID uuid.UUID
}

// NewRail 新しい路線を作成する
func NewRail(railID uuid.UUID) *Rail {
	return &Rail{
		ID: railID,
	}
}

type Rails [consts.RailLimit]*Rail

func NewRails(mainRail *Rail) Rails {
	rails := Rails{}

	// 真ん中にmainを入れる
	rails[consts.RailLimit/2] = mainRail

	return rails
}
