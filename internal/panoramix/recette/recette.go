package recette

import "log/slog"

type Recette struct {
	Nom               string
	Préchauffage      Température
	IngrédientsDeBase []Ingrédient
	Déroulé           []Étape
}

type Plat struct {
	Nom      string
	Quantité int
}

func (plat Plat) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("nom", plat.Nom),
		slog.Int("quantité", plat.Quantité))
}

type Étape struct {
	Nom                 string
	Base                Ingrédient
	Avec                Ingrédient
	Action              Action
	NomIngrédientObtenu string
}

func (étape Étape) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("nom", étape.Nom),
		slog.String("action", étape.Action.String()),
	)
}
