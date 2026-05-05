package debug

import (
	"fmt"
	"log"
	"time"
	"os"
)

type DisplayLogFunc func(string)

type Logger struct {
	NormLogger 			*log.Logger
	ErrLogger 			*log.Logger
	LogFile   			*os.File
	ErrFile   			*os.File
	DisplayLogFunc 	DisplayLogFunc // For NormLogger Only // I hate what I am seeting lol
}

func openLogFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("errors.NewLogger: Cannot open %q: %w", path, err)
	}
	return f, nil
}

func NewLogger(logPath string, errPath string, displayLogFunc DisplayLogFunc) (*Logger, error) {
	lf, err := openLogFile(logPath)
	if err != nil {
		return nil, err
	}
	ef, err := openLogFile(errPath)
	if err != nil {
		return nil, err
	}
	return &Logger{
		NormLogger: log.New(lf, "", 0),
		ErrLogger: log.New(ef, "", 0),
		LogFile: lf,
		ErrFile: ef,
		DisplayLogFunc: displayLogFunc,
	}, nil
}
 
func (l *Logger) Close() error {
	if l.LogFile != nil {
		return l.LogFile.Close()
	}
	return nil
}
 
func (l *Logger) LogErr(err error) { // Separation needed
	if l == nil || err == nil {
		return
	}
	if err != nil {
		if ae, ok := AsError(err); ok {
			l.ErrLogger.Printf(
				"[%s]:%s - Error=%v\nstack:\n%s\n",
				ae.Type,
				ae.Time.Format(time.RFC3339),
				ae.Err,
				ae.Stack,
			)
		} else {
			l.ErrLogger.Printf("[UnknownError]:%s - %v\n", time.Now().Format(time.RFC3339), err)
		}
	}
}

func (l *Logger) Log(logType LogType, context LogContext) {
	if l == nil || l.NormLogger == nil {
		return
	}
	logRes := GetLogFormat(logType, context)
	l.NormLogger.Print(logRes)
	l.DisplayLogFunc(logRes)
}
