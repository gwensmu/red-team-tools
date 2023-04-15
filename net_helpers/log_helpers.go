package net_helpers

import (
	"fmt"
	"log"
	"os"
	"time"
)

// todo this needs to be an absolute path
func InitLogFile(dir string, prefix string) {
	filename := fmt.Sprintf("%s/%s-%s.log", dir, prefix, time.Now())
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
}
