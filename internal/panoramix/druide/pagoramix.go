package druide

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/cuisson"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/stock"
)

type Pagoramix struct {
	stock             stock.Stock
	appareilDeCuisson cuisson.Appareil
}

// Préparer implements Druide.
func (pagoramix *Pagoramix) Préparer(recetteÀPréparer recette.Recette) (*recette.Plat, error) {
	// Vérifier qu'on a tous les ingrédients pour lancer la recette
	for _, ingrédient := range recetteÀPréparer.IngrédientsDeBase {
		ingrédientDisponible := pagoramix.stock.VérifierDisponibilité(ingrédient)
		if !ingrédientDisponible {
			return nil, ErrIngrédientManquant
		}
	}

	// Préchauffer l'appareil de cuisson
	if !recetteÀPréparer.Préchauffage.IsZero() {
		slog.Info("préchauffage de l'appareil de cuisson", slog.Any("température", recetteÀPréparer.Préchauffage))
		préchauffé := pagoramix.appareilDeCuisson.Préchauffer(recetteÀPréparer.Préchauffage)
		ticker := time.NewTicker(100 * time.Millisecond)
		var préchauffageTerminé bool
		for !préchauffageTerminé {
			select {
			case <-ticker.C:
				slog.Info("température actuelle", slog.Any("température", pagoramix.appareilDeCuisson.VérifierTempérature()))
			case préchauffageTerminé = <-préchauffé:
				slog.Info("appareil de cuisson préchauffé", slog.Any("température", pagoramix.appareilDeCuisson.VérifierTempérature()))
			}
		}
	}

	// Parcourir les étapes
	var dernierIngrédientObtenu recette.Ingrédient
	for idx, étape := range recetteÀPréparer.Déroulé {
		slog.Info("réalisation de l'étape", slog.Any("étape", étape), slog.Int("index", idx))
		switch étape.Action {
		case recette.Mélanger:
			// Récupérer les ingrédients du stock
			ingrédientDeBase, err := pagoramix.stock.RécupèrerIngrédient(étape.Base)
			if err != nil {
				return nil, err
			}

			ingrédientAvec, err := pagoramix.stock.RécupèrerIngrédient(étape.Avec)
			if err != nil {
				return nil, err
			}

			// Mélanger les ingrédients
			nouvelIngrédient := pagoramix.Mélange(ingrédientDeBase, ingrédientAvec, étape.NomIngrédientObtenu)

			// Ajouter le nouvel ingrédient dans la liste des ingrédients disponibles
			pagoramix.stock.StockerIngrédient(nouvelIngrédient)
			dernierIngrédientObtenu = nouvelIngrédient
		case recette.Bouillir:
			ingrédientDeBase, err := pagoramix.stock.RécupèrerIngrédient(étape.Base)
			if err != nil {
				return nil, err
			}

			if err := pagoramix.appareilDeCuisson.Cuire(ingrédientDeBase, étape.NomIngrédientObtenu); err != nil {
				return nil, fmt.Errorf("échec de cuisson de l'ingrédient %s: %w", étape.Base.Nom, err)
			}

			dernierIngrédientObtenu = pagoramix.appareilDeCuisson.Prélever()

			pagoramix.stock.StockerIngrédient(dernierIngrédientObtenu)
		}
	}

	return pagoramix.Servir(dernierIngrédientObtenu), nil
}

func (pagoramix *Pagoramix) Mélange(ingrédientDeBase, ingrédientAvec recette.Ingrédient, nomIngrédientObtenu string) recette.Ingrédient {
	return recette.Ingrédient{
		Nom:      nomIngrédientObtenu,
		Quantité: ingrédientDeBase.Quantité + ingrédientAvec.Quantité,
	}
}

func (pagoramix *Pagoramix) Servir(ingrédient recette.Ingrédient) *recette.Plat {
	return &recette.Plat{
		Nom:      ingrédient.Nom,
		Quantité: ingrédient.Quantité,
	}
}

func NouveauPagoramix(
	stock stock.Stock,
	appareilDeCuisson cuisson.Appareil,
) Druide {
	return &Pagoramix{
		stock:             stock,
		appareilDeCuisson: appareilDeCuisson,
	}
}
