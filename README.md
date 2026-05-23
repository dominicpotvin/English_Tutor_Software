# English Tutor Project

Plateforme locale et personnelle de révision d'anglais. Backend Go, frontend React, base PostgreSQL, le tout en conteneurs Docker. Un connecteur MCP permet de co-construire le contenu avec Claude.

---

## Sommaire

- [Prérequis](#prérequis)
- [Démarrage](#démarrage)
- [Ports](#ports)
- [Le cursus](#le-cursus)
- [Connecteur MCP](#connecteur-mcp)
- [Mode développement](#mode-développement-optionnel)
- [Structure du projet](#structure-du-projet)
- [Architecture](#architecture)

---

## Prérequis

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

## Démarrage

```bash
docker compose up --build -d
```

La première construction prend quelques minutes. Ouvrir ensuite : <http://localhost:8095>

| Action | Commande |
|--------|----------|
| Arrêter la plateforme (les données sont conservées) | `docker compose down` |
| Tout réinitialiser (efface la progression et recharge le contenu) | `docker compose down -v` |

## Ports

| Service    | Port hôte | Rôle                    |
|------------|-----------|-------------------------|
| Frontend   | 8095      | Interface web (nginx)   |
| Backend    | 8096      | API REST (Go)           |
| PostgreSQL | 5434      | Base de données         |

## Le cursus

- **Introduction to Grammar** — 6 leçons complètes (16 topics, 96 exercices interactifs).
- **Vocabulary** — leçon 1 complète, plus 88 termes à réviser ; leçons 2 à 6 à compléter.
- **Level 1** — 6 leçons complètes (Plurals, Present, Present Continuous, Past, Future, Yes/No Questions). 14 topics, 87 exercices.
- **Level 2** — 6 leçons complètes (Stative Verbs, Past Continuous, Future Continuous, Present Perfect, Comparison of Adjectives, Possessive). 12 topics, 74 exercices.
- **Level 3** — 6 leçons complètes (Present Perfect, Present Perfect Continuous, Past Perfect, Future Perfect, Conditional Sentences, Passive Voice). 13 topics, 79 exercices.
- **Level 4** — 6 leçons complètes (Gerunds and Infinitives, Indirect Speech, Conditional Sentences, Passive Voice, Modals, The Tenses). 16 topics, 96 exercices.
- **Exam Prep Quiz** — 100 questions couvrant toute l'*Introduction to Grammar*.

> Contenu géré dans `backend/internal/seed/data/seed.json` au premier démarrage et complétable par la suite via le connecteur MCP (voir [`docs/seed-format.md`](docs/seed-format.md) pour le mode de chargement incrémental).

## Connecteur MCP

Le connecteur expose le contenu de la plateforme à Claude : lire et modifier les leçons, sujets, exercices, quiz et vocabulaire, et consulter la progression. C'est le canal prévu pour ajouter le contenu des *Level 1 à 4*.

Installation (une seule fois) :

```bash
cd mcp
npm install
npm run build
```

Le fichier `.mcp.json` à la racine enregistre le connecteur auprès de Claude Code. Redémarrer Claude Code dans ce dossier pour le charger.

> **Note** : la pile Docker doit être démarrée — le connecteur appelle l'API sur le port 8096.

## Mode développement (optionnel)

**Backend**

```bash
cd backend
go run ./cmd/server
```

Nécessite une base PostgreSQL ; voir `DATABASE_URL` dans `backend/internal/config/config.go`.

**Frontend**

```bash
cd frontend
npm install
npm run dev
```

Puis ouvrir <http://localhost:5174>. Les appels `/api` sont relayés vers le backend sur le port 8096.

## Structure du projet

```text
backend/    API Go : net/http, pgx, migrations SQL, chargement du contenu de base
frontend/   application React + Vite + TypeScript
mcp/        connecteur MCP (TypeScript)
docs/       format du fichier de contenu (seed)
```

## Architecture

- Au premier démarrage, le backend applique les migrations SQL puis charge le contenu de base depuis `backend/internal/seed/data/seed.json`.
- Le frontend (nginx) sert l'application et relaie `/api` vers le backend.
- Toute modification de contenu passe par l'API REST ou le connecteur MCP ; la progression de l'apprenant est calculée à partir des réponses enregistrées.
