package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	filename string
}

func (l Logger) LogError(errStr string) {
	f, err := os.OpenFile(
		l.filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(errStr)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func InitTimeStampedLog(prefix string) Logger {
	t := time.Now()
	return Logger{filename: fmt.Sprintf("%v-%v.log", prefix, t)}
}

func InitLog(prefix string) Logger {
	return Logger{filename: fmt.Sprintf("%v.log", prefix)}
}
