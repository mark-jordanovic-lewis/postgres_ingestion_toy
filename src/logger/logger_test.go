package logger

import (
	"os"
	"testing"
)

func TestLogFileInit(t *testing.T) {
	l := InitTimeStampedLog("TestLog")
	if l.filename == "" {
		t.Errorf("Log file name was not initialised")
	}
	m := InitLog("TestLog")
	if m.filename == "" {
		t.Errorf("Log file name was not initialised")
	}
}

func TestLogError(t *testing.T) {
	t.Log("Ensuring file creation and error logging")
	l := InitLog("TestLog")
	errTxt := "Test error\n"
	l.LogError(errTxt)
	logFile, err := os.Open(l.filename)
	if err != nil {
		t.Errorf("Logfile %v was not created", l.filename)
	}
	byteText := make([]byte, 30)
	nBytes, err := logFile.Read(byteText)
	if err != nil {
		t.Errorf("Could not read from log file")
	}
	if nBytes == 0 {
		t.Errorf("Zero bytes read from file")
	}
	if text := string(byteText[:nBytes]); text != errTxt {
		t.Errorf("Log file was not written to correctly")
	}
}
