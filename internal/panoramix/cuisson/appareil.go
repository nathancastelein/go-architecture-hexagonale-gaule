package cuisson

import "github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"

//go:generate go tool mockgen -typed -destination=mocks/appareil.go -package cuissonmocks . Appareil
type Appareil interface {
	Préchauffer(température recette.Température) chan bool
	VérifierTempérature() recette.Température
	Cuire(ingrédient recette.Ingrédient, ingrédientObtenu string) error
	Prélever() recette.Ingrédient
}
