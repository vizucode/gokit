package logger

import (
	"sync"
	"time"
)

// Locker is container data
type Locker struct {
	data sync.Map
}

type (
	// Key of context
	Key int
	// Flags is key for store context
	Flags string
	// ServiceType is type for data logging
	ServiceType string
)

const (
	// logKey is key context for request http rest api
	LogKey = Key(31)

	// cron is type for logging cron
	cron ServiceType = "cron"

	// Flags for key of struct
	_StatusCode   Flags = "StatusCode"
	_Response     Flags = "Response"
	_LogMessages  Flags = "LogMessages"
	_ThirdParties Flags = "ThirdParties"
	_ErrorMessage Flags = "ErrorMessage"
	_UserCode     Flags = "UserCode"
	_Device       Flags = "Device"
	RequestId     Flags = "RequestId"
	_SaltKey      Flags = "SaltKey"

	// list type of logger
	debug   = "DEBUG"
	print   = "PRINT"
	err     = "ERROR"
	success = "Success request"
)

// DataLogger is standard output to terminal
type DataLogger struct {
	RequestId     string       `json:"request_id"`
	UserCode      string       `json:"user_code"`
	Device        string       `json:"device"`
	Type          ServiceType  `json:"type"`
	TimeStart     time.Time    `json:"time_start"`
	Service       string       `json:"service"`
	Host          string       `json:"host"`
	Endpoint      string       `json:"endpoint"`
	RequestMethod string       `json:"request_method"`
	RequestHeader string       `json:"request_header"`
	RequestBody   string       `json:"request_body"`
	StatusCode    int          `json:"status_code"`
	Response      interface{}  `json:"response"`
	ErrorMessage  string       `json:"error_message"`
	ExecTime      float64      `json:"exec_time"`
	LogMessages   []LogMessage `json:"log_message"`
	ThirdParties  []ThirdParty `json:"outgoing_log"`
}

// LogMessage is data logging for developer want to debug or error
type LogMessage struct {
	File    string `json:"file"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

// ThirdParty is data logging for any request to third party
type ThirdParty struct {
	ServiceTarget string  `json:"service_target"`
	URL           string  `json:"url"`
	RequestHeader string  `json:"request_header"`
	RequestBody   string  `json:"request_body"`
	Response      string  `json:"response"`
	Method        string  `json:"method"`
	StatusCode    int     `json:"status_code"`
	ExecTime      float64 `json:"exec_time"`
}

// String convert
func (st ServiceType) String() string {
	return string(st)
}
