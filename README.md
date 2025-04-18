// 📁 README.md
# cvaas-cli CLI (CloudVision CLI)

Une CLI modulaire en Go basée sur [Cobra](https://github.com/spf13/cobra) pour interagir avec Arista CloudVision-as-a-Service (cvaas-cli).

## 🚀 Fonctionnalités

- Lire l'inventaire des équipements

## 🛠 Prérequis

- Go 1.18+
- Un token cvaas-cli (`token.txt`)
- Une URL cvaas-cli (`url.txt`)

## 📦 Installation

```bash
go mod tidy
go build -o cvaas-cli
```
## 📁 Structure du projet

```
cvaas-cli/
├── main.go
├── internal   
|   ├── client.go              # Connexion gRPC + lecture fichiers
|   └── actions.go             # Fonctions CloudVision (create, tag, assign...)
└── cmd/
    ├── root.go
    ├── create.go
    ├── get.go
    └── run.go
```

## 📚 Utilisation

```bash
./cvaas-cli --token token.txt --url url.txt [commande]
```

## 📟 Commande `get devices`

Cette commande permet d'afficher l'inventaire des équipements (devices) connus par CVaaS.

---

### 🔧 Syntaxe

```bash
cvaas-cli get devices [options]
```

---

### 🔍 Options disponibles

| Option           | Description                                                                 |
|------------------|-----------------------------------------------------------------------------|
| `--model`        | Filtrer les équipements par modèle (ex: `cEOSLab`, `vEOS`, etc.)            |
| `--mlag`         | Afficher uniquement les équipements avec la fonctionnalité **MLAG** activée |
| `--danz`         | Afficher uniquement les équipements avec la fonctionnalité **DANZ** activée |

> ⚠️ Les filtres `--mlag` et `--danz` sont **mutuellement exclusifs** (ne peuvent pas être utilisés ensemble).

---

### ✅ Exemples

- 🔎 Afficher tous les équipements sans filtre :

```bash
cvaas-cli get devices --token token.txt --url url.txt
```

- 🔎 Filtrer par modèle `cEOSLab` :

```bash
cvaas-cli get devices --model cEOSLab --token token.txt --url url.txt
```

- 🔎 Afficher uniquement les équipements avec MLAG activé :

```bash
cvaas-cli get devices --mlag --token token.txt --url url.txt
```

- 🔎 Afficher uniquement ceux avec DANZ activé :

```bash
cvaas-cli get devices --danz --token token.txt --url url.txt
```

---

### 🛑 Erreurs possibles

- Si `--mlag` et `--danz` sont utilisés en même temps, la commande retournera une erreur :

```text
❌ Les filtres --mlag et --danz ne peuvent pas être utilisés en même temps.
```

---

### 🧠 Notes

- Le modèle (`--model`) peut être combiné avec l’un des deux filtres `--mlag` ou `--danz`
- La commande interroge directement l’inventaire CloudVision via gRPC

---







## 📌 Exemple de token.txt
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 📌 Exemple d'url.txt
```
cvaas-cli.example.arista.io:443
```

---

🔐 **NB :** Ne partagez jamais vos `token.txt` publiquement.

---

Développé pour automatiser les tâches CloudVision.
