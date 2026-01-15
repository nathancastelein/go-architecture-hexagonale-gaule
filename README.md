# Pagoramix - Le Gestionnaire de Recettes Gauloises

*Par Cétautomatix, forgeron du village et développeur Go*

---

## Par Toutatis, quelle histoire !

Écoutez-moi bien, vous autres ! Moi, Cétautomatix, le meilleur forgeron de toute la Gaule (et accessoirement le seul à savoir coder proprement), je vais vous raconter comment j'ai dû sauver notre village d'une catastrophe sans précédent.

Figurez-vous que notre druide Panoramix a eu un petit... disons... accident de mémoire. Résultat ? Plus personne ne connaissait la recette de la potion magique ! Vous imaginez ? Notre village, sans potion magique, face aux Romains ? Autant dire qu'Ordralfabétix aurait pu leur vendre son poisson pas frais pour les faire fuir, ça aurait eu le même effet !

Bref, notre cher Astérix a décidé qu'il fallait créer une application pour **conserver à jamais nos précieuses recettes gauloises**. Et qui s'est porté volontaire ? Votre serviteur, évidemment !

Bon, Ordralfabétix aussi s'est proposé avec son Java... Pfff ! Du Java ! Aussi frais que son poisson, son code ! Moi, je martèle du **Go** tous les jours à ma forge, un langage solide comme mes épées, simple comme une bonne baffe sur un Romain !

---

## L'Architecture Hexagonale, ou comment forger du code solide

Alors, j'aurais pu faire comme Ordralfabétix et empiler du code comme il empile ses poissons pourris. Mais non ! Moi, Cétautomatix, j'ai de la rigueur ! J'ai utilisé l'**architecture hexagonale**, aussi appelée "Ports et Adaptateurs".

C'est comme quand je forge une épée :
- Le **cœur** (le domaine), c'est la lame : le métal pur, sans compromis
- Les **ports**, ce sont les manches : l'interface entre la lame et le monde extérieur
- Les **adaptateurs**, ce sont les différentes poignées qu'on peut y mettre : cuir, bois, or...

```
                    ┌─────────────────┐
                    │   CLI (CLIx)    │  ← Vous tapez vos commandes ici
                    │    (Cobra)      │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │    PAGORAMIX    │  ← Le cerveau de l'opération
                    │ (Service Druide)│     (oui, j'ai fait un jeu de mots)
                    └───┬─────────┬───┘
                        │         │
               ┌────────▼───┐ ┌───▼────────┐
               │   Stock    │ │  Appareil  │  ← Les ports
               │   (Port)   │ │   (Port)   │
               └────────┬───┘ └───┬────────┘
                        │         │
        ┌───────────────┼─────────┼───────────────┐
        │               │         │               │
   ┌────▼─────┐   ┌─────▼───┐ ┌───▼─────┐        │
   │GardeManger│  │  Redix  │ │ Chaudron│        │
   │ (Mémoire) │  │ (Redis) │ │         │  ← Les adaptateurs
   └──────────┘   └─────────┘ └─────────┘
```

---

## Structure du Projet

```
go-architecture-hexagonale-gaule/
├── cmd/druide/              # Point d'entrée - là où tout se build
│   └── main.go
├── internal/
│   └── panoramix/           # Package principal (en hommage à notre druide)
│       ├── recette/
│       │   ├── recette.go         # Qu'est-ce qu'une recette ?
│       │   ├── ingrédient.go      # Qu'est-ce qu'un ingrédient ?
│       │   ├── action.go          # Mélanger, Bouillir...
│       │   ├── température.go     # Gestion de la température
│       │   └── potion_magique.go  # LA recette secrète !
│       │
│       ├── druide/          # Le SERVICE - l'orchestrateur
│       │   ├── druide.go          # L'interface du druide
│       │   └── pagoramix.go       # L'implémentation (Go + Panoramix, vous l'avez ?)
│       │
│       ├── stock/           # PORT & ADAPTATEURS pour le stockage
│       │   ├── stock.go           # L'interface (le port)
│       │   ├── garde_manger.go    # Adaptateur stockage en mémoire
│       │   └── redix.go           # Adaptateur Redis
│       │
│       ├── cuisson/         # PORT & ADAPTATEURS pour la cuisson
│       │   ├── appareil.go        # L'interface (le port)
│       │   └── chaudron.go        # Le bon vieux chaudron gaulois
│       │
│       └── clix/            # ADAPTATEUR d'entrée (CLI)
│           └── clix.go            # Interface ligne de commande
│
├── go.mod
└── go.sum
```

