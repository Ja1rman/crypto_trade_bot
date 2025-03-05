package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func LoggerSetup() {
	logFile, err := os.OpenFile("./logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла логов: %v", err)
	}

	Logger = log.New(logFile, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
}
