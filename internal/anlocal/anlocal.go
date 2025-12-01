package anlocal

import (
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/Aninetix/core/aninterface"
)

type staticDataImpl struct {
	process struct {
		PID        int
		PPID       int
		Executable string
		WorkingDir string
		Goroutines int
		MemStats   runtime.MemStats
	}
	machine struct {
		Hostname  string
		OS        string
		Arch      string
		GoVersion string
		Uptime    time.Duration
	}
	user struct {
		Username string
		UID      string
		GID      string
		HomeDir  string
	}
	cpu struct {
		NumCPU     int
		GOMAXPROCS int
	}
	memory struct {
		HeapUsed   uint64
		StackUsage uint64
	}
	network struct {
		Interfaces []net.Interface
	}
	anCoreID      string
	anCoreVersion string
}

// Factory pour créer l’interface
func LoadStaticData() aninterface.StaticData {
	data := &staticDataImpl{}

	// Processus
	data.process.PID = os.Getpid()
	data.process.PPID = os.Getppid()
	data.process.Executable, _ = os.Executable()
	data.process.WorkingDir, _ = os.Getwd()
	data.process.Goroutines = runtime.NumGoroutine()
	runtime.ReadMemStats(&data.process.MemStats)

	// Machine
	data.machine.Hostname, _ = os.Hostname()
	data.machine.OS = runtime.GOOS
	data.machine.Arch = runtime.GOARCH
	data.machine.GoVersion = runtime.Version()
	data.machine.Uptime = time.Since(time.Now().Add(-1 * time.Duration(data.process.MemStats.Sys)))

	// Utilisateur
	if currentUser, err := user.Current(); err == nil {
		data.user.Username = currentUser.Username
		data.user.UID = currentUser.Uid
		data.user.GID = currentUser.Gid
		data.user.HomeDir = currentUser.HomeDir
	}

	// CPU
	data.cpu.NumCPU = runtime.NumCPU()
	data.cpu.GOMAXPROCS = runtime.GOMAXPROCS(0)

	// Mémoire
	data.memory.HeapUsed = data.process.MemStats.HeapAlloc
	data.memory.StackUsage = data.process.MemStats.StackInuse

	// Réseau
	data.network.Interfaces, _ = net.Interfaces()

	data.anCoreID = generateAnId()
	data.anCoreVersion = getGitVersion()

	return data
}

// --- Méthodes pour implémenter l'interface ---

func (d *staticDataImpl) PID() int                    { return d.process.PID }
func (d *staticDataImpl) PPID() int                   { return d.process.PPID }
func (d *staticDataImpl) Executable() string          { return d.process.Executable }
func (d *staticDataImpl) WorkingDir() string          { return d.process.WorkingDir }
func (d *staticDataImpl) Goroutines() int             { return d.process.Goroutines }
func (d *staticDataImpl) MemStats() *runtime.MemStats { return &d.process.MemStats }

func (d *staticDataImpl) Hostname() string      { return d.machine.Hostname }
func (d *staticDataImpl) OS() string            { return d.machine.OS }
func (d *staticDataImpl) Arch() string          { return d.machine.Arch }
func (d *staticDataImpl) GoVersion() string     { return d.machine.GoVersion }
func (d *staticDataImpl) Uptime() time.Duration { return d.machine.Uptime }

func (d *staticDataImpl) Username() string { return d.user.Username }
func (d *staticDataImpl) UID() string      { return d.user.UID }
func (d *staticDataImpl) GID() string      { return d.user.GID }
func (d *staticDataImpl) HomeDir() string  { return d.user.HomeDir }

func (d *staticDataImpl) NumCPU() int     { return d.cpu.NumCPU }
func (d *staticDataImpl) GOMAXPROCS() int { return d.cpu.GOMAXPROCS }

func (d *staticDataImpl) HeapUsed() uint64   { return d.memory.HeapUsed }
func (d *staticDataImpl) StackUsage() uint64 { return d.memory.StackUsage }

func (d *staticDataImpl) Interfaces() []net.Interface { return d.network.Interfaces }

func (d *staticDataImpl) AnCoreID() string      { return d.anCoreID }
func (d *staticDataImpl) AnCoreVersion() string { return d.anCoreVersion }

// --- Fonctions internes ---
func generateAnId() string {
	return time.Now().Format("20060102150405.000000000")
}

func getGitVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--always")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}
