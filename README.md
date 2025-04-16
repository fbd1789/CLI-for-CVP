// 📁 README.md
# cvaas-cli CLI (CloudVision CLI)

Une CLI modulaire en Go basée sur [Cobra](https://github.com/spf13/cobra) pour interagir avec Arista CloudVision-as-a-Service (cvaas-cli).

## 🚀 Fonctionnalités

- Créer des workspaces
- Créer et assigner des tags
- Lire l'inventaire des équipements
- Exécuter un processus complet (création + tag + assignation)

## 🛠 Prérequis

- Go 1.18+
- Un token cvaas-cli (`token.txt`)
- Une URL cvaas-cli (`url.txt`)

## 📦 Installation

```bash
go mod tidy
go build -o cvaas-cli
```

## 📚 Utilisation

```bash
./cvaas-cli --token token.txt --url url.txt [commande]
```

## 🔧 Commandes disponibles

### Créer un workspace
```bash
./cvaas-cli --token token.txt --url url.txt create workspace
```

### Créer un tag
```bash
./cvaas-cli --token token.txt --url url.txt create tag
```

### Voir l'inventaire des équipements
```bash
./cvaas-cli --token token.txt --url url.txt get devices
```

### Exécuter le process complet (création + tag + assignation)
```bash
./cvaas-cli --token token.txt --url url.txt run process
```

## 📁 Structure du projet

```
cvaas-cli/
├── main.go
├── client.go              # Connexion gRPC + lecture fichiers
├── actions.go             # Fonctions CloudVision (create, tag, assign...)
└── cmd/
    ├── root.go
    ├── create.go
    ├── get.go
    └── run.go
```

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

Développé avec ❤️ pour automatiser les tâches CloudVision.