---

## Utilisation

### Prérequis
- Go 1.25.5 ou supérieur (on n'utilise pas de vieilleries ici, contrairement au poisson d'Ordralfabétix)

### Lancer l'application

```bash
# Préparer la potion magique
go run ./cmd/druide potion-magique

# Afficher l'aide
go run ./cmd/druide --help
```

### Lancer les tests

```bash
# Tous les tests
go test ./...
```

---

## Technologies Utilisées

| Technologie        | Utilisation         | Commentaire de Cétautomatix                                   |
| ------------------ | ------------------- | ------------------------------------------------------------- |
| **Go 1.25**        | Langage             | Solide comme mes enclumes !                                   |
| **testify**        | Assertions          | Pour être sûr que ça marche                                   |
| **gomock**         | Mocks               | Pour tester sans dépendances                                  |
| **Cobra**          | Framework CLI       | Pour taper des commandes comme je tape sur le fer             |
| **Redis**          | Stockage persistant | Parce que la mémoire, c'est pas fiable (demandez à Panoramix) |
| **Testcontainers** | Tests d'intégration | Des vrais tests, pas du poisson !                             |

---

## Les Principes de l'Architecture Hexagonale

### 1. Isolation du domaine
Le code métier (dans `recette/` et `druide/`) n'a **aucune dépendance externe**. Il ne connaît ni Redis, ni Cobra, ni rien d'autre. C'est du code pur, comme le métal avant que je le forge.

### 2. Inversion des dépendances
Les adaptateurs dépendent des ports (interfaces), jamais l'inverse. C'est comme ça qu'on obtient du code flexible !

### 3. Testabilité
Grâce aux interfaces, on peut tout mocker et tout tester. Mes tests sont verts comme les forêts gauloises !

### 4. Interchangeabilité
Envie de remplacer Redis par PostgreSQL ? Pas de souci, il suffit de créer un nouvel adaptateur. Le reste du code n'en saura rien !

---

## Comparaison Go vs Java

Ordralfabétix a fait la même application en Java. Voici pourquoi ma version est meilleure :

| Aspect                   | Go (moi)                  | Java (Ordralfabétix)              |
| ------------------------ | ------------------------- | --------------------------------- |
| Interfaces               | Implicites, élégantes     | Explicites, verbeuses             |
| Injection de dépendances | Simple, via constructeurs | Framework lourd (Spring)          |
| Compilation              | Rapide comme Astérix      | Lent comme un légionnaire fatigué |
| Binaire final            | Un seul fichier           | Une armée de JARs                 |

*Note : Ordralfabétix prétend que son poisson... euh, son code Java est "plus frais". Ne l'écoutez pas.*

---

## La Conférence

Cette application a été créée pour une conférence :

> **L'architecture hexagonale au pays des irréductibles développeurs**
>
> Une conférence épique, pragmatique et pleine d'humour, pour découvrir (ou redécouvrir) l'architecture hexagonale par la pratique !

Présentée par deux irréductibles développeurs :
- **Cétautomatix** - Version Go, disponible ici même
- **Ordralfabétix** - Version Java, disponible [ici](https://github.com/ambreperson/java-architecture-hexagonale-gaule)

Les deux rivaux explorent les concepts fondamentaux de l'architecture hexagonale et mettent en lumière, à travers des implémentations concrètes en Java et en Go, les différences d'approche selon les écosystèmes !

---

## Contribution

Si vous voulez contribuer, n'hésitez pas ! Mais attention :
- Du code propre, ou je vous forge une nouvelle tête
- Des tests, sinon c'est comme vendre du poisson pas frais
- Respectez l'architecture hexagonale

---

*"Ils sont fous ces Romains !" - Obélix*

*"Ils sont fous ces développeurs qui n'utilisent pas l'architecture hexagonale !" - Cétautomatix*

---

**Fait avec force coups de marteau par Cétautomatix, forgeron et développeur Go du village gaulois**
