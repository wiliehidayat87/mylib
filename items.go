package mylib

import (
	"os"
	"time"
)

type (
	Utils struct {
		LogPath             string
		LogLevelInit        int
		LogName             string
		LogOS               *os.File
		LogThread           string
		AccessLogFormat     string
		AccessLogTimeFormat string
		TimeZone            string
	}

	PHttp struct {
		Timeout            time.Duration
		KeepAlive          time.Duration
		IsDisableKeepAlive bool
		MaxIdleConns       int
		IdleConnTimeout    time.Duration
		DisableCompression bool
	}
)
