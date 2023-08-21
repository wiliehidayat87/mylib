package mylib

import (
	"os"
	"time"
)

type (
	Utils struct {
		LogPath             string
		LogLevel            int
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
	}
)
