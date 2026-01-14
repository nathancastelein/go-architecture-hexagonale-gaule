package stock

import (
	"slices"
	"testing"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
	"github.com/stretchr/testify/suite"
)

type GardeMangerTestSuite struct {
	suite.Suite
}

func (suite *GardeMangerTestSuite) TestVérifierDisponibilité_IngrédientEnStock() {
	// Arrange
	require := suite.Require()
	gardeManger := NouveauGardeManger([]*recette.Ingrédient{
		{
			Nom:      "sanglier",
			Quantité: 3,
		},
		{
			Nom:      "gui",
			Quantité: 15,
		},
	})

	// Act
	ingrédientDisponible := gardeManger.VérifierDisponibilité(recette.Ingrédient{
		Nom:      "sanglier",
		Quantité: 2,
	})

	// Assert
	require.True(ingrédientDisponible)
}

func (suite *GardeMangerTestSuite) TestVérifierDisponibilité_IngrédientManquant() {
	// Arrange
	require := suite.Require()
	gardeManger := NouveauGardeManger([]*recette.Ingrédient{
		{
			Nom:      "sanglier",
			Quantité: 3,
		},
		{
			Nom:      "gui",
			Quantité: 15,
		},
	})

	// Act
	ingrédientDisponible := gardeManger.VérifierDisponibilité(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	require.False(ingrédientDisponible)
}

func (suite *GardeMangerTestSuite) TestVérifierDisponibilité_QuantitéInsuffisante() {
	// Arrange
	require := suite.Require()
	gardeManger := NouveauGardeManger([]*recette.Ingrédient{
		{
			Nom:      "sanglier",
			Quantité: 3,
		},
		{
			Nom:      "gui",
			Quantité: 15,
		},
	})

	// Act
	ingrédientDisponible := gardeManger.VérifierDisponibilité(recette.Ingrédient{
		Nom:      "sangliers",
		Quantité: 5,
	})

	// Assert
	require.False(ingrédientDisponible)
}

func (suite *GardeMangerTestSuite) TestRécupérerIngrédient_IngrédientEnStock() {
	// Arrange
	require := suite.Require()
	gardeManger := GardeManger{
		ingrédients: []*recette.Ingrédient{
			{
				Nom:      "sanglier",
				Quantité: 3,
			},
			{
				Nom:      "gui",
				Quantité: 15,
			},
		},
	}

	// Act
	ingrédientObtenu, err := gardeManger.RécupèrerIngrédient(recette.Ingrédient{
		Nom:      "sanglier",
		Quantité: 2,
	})

	// Assert
	require.NoError(err)
	require.Equal(ingrédientObtenu, recette.Ingrédient{Nom: "sanglier", Quantité: 2})
	require.True(slices.ContainsFunc(gardeManger.ingrédients, func(ingrédient *recette.Ingrédient) bool {
		return ingrédient.Nom == "sanglier" && ingrédient.Quantité == 1
	}))
}

func (suite *GardeMangerTestSuite) TestRécupérerIngrédient_IngrédientEnStock_ToutRetirer() {
	// Arrange
	require := suite.Require()
	gardeManger := GardeManger{
		ingrédients: []*recette.Ingrédient{
			{
				Nom:      "sanglier",
				Quantité: 3,
			},
			{
				Nom:      "gui",
				Quantité: 15,
			},
		},
	}

	// Act
	ingrédientObtenu, err := gardeManger.RécupèrerIngrédient(recette.Ingrédient{
		Nom:      "sanglier",
		Quantité: 3,
	})

	// Assert
	require.NoError(err)
	require.Equal(ingrédientObtenu, recette.Ingrédient{Nom: "sanglier", Quantité: 3})
	require.Len(gardeManger.ingrédients, 1)
	require.Equal(&recette.Ingrédient{
		Nom:      "gui",
		Quantité: 15,
	}, gardeManger.ingrédients[0])
}

func (suite *GardeMangerTestSuite) TestRécupérerIngrédient_IngrédientManquant() {
	// Arrange
	require := suite.Require()
	gardeManger := GardeManger{
		ingrédients: []*recette.Ingrédient{
			{
				Nom:      "sanglier",
				Quantité: 3,
			},
			{
				Nom:      "gui",
				Quantité: 15,
			},
		},
	}

	// Act
	ingrédientObtenu, err := gardeManger.RécupèrerIngrédient(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	require.ErrorIs(err, ErrIngrédientNonTrouvé)
	require.Equal(ingrédientObtenu, recette.Ingrédient{})
}

func (suite *GardeMangerTestSuite) TestStockerIngrédient_GardeMangerVide() {
	// Arrange
	require := suite.Require()
	gardeManger := GardeManger{
		ingrédients: []*recette.Ingrédient{},
	}

	// Act
	gardeManger.StockerIngrédient(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	require.Len(gardeManger.ingrédients, 1)
	require.Equal(&recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	}, gardeManger.ingrédients[0])
}

func (suite *GardeMangerTestSuite) TestStockerIngrédient_GardeMangerDéjàRempli() {
	// Arrange
	require := suite.Require()
	gardeManger := GardeManger{
		ingrédients: []*recette.Ingrédient{
			{
				Nom:      "cervoise",
				Quantité: 5,
			},
		},
	}

	// Act
	gardeManger.StockerIngrédient(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	require.Len(gardeManger.ingrédients, 1)
	require.Equal(&recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 10,
	}, gardeManger.ingrédients[0])
}

func TestGardeManger(t *testing.T) {
	suite.Run(t, &GardeMangerTestSuite{})
}
