# Aninetix-Core

[![Go Reference](https://pkg.go.dev/badge/github.com/Aninetix/aninetix-core.svg)](https://pkg.go.dev/github.com/Aninetix/aninetix-core)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aninetix/aninetix-core)](https://goreportcard.com/report/github.com/Aninetix/aninetix-core)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**Aninetix-Core** est un framework Go modulaire servant de socle à des applications **robustes, extensibles et fortement typées**. Il fournit un **core minimal**, un **système de modules auto‑enregistrés**, une **gestion stricte de la configuration**, et une **communication événementielle** claire entre modules.

---

## ✨ Fonctionnalités clés

* 🔧 Chargement automatique de configuration JSON **typée**
* 🚩 Parsing des flags CLI avec valeurs par défaut
* 📝 Logging centralisé et extensible
* 🔌 Architecture modulaire auto‑enregistrée (plugin‑like, sans dépendance directe)
* 🔍 Validation **stricte** des configurations de modules
* 🔄 Gestion native du cycle de vie via `context.Context`

---

## 🧠 Architecture globale

Aninetix‑Core repose sur **3 piliers clairement séparés**.

### 1️⃣ AnCore — le socle applicatif

Responsable du **boot de l’application**, AnCore ne contient **aucune logique métier**.

Responsabilités :

* Parsing des flags globaux
* Chargement de la configuration applicative
* Initialisation du logger
* Création et lancement du système modulaire (`AnWare`)

Fonctions principales :

* `InitCore[F, C]()` — prépare flags, config et logger
* `BootCore()` — instancie le core runtime
* `Run()` — déclenche le chargement des modules

---

### 2️⃣ AnInterface — le contrat public

Expose les **interfaces partagées** entre le core et les modules :

* `AnLogger`
* `StaticData`
* Types d’événements

➡️ Garantit un **typage fort**, sans couplage entre modules et core.

---

### 3️⃣ AnWare — le système de modules

**Cœur du système modulaire**.

Responsabilités :

* Registre global des modules
* Auto‑chargement dynamique au runtime
* Injection de configuration typée
* Validation stricte des contrats modules
* Orchestration et communication événementielle

---

## 🚀 Cycle de vie d’une application

### 1. `main.go`

```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flags, config, logger := ancore.InitCore[anparam.Flags, anparam.Config]()
	core := ancore.BootCore(flags, config, logger, ctx, cancel)

	core.Run()
	<-ctx.Done()
}
```

➡️ Le `main` **ne connaît aucun module**.

---

## 🧩 Paramétrage global de l’application (`anparam`)

Le package `anparam` est **l’unique point d’entrée de l’application** pour :

* les flags CLI
* la configuration JSON
* la liste **exacte** des modules disponibles

### Exemple

```go
package anparam

import (
	anconsol "github.com/Aninetix/core_test/anmodules/anConsol"
	antest   "github.com/Aninetix/core_test/anmodules/anTest"
)

// Configuration applicative
// Reflète EXACTEMENT les modules disponibles
type Config struct {
	AnTest   antest.Config   `json:"anTest"`
	AnConsol anconsol.Config `json:"anConsol"`
}

// Flags globaux de l’application
type Flags struct {
	ConfigPath string `flag:"config_path" default:"data/config.json"`
	Debug      bool   `flag:"debug" default:"true"`
	LogPath    string `flag:"log_path" default:"data/server.log"`
}
```

➡️ **Une seule source de vérité**
➡️ Aucun doublon entre app et modules

---

## 🔌 Définition d’un module

### Interface standard

```go
type AnModule interface {
	Name() string
	Param(ctx context.Context, in <-chan AnWareEvent, mw *AnWare)
	Start()
	Stop() error
}
```

---

## 🧠 Auto‑enregistrement d’un module

Chaque module s’enregistre **automatiquement au build**, via `init()`.

```go
func init() {
	anware.RegisterModule(anware.ModuleDescriptor{
		Name:       "anTest",
		New:        NewModule,
		ConfigType: Config{},
	})
}
```

➡️ Le core **ne référence jamais explicitement un module**

---

## 🧪 Validation stricte de configuration (IMPORTANT)

Un module peut déclarer des **pré‑requis obligatoires**.

### Interface

```go
type ConfigValidator interface {
	Validate() error
}
```

### Exemple côté module

```go
func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("host is required")
	}
	if c.Port == 0 {
		return errors.New("port is required")
	}
	return nil
}
```

### Comportement

| Situation             | Résultat        |
| --------------------- | --------------- |
| Module absent du JSON | ❌ non chargé    |
| Champ requis manquant | ❌ non chargé    |
| Config valide         | ✅ module chargé |

➡️ **Pas de fallback silencieux**
➡️ **La configuration est un contrat**

---

## ⚙️ Auto‑chargement des modules

Lors du `Run()` :

1. Extraction de la sous‑configuration
2. Validation du contrat (`Validate()`)
3. Instanciation du module
4. Wiring des channels et du contexte

Les modules invalides sont **ignorés proprement**, sans panic.

---

## 📡 Communication inter‑modules

### Asynchrone

```go
m.mw.Send(anware.AnWareEvent{
	Source: m.Name(),
	Target: "anWare",
	Type:   "exit",
})
```

### Synchrone

```go
result, err := m.mw.SendSync(
	m.Name(),
	"anTest",
	"test_string",
	payload,
)
```

➡️ Le mode synchrone permet un **retour immédiat typé**
➡️ Le mode asynchrone reste non bloquant

---

## 📁 Structure du projet

```
aninetix-core/
├── ancore/          # Boot & orchestration
├── aninterface/     # Interfaces publiques
├── aninternal/      # Implémentations internes
├── anware/          # Système modulaire
├── examples/        # Exemples & modules de référence
└── README.md
```

---

## 🎯 Philosophie

* **Le core ne dépend de rien**
* **Les modules déclarent leurs besoins**
* **La configuration est le contrat**
* **L’import suffit pour activer**

---

Made with ❤️ by the Aninetix Team
