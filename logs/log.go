package logs

import (
	"bytes"
	"log"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	logs := bytes.Buffer{}
	InfoLogger = log.New(&logs, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(&logs, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(&logs, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
