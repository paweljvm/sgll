package logging

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	info  = "[INFO] %s %s \r\n"
	warn  = "[WARN] %s %s \r\n"
	debug = "[DEBUG] %s %s \r\n"
	err   = "[ERROR] %s %s \r\n"
)

type Level int

type OutputType int

const (
	DEBUG Level = 3
	INFO  Level = 2
	WARN  Level = 1
	ERROR Level = 0
)

var level = INFO
var logFile *os.File
var lock = &sync.Mutex{}

func SetLevel(newLevel Level) {
	level = newLevel
}

func LogToFile(filePath string) {
	CloseLogFile()
	if len(filePath) == 0 {
		panic("Specify log file path")
	}
	var err error
	logFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to use file " + filePath + " due to error: " + err.Error())
	}
}

func Info(messageTemplate string, params ...interface{}) {
	writeLog(info, INFO, messageTemplate, isFileSpecified(), params...)
}

func Warn(messageTemplate string, params ...interface{}) {
	writeLog(warn, WARN, messageTemplate, isFileSpecified(), params...)
}

func Debug(messageTemplate string, params ...interface{}) {
	writeLog(debug, DEBUG, messageTemplate, isFileSpecified(), params...)
}

func Error(messageTemplate string, params ...interface{}) {
	writeLog(err, ERROR, messageTemplate, isFileSpecified(), params...)
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}

func writeLog(template string, minLevel Level, messageTemplate string, writeToFile bool, params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if level >= minLevel {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		logMessage := fmt.Sprintf(template, currentTime, resolveMessage(messageTemplate, params...))
		fmt.Println(logMessage)
		if logFile != nil {
			if _, erro := logFile.WriteString(logMessage); erro != nil {
				writeLog(err, ERROR, "Unable to save log to file {} due to {}", false, logFile.Name(), erro.Error())
			}
		}
	}
}

func resolveMessage(messageTemplate string, params ...interface{}) string {
	var message = messageTemplate
	for _, param := range params {
		message = strings.Replace(message, "{}", fmt.Sprintf("%v", param), 1)
	}
	return message
}

func isFileSpecified() bool {
	return logFile != nil
}
