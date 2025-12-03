package anlogger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Aninetix/core/aninterface"
)

type AnLoggerImpl struct {
	dir      string
	fileName string // optionnel
	debugOn  bool
}

var _ aninterface.AnLogger = (*AnLoggerImpl)(nil)

// ---- constructeur principal ----
func NewLogger(logDir string, debugOn bool) aninterface.AnLogger {
	os.MkdirAll(logDir, 0755)

	return &AnLoggerImpl{
		dir:     logDir,
		debugOn: debugOn,
	}
}

// ---- génère automatiquement YYYY-MM-DD-type.log si aucun fileName ----
func (l *AnLoggerImpl) getFilePath(t string) string {
	if l.fileName != "" {
		return filepath.Join(l.dir, l.fileName)
	}
	date := time.Now().Format("2006-01-02")
	return filepath.Join(l.dir, fmt.Sprintf("%s-%s.log", date, t))
}

// ---- crée un writer ----
func (l *AnLoggerImpl) writerFor(t string) *log.Logger {
	filePath := l.getFilePath(t)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return log.New(f, "", log.LstdFlags)
}

// ---- file:line ----
func callerInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// ---- Logs ----
func (l *AnLoggerImpl) Info(msg string) {
	l.writerFor("info").Printf("[INFO]  [%s] %s", callerInfo(), msg)
}

func (l *AnLoggerImpl) Error(msg string) {
	l.writerFor("error").Printf("[ERROR] [%s] %s", callerInfo(), msg)
}

func (l *AnLoggerImpl) Debug(msg string) {
	if !l.debugOn {
		return
	}
	l.writerFor("debug").Printf("[DEBUG] [%s] %s", callerInfo(), msg)
}

// ---- clone pour usage custom ----
func (l *AnLoggerImpl) WithFile(filename string) aninterface.AnLogger {
	return &AnLoggerImpl{
		dir:      l.dir,
		fileName: filename,
		debugOn:  l.debugOn,
	}
}
