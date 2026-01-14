package recette

var PotionMagique = Recette{
	Nom:          "Potion Magique",
	Préchauffage: NouvelleTempérature(90),
	IngrédientsDeBase: []Ingrédient{
		{
			Nom:      "trèfle à 4 feuilles",
			Quantité: 8,
		},
		{
			Nom:      "fraise",
			Quantité: 12,
		},
		{
			Nom:      "once de lait de sanglier",
			Quantité: 4,
		},
		{
			Nom:      "pincée de curcuma",
			Quantité: 10,
		},
		{
			Nom:      "feuille de gui coupée à la serpe d'or",
			Quantité: 1,
		},
	},
	Déroulé: []Étape{
		{
			Nom: "Mélanger les trèfles et les fraises",
			Base: Ingrédient{
				Nom:      "trèfle à 4 feuilles",
				Quantité: 8,
			},
			Avec: Ingrédient{
				Nom:      "fraise",
				Quantité: 12,
			},
			Action:              Mélanger,
			NomIngrédientObtenu: "fraifles",
		},
		{
			Nom: "Préparer le lait doré",
			Base: Ingrédient{
				Nom:      "once de lait de sanglier",
				Quantité: 4,
			},
			Avec: Ingrédient{
				Nom:      "pincée de curcuma",
				Quantité: 10,
			},
			Action:              Mélanger,
			NomIngrédientObtenu: "lait doré",
		},
		{
			Nom: "Faire bouillir les fraifles",
			Base: Ingrédient{
				Nom:      "fraifles",
				Quantité: 14,
			},
			Action:              Bouillir,
			NomIngrédientObtenu: "fraifles bouillies",
		},
		{
			Nom: "Ajouter le gui",
			Base: Ingrédient{
				Nom:      "fraifles bouillies",
				Quantité: 14,
			},
			Avec: Ingrédient{
				Nom:      "feuille de gui coupée à la serpe d'or",
				Quantité: 1,
			},
			Action:              Mélanger,
			NomIngrédientObtenu: "fraifles au gui",
		},
		{
			Nom: "Verser le lait doré",
			Base: Ingrédient{
				Nom:      "fraifles au gui",
				Quantité: 14,
			},
			Avec: Ingrédient{
				Nom:      "lait doré",
				Quantité: 14,
			},
			Action:              Mélanger,
			NomIngrédientObtenu: "précurseur de potion magique",
		},
		{
			Nom: "Faire bouillir le précurseur",
			Base: Ingrédient{
				Nom:      "précurseur de potion magique",
				Quantité: 28,
			},
			Action:              Bouillir,
			NomIngrédientObtenu: "potion magique",
		},
	},
}
