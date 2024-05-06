package mylib

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

// Instance for log setup method
// param :
// 1. @threadlog ( number string info for logging ) -> string
// 2. @logname ( string info for logging ) -> string
// 3. @logerr ( string info for logging error ) -> string
// returns :
// 1. @Logging -> struct interface

func InitLog(Log Utils) *Utils {

	var l Utils

	l.AccessLogFormat = Log.AccessLogFormat
	l.AccessLogTimeFormat = Log.AccessLogTimeFormat
	l.LogPath = Log.LogPath
	l.LogLevelInit = Log.LogLevelInit
	l.TimeZone = Log.TimeZone

	os.Setenv("TZ", Log.TimeZone)

	return &l
}

func (l *Utils) SetUpLog(Log Utils) {
	l.LogThread = Log.LogThread + " "
	l.LogName = Log.LogName

	fullPathLog := l.GetStringPathLog(l.LogName)
	f, err := os.OpenFile(fullPathLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println(err)
	}
	l.LogOS = f
	l.LogOS.Chmod(fs.ModePerm)
	os.Chmod(fullPathLog, 0777)
}

// Used to define a full path of a log
func (l *Utils) GetStringPathLog(logName string) string {

	modLogPath := l.LogPath + "/" + logName

	//fmt.Println("logpath : " + logpath)

	if _, err := os.Stat(modLogPath); os.IsNotExist(err) {

		_ = os.Mkdir(modLogPath, 0777)
	}

	// Return log path with modify full path
	l.LogFileName = l.GetFormatTime("20060102")
	return modLogPath + "/" + l.LogFileName + ".log"
}

// Write method
// param :
// 1. @loglevel ( option : 'info', 'debug', & 'error' ) -> string
// 2. @logMsg ( a message string appear in a log file ) -> string
func (l *Utils) Write(logName string, logLevel string, logMsg string) {

	var (
		level        int
		allowLogging bool
	)

	// Parsing loglevel

	if logLevel == "info" {
		level = 1
	} else if logLevel == "debug" {
		level = 2
	} else if logLevel == "error" {
		level = 3
	}

	if l.LogLevelInit == 0 || l.LogLevelInit > 2 {

		allowLogging = true

	} else if l.LogLevelInit == 1 {

		if level == 1 {
			allowLogging = true
		} else {
			allowLogging = false
		}

	} else if l.LogLevelInit == 2 {

		if level == 1 || level == 2 {
			allowLogging = true
		} else {
			allowLogging = false
		}

	} else {

		allowLogging = false
	}

	if level == 3 {

		allowLogging = false

		fullPathLog := l.GetStringPathLog("error")
		f, err := os.OpenFile(fullPathLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			log.Println(err)
		}
		f.Chmod(fs.ModePerm)
		os.Chmod(fullPathLog, 0777)

		threadlogging := l.LogThread + " " + l.GetFormatTime("2006-01-02 15:04:05")

		logger := log.New(f, threadlogging, 0)
		logger.Println(" " + logLevel + " - " + logMsg)

		defer f.Close()

		fmt.Println(logMsg)
	}

	if allowLogging {

		if l.LogFileName != l.GetFormatTime("20060102") {

			l = InitLog(Utils{
				LogPath:      l.LogPath,
				LogLevelInit: l.LogLevelInit,
				TimeZone:     l.TimeZone,
			})
			l.SetUpLog(Utils{LogThread: l.GetUniqId(), LogName: logName})

			l.LogFileName = l.GetFormatTime("20060102")
		}

		threadlogging := l.LogThread + " " + l.GetFormatTime("2006-01-02 15:04:05")

		logger := log.New(l.LogOS, threadlogging, 0)
		logger.Println(" " + logLevel + " - " + logMsg)

		fmt.Println(logMsg)
	}
}
