package util

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// Debugging
const Debug = true

type logTopic string

const (
	DClient  logTopic = "CLNT"
	DCommit  logTopic = "CMIT"
	DDrop    logTopic = "DROP"
	DError   logTopic = "ERRO"
	DInfo    logTopic = "INFO"
	DLeader  logTopic = "LEAD"
	DLog     logTopic = "LOG1"
	DLog2    logTopic = "LOG2"
	DPersist logTopic = "PERS"
	DSnap    logTopic = "SNAP"
	DTerm    logTopic = "TERM"
	DTest    logTopic = "TEST"
	DTimer   logTopic = "TIMR"
	DTrace   logTopic = "TRCE"
	DVote    logTopic = "VOTE"
	DWarn    logTopic = "WARN"
)

var debugStart time.Time

func init() {
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetFlags(log.Flags() & ^(log.Ldate | log.Ltime))
	debugStart = time.Now()
}

func DPrintf(format string, a ...interface{}) {
	if !Debug {
		return
	}
	t := time.Since(debugStart).Microseconds()
	t /= 100
	prefix := fmt.Sprintf("%06d ", t)
	format = prefix + format + "\n"
	log.Printf(format, a...)
	return
}

func StartHTTPDebugger() {
	// http://localhost:6060/debug/pprof/
	DPrintf("%+v", http.ListenAndServe("localhost:6060", nil))
}
