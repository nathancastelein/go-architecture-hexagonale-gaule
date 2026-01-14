package stock

import (
	"slices"
	"sync"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
)

type GardeManger struct {
	mutex       sync.Mutex
	ingrédients []*recette.Ingrédient
}

// RécupèrerIngrédient implements Stock.
func (gardeManger *GardeManger) RécupèrerIngrédient(ingrédientDemandé recette.Ingrédient) (recette.Ingrédient, error) {
	gardeManger.mutex.Lock()
	defer gardeManger.mutex.Unlock()

	for idx, ingrédientDuStock := range gardeManger.ingrédients {
		if peutRépondreÀLaDemande(*ingrédientDuStock, ingrédientDemandé) {
			ingrédientDuStock.Quantité = ingrédientDuStock.Quantité - ingrédientDemandé.Quantité
			if ingrédientDuStock.Quantité == 0 {
				gardeManger.ingrédients = slices.Delete(gardeManger.ingrédients, idx, idx+1)
			}
			return ingrédientDemandé, nil
		}
	}

	return recette.Ingrédient{}, ErrIngrédientNonTrouvé
}

func (gardeManger *GardeManger) StockerIngrédient(ingrédientÀStocker recette.Ingrédient) {
	gardeManger.mutex.Lock()
	defer gardeManger.mutex.Unlock()

	indexIngrédientDéjàEnStock := slices.IndexFunc(gardeManger.ingrédients, func(ingrédient *recette.Ingrédient) bool {
		return ingrédient.Nom == ingrédientÀStocker.Nom
	})
	if indexIngrédientDéjàEnStock >= 0 {
		gardeManger.ingrédients[indexIngrédientDéjàEnStock].Quantité += ingrédientÀStocker.Quantité
	} else {
		gardeManger.ingrédients = append(gardeManger.ingrédients, &ingrédientÀStocker)
	}
}

// VérifierDisponibilité implements Stock.
func (g *GardeManger) VérifierDisponibilité(ingrédientDemandé recette.Ingrédient) bool {
	return slices.ContainsFunc(g.ingrédients, func(ingrédientEnStock *recette.Ingrédient) bool {
		return peutRépondreÀLaDemande(*ingrédientEnStock, ingrédientDemandé)
	})
}

func peutRépondreÀLaDemande(ingrédientEnStock, ingrédientDemandé recette.Ingrédient) bool {
	return ingrédientEnStock.Nom == ingrédientDemandé.Nom && ingrédientEnStock.Quantité >= ingrédientDemandé.Quantité
}

func NouveauGardeManger(
	ingrédients []*recette.Ingrédient,
) Stock {
	return &GardeManger{
		ingrédients: ingrédients,
	}
}
