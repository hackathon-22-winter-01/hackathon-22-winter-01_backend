package domain

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
)

// Rail 路線の情報
type Rail struct {
	Index int
}

// NewRail 新しい路線を作成する
func NewRail(railIndex int) *Rail {
	return &Rail{
		Index: railIndex,
	}
}

type Rails [consts.RailLimit]*Rail

func NewRails(mainRail *Rail) Rails {
	rails := Rails{}

	// 真ん中にmainを入れる
	rails[consts.RailLimit/2] = mainRail

	return rails
}
