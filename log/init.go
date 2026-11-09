package log

import (
	"github.com/cinemaxxi-team/xxi-library/logging"
)

var Logging *logging.Logger

func InitLog() {
	cfg := [2]logging.LogConfig{
		{
			FileName:   "logs/access.log",
			MaxSize:    50,
			MaxBackups: 3,
			MaxAge:     28,
		},
		{
			FileName:   "logs/error.log",
			MaxSize:    50,
			MaxBackups: 3,
			MaxAge:     28,
		},
	}

	Logging = logging.InitLogger(cfg)
}
