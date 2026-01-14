package druide

import (
	"errors"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
)

var (
	ErrIngrédientManquant = errors.New("ingrédient manquant")
)

//go:generate go tool mockgen -typed -destination=mocks/druide.go -package druidemocks . Druide
type Druide interface {
	Préparer(recette recette.Recette) (*recette.Plat, error)
}
