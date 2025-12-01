# Aninetix-Core

[![Go Reference](https://pkg.go.dev/badge/github.com/Aninetix/aninetix-core.svg)](https://pkg.go.dev/github.com/Aninetix/aninetix-core)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aninetix/aninetix-core)](https://goreportcard.com/report/github.com/Aninetix/aninetix-core)
[![CI](https://github.com/Aninetix/aninetix-core/actions/workflows/ci.yml/badge.svg)](https://github.com/Aninetix/aninetix-core/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**Aninetix-Core** est un framework Go modulaire conçu pour servir de base à tout développement d'applications Go. Il fournit une architecture extensible avec gestion de configuration, logging, et un système de modules basé sur les événements.

## ✨ Fonctionnalités

- 🔧 **Configuration JSON** - Chargement automatique de configuration typée
- 🚩 **Parsing de flags** - Support des flags CLI avec valeurs par défaut
- 📝 **Logging flexible** - Système de logging avec niveaux (Info, Error, Debug)
- 🔌 **Architecture modulaire** - Système de plugins avec communication par événements
- 📊 **Données système** - Accès aux informations système (OS, CPU, mémoire, réseau)
- 🔄 **Gestion du contexte** - Support natif de context.Context pour les arrêts gracieux

## 📦 Installation

```bash
go get github.com/Aninetix/aninetix-core
```

## 📚 Documentation

### Structure interne

La structure d’Aninet-Core est divisée en **3 parties principales** :

### 1. AnCore

* Contient le fichier `ancore.go` :

  * Définit la **struct `AnCore`**.
  * Contient 3 fonctions principales :

    * `InitCore()`
    * `BootCore()`
    * `Run()`
	
* Contient le dossier `AnWare` :

  * Sert de **middleware pour les modules** en dehors du core.
  * Permet de gérer la logique métier non standard.

### 2. AnInterface

* Contient les interfaces nécessaires pour utiliser les structs et données du core et etre disponible au module pour le typage.

* Exemples :

  * Logger
  * Données `StaticLocal`

### 3. Internal

* **AnConfig** : loader JSON pour configuration custom.
* **AnFlags** : loader pour arguments/flags du binaire.
* **AnLocal** : loader des données statiques (process ID, IP, etc.).
* **AnLogger** : loader du logger.
* **Helpers** : fonctions communes et utilitaires disponibles pour AnCore.

---

## Utilisation

### 1. Fichier `main.go`

Appel du core avec contexte :

```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flg, config, logger := ancore.Init_Core[param.Flags, param.Config]()

	AnCore := ancore.Boot_Core(flg, config, logger, ctx, cancel)

	AnCore.Run()

	<-ctx.Done()
	logger.Info("[MAIN] stop func Main(), context finish")
}
```

---

### 2. Création du dossier parametre  `param/param.go`

Définir les structs custom pour la configuration et les flags :

Exemple :

```go
type Config struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	PeerAnCluster string `json:"peer_cluster"`
}

type Flags struct {
	ConfigPath string        `flag:"config_path" default:"config.json" usage:"Chemin du fichier de configuration"`
	Debug      bool          `flag:"debug" default:"true" usage:"Activer le mode debug"`
	LogPath    string        `flag:"log_path" default:"_Data/server_default.log" usage:"Port d'écoute"`
	Timeout    time.Duration `flag:"timeout" default:"30s" usage:"Timeout (ex: 10s, 1m)"`
}
```

---

### 3. Création d’un module et structure d'appel

Chaque module doit implémenter l’interface standard :

```go
type AnModule interface {
	Name() string
	Param(ctx context.Context, in <-chan AnWareEvent, mw *AnWare)
	Start()
	Stop() error
}
```

##### Auto-enregistrement du module

```go
func init() {
	anware.RegisterModule("moduletest", NewAnModule)
}
```

##### Exemple complet de module minimal et prérequis

Info : le anware est générics avec les params, il faut donc les pointers avec le init et le NewAnmodule.

```go
package anmoduletest

import (
	"aninet-core/aninterface"
	"aninet-core/ancore/anware"
	"context"
	"fmt"
)

type AnConsolModule struct {
	ctx context.Context
	in  <-chan anware.AnWareEvent
	aw  *anware.AnWare

	localData aninterface.StaticData
	logger    aninterface.Logger

	Flags  *anparam.Flags
	Config *anparam.Config
}

// ---------- AUTO-ENREGISTREMENT DU MODULE ----------
func init() {
	anware.RegisterModule[*param.Flags, *param.Config](
		"moduletest",
		NewAnModule[*anparam.Flags, *anparam.Config],
	)
}

// ---------- CONSTRUCTEUR ----------
func NewAnModule[F *anparam.Flags, C *anparam.Config](
	local aninterface.StaticData,
	config C,
	flags F,
	logger aninterface.AnLogger,
) anware.AnModule {
	return &AnModule{
		anLocal:  local,
		anLogger: logger,
		Flags:    flags,
		Config:   config,
	}
}

// ---------- MÉTHODES INTERFACE AnModule ----------
func (m *AnConsolModule) Name() string {
	return m.name
}

func (m *AnConsolModule) Param(ctx context.Context, in <-chan anware.AnWareEvent, aw *anware.AnWare) {
	m.ctx = ctx
	m.in = in
	m.aw = aw
}

func (m *AnConsolModule) Stop() error {
	return nil
}

func (m *AnConsolModule) Start() {
	// A implémenter selon les besoins
	m.HandlePersonalisable()
}

```


##### Structure personalible et minimum pour intéragir avec les autres modules ou le core

```go

func (m *AnConsolModule) HandlePersonalisable() {
	m.logger.Info("[module test]")
	go func() {
		for {
			select {
			case <-m.ctx.Done():
				return
			case msg := <-m.in:
				m.msgChanAnWare(msg)
			}
		}
	}()
}
	

func (m *AnConsolModule) msgChanAnWare(msg anware.AnWareEvent) {
	switch msg.Type {
	case "status_response":
		fmt.Printf("[anChain] Status: %+v\n", msg.Data)
	case "peers_list_response":
		fmt.Printf("[memberlist] Peers: %+v\n", msg.Data)
	default:
		fmt.Printf("[console] Message reçu: %+v\n", msg)
	}
}

```

---

### 3. Import des modules

#### Structure d'appel

Pour l'appel, il faut créer un fichier contenant un import de chaque dossier de module.

Conseil, le mettre dans un dossier racine ou seront les modules `aninet_v2/AnModule/`

```go
package anmodule

import (
	_ "aninet_v2/AnModule/anmoduletest"
	_ "aninet_v2/AnModule/anmoduletest2"
)
```

Et dans `main.go` :

```go
import _ "aninet_v2/AnModule"
```

Cela permet au build d’intégrer les modules dans le core et de communiquer avec eux.

---

## Bonus : envoyer un message dans le core pour terminé le context (a implémenter en fonction des besoins)

```go
m.mw.Send(anware.AnWareEvent{
	Source: m.name,
	Target: "anWare",
	Type:   "exit",
})
```

---

Cette documentation couvre la structure, l’utilisation du core et la création/intégration de modules de manière claire et prête à l’emploi.

## 📊 Benchmarks

Les benchmarks sont exécutés sur chaque PR via GitHub Actions:

| Package | Operation | ns/op | B/op | allocs/op |
|---------|-----------|-------|------|-----------|
| anconfig | LoadConfig | ~8,334 | 1,144 | 12 |
| anflags | ParseFlags | ~1,623 | 704 | 12 |
| anlocal | LoadStaticData | ~1,862,444 | 93,361 | 169 |
| anlogger | Info | ~2,498 | 320 | 6 |
| helpers | GetFieldString | ~74 | 0 | 0 |

## 🧪 Tests

```bash
# Exécuter tous les tests
go test ./... -v

# Avec couverture
go test ./... -cover

# Exécuter les benchmarks
go test ./... -bench=. -benchmem
```

### Couverture actuelle

| Package | Couverture |
|---------|------------|
| anconfig | 100% |
| anflags | 90.3% |
| anlocal | 86.8% |
| anlogger | 88.9% |
| helpers | 100% |
| anware | 58.5% |

## 📁 Structure du projet

```
aninetix-core/
├── ancore/              # Core principal
│   ├── ancore.go        # InitCore, BootCore, Run
├── aninterface/         # Interfaces publiques
│   ├── AnLogger.go      # Interface Logger
│   ├── anData.go        # Interface StaticData
├── aninternal/          # Packages internes
│   ├── anconfig/        # Chargement config JSON
│   ├── anflags/         # Parsing des flags CLI
│   ├── anlocal/         # Données système
│   ├── anlogger/        # Implémentation du logger
│   └── helpers/         # Fonctions utilitaires
├── anware/              # Système de modules
│   ├── anware.go        # AnWare struct et constructeur
│   ├── method.go        # Méthodes (Run, Send, Broadcast)
│   ├── registry.go      # Registre des modules
├── examples/            # Exemples d'utilisation
├── LICENSE              # Licence Apache 2.0
└── README.md            # Ce fichier
```

## 🤝 Contribution

Les contributions sont les bienvenues ! Donc allez-y !

## 📄 Licence

Ce projet est sous licence [Apache 2.0](LICENSE).

## 📞 Support

- 📝 [Ouvrir une issue](https://github.com/Aninetix/aninetix-core/issues)
- 📖 [Documentation](https://pkg.go.dev/github.com/Aninetix/aninetix-core)

---

Made with ❤️ by the Aninetix Team