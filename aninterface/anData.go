package aninterface

import (
	"net"
	"runtime"
	"time"
)

type StaticData interface {
	// Processus
	PID() int
	PPID() int
	Executable() string
	WorkingDir() string
	Goroutines() int
	MemStats() *runtime.MemStats

	// Machine
	Hostname() string
	OS() string
	Arch() string
	GoVersion() string
	Uptime() time.Duration

	// Utilisateur
	Username() string
	UID() string
	GID() string
	HomeDir() string

	// CPU
	NumCPU() int
	GOMAXPROCS() int

	// Mémoire
	HeapUsed() uint64
	StackUsage() uint64

	// Réseau
	Interfaces() []net.Interface

	// Core info
	AnCoreID() string
	AnCoreVersion() string
}
