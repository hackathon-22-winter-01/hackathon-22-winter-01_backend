package config

import "flag"

const (
	// PlayerLimit 部屋に入れるプレイヤーの上限数
	PlayerLimit = 4
	// RailLimit レール本数の上限
	RailLimit = 7
	// MaxLife プレイヤーの最大ライフ
	MaxLife float32 = 100.0
)

var (
	Port       = flag.Int("port", 8080, "port number")
	Production = flag.Bool("production", false, "enable production mode (default false)")
)

func ParseFlags() {
	flag.Parse()
}
