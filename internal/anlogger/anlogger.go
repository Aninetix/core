package anlogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Aninetix/core/aninterface"
)

// AnLoggerImpl structure principale
type AnLoggerImpl struct {
	info    *log.Logger
	error   *log.Logger
	debug   *log.Logger
	debugOn bool
}

// Assurer que AnLoggerImpl implémente l'interface
var _ aninterface.AnLogger = (*AnLoggerImpl)(nil)

// NewLogger crée un logger avec fichiers ou console et debug activable
func NewLogger(pathLog string, debugOn bool) aninterface.AnLogger {
	var infoWriter, errorWriter, debugWriter io.Writer

	if pathLog != "" {
		_ = os.MkdirAll(filepath.Dir(pathLog), 0755)
		f, err := os.OpenFile(pathLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Erreur ouverture du fichier log: %v", err)
			os.Exit(1)
		}
		infoWriter = f
		errorWriter = io.MultiWriter(os.Stderr, f)
		debugWriter = io.MultiWriter(os.Stdout, f)
	} else {
		infoWriter = os.Stdout
		errorWriter = os.Stderr
		debugWriter = os.Stdout
	}

	return &AnLoggerImpl{
		info:    log.New(infoWriter, "[INFO]  ", log.LstdFlags),
		error:   log.New(errorWriter, "[ERROR] ", log.LstdFlags),
		debug:   log.New(debugWriter, "[DEBUG] ", log.LstdFlags),
		debugOn: debugOn,
	}
}

// getCaller retourne fichier:ligne de la fonction appelante
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Info affiche un message d'information
func (l *AnLoggerImpl) Info(msg string) {
	caller := getCaller(2)
	l.info.Printf("[%s] %s", caller, msg)
}

// Error affiche un message d'erreur
func (l *AnLoggerImpl) Error(msg string) {
	caller := getCaller(2)
	l.error.Printf("[%s] %s", caller, msg)
}

// Debug affiche un message de debug si activé
func (l *AnLoggerImpl) Debug(msg string) {
	if !l.debugOn {
		return
	}
	caller := getCaller(2)
	l.debug.Printf("[%s] %s", caller, msg)
}

// EnableDebug permet d'activer/désactiver le debug dynamiquement
func (l *AnLoggerImpl) EnableDebug(on bool) {
	l.debugOn = on
}
