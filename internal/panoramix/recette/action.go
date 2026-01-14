package recette

type Action string

const (
	Mélanger Action = "Mélanger"
	Bouillir Action = "Bouillir"
)

func (a Action) String() string {
	return string(a)
}
