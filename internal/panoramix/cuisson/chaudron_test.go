package cuisson

import (
	"testing"
	"testing/synctest"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
	"github.com/stretchr/testify/suite"
)

type ChaudronTestSuite struct {
	suite.Suite
}

func TestChaudron(t *testing.T) {
	suite.Run(t, &ChaudronTestSuite{})
}

func (suite *ChaudronTestSuite) TestCuire_NonPréchauffé() {
	// Arrange
	require := suite.Require()
	chaudron := NouveauChaudron()

	// Act
	err := chaudron.Cuire(recette.Ingrédient{}, "soupe froide")

	// Assert
	require.ErrorIs(err, ErrChaudronNonPréchauffé)
}

func (suite *ChaudronTestSuite) TestCuire() {
	// Arrange
	require := suite.Require()
	chaudron := &Chaudron{
		TempératureActuelle: recette.NouvelleTempérature(180),
	}

	// Act
	err := chaudron.Cuire(recette.Ingrédient{
		Nom:      "eau",
		Quantité: 10,
	}, "eau chaude")

	// Assert
	require.NoError(err)
	require.Equal(chaudron.IngrédientEnCuisson, recette.Ingrédient{
		Nom:      "eau chaude",
		Quantité: 10,
	})
}

func (suite *ChaudronTestSuite) TestPrélever() {
	// Arrange
	require := suite.Require()
	chaudron := &Chaudron{
		TempératureActuelle: recette.NouvelleTempérature(180),
		IngrédientEnCuisson: recette.Ingrédient{
			Nom:      "soupe chaude",
			Quantité: 10,
		},
	}

	// Act
	ingrédientPrélevé := chaudron.Prélever()

	// Assert
	require.Equal(ingrédientPrélevé, recette.Ingrédient{
		Nom:      "soupe chaude",
		Quantité: 10,
	})
	require.True(chaudron.IngrédientEnCuisson.IsZero())
}

func (suite *ChaudronTestSuite) TestPréchauffer() {
	synctest.Test(suite.T(), func(t *testing.T) {
		// Arrange
		require := suite.Require()
		chaudron := Chaudron{}
		températurePréchauffage := recette.NouvelleTempérature(180)

		// Act
		ch := chaudron.Préchauffer(températurePréchauffage)

		// Assert
		require.NotNil(ch)
		done := <-ch
		require.True(done)
		require.Equal(températurePréchauffage, chaudron.VérifierTempérature())
	})
}
