package main

import (
	"os"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/clix"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/cuisson"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/druide"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/stock"
)

func main() {
	// Création du garde manger et du chaudron
	gardeManger := stock.NouveauGardeManger([]*recette.Ingrédient{
		{
			Nom:      "trèfle à 4 feuilles",
			Quantité: 8,
		},
		{
			Nom:      "fraise",
			Quantité: 12,
		},
		{
			Nom:      "once de lait de sanglier",
			Quantité: 4,
		},
		{
			Nom:      "pincée de curcuma",
			Quantité: 10,
		},
		{
			Nom:      "feuille de gui coupée à la serpe d'or",
			Quantité: 1,
		},
	})
	chaudron := cuisson.NouveauChaudron()

	// Création du service métier avec injection
	druidePanoramix := druide.NouveauPagoramix(gardeManger, chaudron)

	// Création du CLIx
	cliAdapter := clix.NouveauCLIx(druidePanoramix)

	// Exécution
	if err := cliAdapter.Exécuter(); err != nil {
		os.Exit(1)
	}
}
