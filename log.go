package mylib

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Logging struct {
	LogPath             string
	LogLevel            int
	LogDefault          string
	AccessLogName       string
	AccessErrorLogName  string
	AccessLogFormat     string
	AccessLogTimeFormat string
	AccessLogTimeZone   string
	BroadcastLogName    string
	CmpProcessorLogName string
	MoReceiverLogName   string
	MoProcessorLogName  string
	MtProcessorLogName  string
	DrReceiverLogName   string
	DrProcessorLogName  string
	ChargingLogName     string
	RetryLogName	      string
	PortalLogName    		string
	TrxLogName			    string
	LogThread           string
	LogFileName         string
	LogFileErr          string
	LogBehaviour        bool
}

// Used to define a full path of a log
func (l *Logging) GetStringPathLog(logName string) string {

	modLogPath := l.LogPath + "/" + logName

	//fmt.Println("logpath : " + logpath)

	if _, err := os.Stat(modLogPath); os.IsNotExist(err) {

		_ = os.Mkdir(modLogPath, 0777)
	}

	// Return log path with modify full path
	return modLogPath + "/" + GetFormatTime("20060102") + ".log"
}

// Instance for log setup method
// param :
// 1. @threadlog ( number string info for logging ) -> string
// 2. @logname ( string info for logging ) -> string
// 3. @logerr ( string info for logging error ) -> string
// returns :
// 1. @Logging -> struct interface

func (l *Logging) SetUpLog(appname string, threadlog string, logname string, logerr string) {

	// Set the file name of the configurations file
	viper.SetConfigName("logging")

	// Set the path to look for the configurations file
	viper.AddConfigPath(appname + "/config/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	prefix := "LOG"

	l.LogPath = viper.GetString(prefix + ".PATH")

	l.LogLevel = viper.GetInt(prefix + ".LEVEL")
	l.LogDefault = viper.GetString(prefix + ".DEFAULTNAME")

	l.AccessLogName = viper.GetString(prefix + ".ACCESSLOG")
	l.AccessErrorLogName = viper.GetString(prefix + ".ERRORLOG")
	l.AccessLogFormat = viper.GetString(prefix + ".ACCESSLOGFORMAT")
	l.AccessLogTimeFormat = viper.GetString(prefix + ".ACCESSLOGTIMEFORMAT")
	l.AccessLogTimeZone = viper.GetString(prefix + ".ACCESSLOGTIMEZONE")

	l.BroadcastLogName = viper.GetString(prefix + ".RENEWAL")
	l.CmpProcessorLogName = viper.GetString(prefix + ".CMP")
	l.MoReceiverLogName = viper.GetString(prefix + ".MORECEIVER")
	l.MoProcessorLogName = viper.GetString(prefix + ".MOPROCESSOR")

	l.MtProcessorLogName = viper.GetString(prefix + ".MTPROCESSOR")
	l.DrReceiverLogName = viper.GetString(prefix + ".DRRECEIVER")
	l.DrProcessorLogName = viper.GetString(prefix + ".DRPROCESSOR")
	l.ChargingLogName = viper.GetString(prefix + ".CHARGINGLOG")

	l.RetryLogName = viper.GetString(prefix + ".RETRY")
	l.PortalLogName = viper.GetString(prefix + ".PORTAL")
	l.TrxLogName = viper.GetString(prefix + ".TRX")

	if threadlog != "" {
		l.LogThread = Concat(threadlog, " ")
	} else {
		l.LogThread = Concat(GetLogId(), " ")
	}

	if logname != "" {

		switch strings.ToLower(logname) {
		case "access_log":
			l.LogFileName = l.AccessLogName
		case "error_log":
			l.LogFileName = l.AccessErrorLogName
		case "mo_receiver":
			l.LogFileName = l.MoReceiverLogName
		case "mo_processor":
			l.LogFileName = l.MoProcessorLogName
		case "dr_receiver":
			l.LogFileName = l.DrReceiverLogName
		case "dr_processor":
			l.LogFileName = l.DrProcessorLogName
		case "mt_processor":
			l.LogFileName = l.MtProcessorLogName
		case "broadcast":
			l.LogFileName = l.BroadcastLogName
		case "cmp_processor":
			l.LogFileName = l.CmpProcessorLogName
		case "charginglog":
			l.LogFileName = l.ChargingLogName
		case "retry":
			l.LogFileName = l.RetryLogName
		case "portal":
			l.LogFileName = l.PortalLogName
		case "trx":
			l.LogFileName = l.TrxLogName
		default:
			l.LogFileName = logname
		}

	} else if logname == "" {
		l.LogFileName = l.LogDefault
	}

	if logerr != "" {

		switch strings.ToLower(logerr) {
		case "error_log":
			l.LogFileErr = logerr
		default:
			l.LogFileErr = l.LogDefault
		}

	} else if logerr == "" {
		l.LogFileErr = l.LogFileName
	}

}

// LogWrite method
// param :
// 1. @loglevel ( option : 'info', 'debug', & 'error' ) -> string
// 2. @behaviour ( a boolean log wheter normal logging or error log ) -> bool
// 3. @logMsg ( a message string appear in a log file ) -> string
// returns :
// 1. -

func (l *Logging) Write(logLevel string, behaviour bool, logMsg string) {

	var logging bool
	logging = false

	var level int

	// Parsing loglevel

	if logLevel == "info" {
		level = 1
	} else if logLevel == "debug" {
		level = 2
	} else if logLevel == "error" {
		level = 3
	}

	// Check if log level is 0 (zero), will appear all
	if l.LogLevel == 0 {

		logging = true

	} else {

		// The log means if level log parameter is low than config set
		// then log will appear, if level log parameter is upper than config set
		// then log will never appear
		if level <= l.LogLevel {

			logging = true

		}

	}

	if level == 3 {

		logging = true

	}

	// Setup whether the error or normal
	// logpath of log file name
	l.LogBehaviour = behaviour

	if logging {

		var fullLogPath string

		if l.LogBehaviour {
			fullLogPath = l.GetStringPathLog(l.LogFileName)
		} else {
			fullLogPath = l.GetStringPathLog(l.LogFileErr)
		}

		f, err := os.OpenFile(fullLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 077)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		var threadlogging string

		if l.LogThread != "" {

			threadlogging = l.LogThread

		} else {

			threadlogging = Concat(GetLogId(), " ")
		}

		logger := log.New(f, threadlogging, log.LstdFlags)
		logger.Println(logLevel + " - " + logMsg)

	}
}
