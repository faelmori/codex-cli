package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

func Warn(v ...interface{}) {
	WarnLogger.Println(v...)
}

func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}
