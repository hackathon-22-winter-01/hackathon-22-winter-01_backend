package log

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/config"
	"go.uber.org/zap"
)

var l *zap.Logger

func init() {
	var (
		_l  *zap.Logger
		err error
	)

	if *config.Production {
		_l, err = zap.NewProduction(zap.AddCaller())
	} else {
		_l, err = zap.NewDevelopment(zap.AddCaller())
	}

	if err != nil {
		panic(err)
	}

	l = _l
}

func L() *zap.Logger {
	return l
}
