package monitorservice

import (
	"github.com/google/uuid"
	"time"
)

type LogType int32
type SeverityType int32

const (
	ERROR    LogType      = 0
	LOG      LogType      = 1
	WARNING  LogType      = 2
	LOW      SeverityType = 0
	MEDIUM   SeverityType = 1
	HIGH     SeverityType = 2
	CRITICAL SeverityType = 3
)

type Log struct {
	Id            uuid.UUID
	Name          string
	UserId	 	  uuid.UUID
	Type          LogType
	Content       interface {}
	Environment   string
	RuntimeInfo   interface { }
	Hostname      string
	StackTrace    interface { }
	Severity      SeverityType
	OccurredAt    time.Time

}

