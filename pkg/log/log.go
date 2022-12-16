package log

import "go.uber.org/zap"

var l *zap.Logger

func init() {
	_l, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		panic(err)
	}

	l = _l
}

func L() *zap.Logger {
	return l
}
