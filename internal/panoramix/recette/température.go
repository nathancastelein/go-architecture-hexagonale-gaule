package recette

import "log/slog"

type Température struct {
	degré int
}

func (température Température) Valeur() int {
	return température.degré
}

func (température Température) Égale(températureÀComparer Température) bool {
	return température.degré == températureÀComparer.degré
}

func (température *Température) Chauffe(degrésGagnés int) {
	température.degré += degrésGagnés
}

func (température Température) LogValue() slog.Value {
	return slog.IntValue(température.degré)
}

func (température Température) IsZero() bool {
	return température.degré == 0
}

func NouvelleTempérature(degré int) Température {
	return Température{
		degré: degré,
	}
}
