package stock

import (
	"errors"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
)

var (
	ErrIngrédientNonTrouvé = errors.New("ingrédient non trouvé, ou quantité insuffisante")
)

//go:generate go tool mockgen -typed -destination=mocks/stock.go -package stockmocks . Stock
type Stock interface {
	VérifierDisponibilité(ingrédient recette.Ingrédient) bool
	RécupèrerIngrédient(ingrédient recette.Ingrédient) (recette.Ingrédient, error)
	StockerIngrédient(ingrédient recette.Ingrédient)
}
