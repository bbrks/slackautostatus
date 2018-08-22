package main

import (
	"log"
	"os"
)

func logf(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}

var errorLogger = log.New(os.Stderr, "", log.LstdFlags)

func errorf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func fatalf(format string, v ...interface{}) {
	errorLogger.Fatalf(format, v...)
}
