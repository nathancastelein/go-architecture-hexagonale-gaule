package recette

type Ingrédient struct {
	Nom      string
	Quantité int
}

func (ingrédient Ingrédient) IsZero() bool {
	return ingrédient.Quantité == 0 && ingrédient.Nom == ""
}
