package jst

import "time"

func Now() time.Time {
	return time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
}
