package cuisson

import (
	"errors"
	"time"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
)

var (
	ErrChaudronNonPréchauffé = errors.New("le chaudron n'a pas été préchauffé")
)

type Chaudron struct {
	TempératureActuelle recette.Température
	IngrédientEnCuisson recette.Ingrédient
}

// Cuire implements Appareil.
func (c *Chaudron) Cuire(ingrédient recette.Ingrédient, ingrédientObtenu string) error {
	if c.TempératureActuelle.IsZero() {
		return ErrChaudronNonPréchauffé
	}

	c.IngrédientEnCuisson = ingrédient
	c.IngrédientEnCuisson.Nom = ingrédientObtenu

	return nil
}

// Préchauffer implements Appareil.
func (c *Chaudron) Préchauffer(température recette.Température) chan bool {
	ch := make(chan bool, 1)
	go func() {
		for range température.Valeur() {
			time.Sleep(10 * time.Millisecond)
			c.TempératureActuelle.Chauffe(1)
		}
		ch <- c.TempératureActuelle.Égale(température)
	}()
	return ch
}

// Prélever implements Appareil.
func (c *Chaudron) Prélever() recette.Ingrédient {
	ingrédient := c.IngrédientEnCuisson
	c.IngrédientEnCuisson = recette.Ingrédient{}
	return ingrédient
}

// VérifierTempérature implements Appareil.
func (c *Chaudron) VérifierTempérature() recette.Température {
	return c.TempératureActuelle
}

func NouveauChaudron() Appareil {
	return &Chaudron{}
}
