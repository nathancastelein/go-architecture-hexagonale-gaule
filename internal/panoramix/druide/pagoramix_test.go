package druide

import (
	"testing"

	cuissonmocks "github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/cuisson/mocks"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
	stockmocks "github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/stock/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PréparerTraiflesTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
}

func (suite *PréparerTraiflesTestSuite) SetupSuite() {
	suite.mockController = gomock.NewController(suite.T())
}

type fauxAppareilCuisson struct {
	*cuissonmocks.MockAppareil
	ingrédientEnCuisson recette.Ingrédient
}

func (appareil *fauxAppareilCuisson) Cuire(ingrédient recette.Ingrédient, nomIngrédientObtenu string) error {
	appareil.ingrédientEnCuisson = ingrédient
	appareil.ingrédientEnCuisson.Nom = nomIngrédientObtenu
	return nil
}

func (appareil *fauxAppareilCuisson) Prélever() recette.Ingrédient {
	return appareil.ingrédientEnCuisson
}

func (suite *PréparerTraiflesTestSuite) TestPlatPréparé() {
	// Arrange
	require := suite.Require()

	températureDePréchauffage := recette.NouvelleTempérature(180)

	recetteFraifles := recette.Recette{
		Préchauffage: températureDePréchauffage,
		IngrédientsDeBase: []recette.Ingrédient{
			{
				Nom:      "fraise",
				Quantité: 8,
			},
			{
				Nom:      "trèfle à 4 feuilles",
				Quantité: 13,
			},
		},
		Déroulé: []recette.Étape{
			{
				Nom: "Mélanger les trèfles et les fraises",
				Base: recette.Ingrédient{
					Nom:      "fraise",
					Quantité: 8,
				},
				Avec: recette.Ingrédient{
					Nom:      "trèfle à 4 feuilles",
					Quantité: 13,
				},
				Action:              recette.Mélanger,
				NomIngrédientObtenu: "fraifles",
			},
			{
				Nom: "Faire bouillir les fraifles",
				Base: recette.Ingrédient{
					Nom:      "fraifles",
					Quantité: 21,
				},
				Avec:                recette.Ingrédient{},
				Action:              recette.Bouillir,
				NomIngrédientObtenu: "fraifles bouillies",
			},
		},
	}

	fauxStock := stockmocks.NewMockStock(suite.mockController)
	fauxStock.EXPECT().VérifierDisponibilité(gomock.Any()).AnyTimes().Return(true)
	fauxStock.EXPECT().RécupèrerIngrédient(gomock.Any()).AnyTimes().DoAndReturn(func(ingrédient recette.Ingrédient) (recette.Ingrédient, error) {
		return ingrédient, nil
	})
	fauxStock.EXPECT().StockerIngrédient(gomock.Any()).AnyTimes().Return()

	fauxAppareilCuisson := &fauxAppareilCuisson{
		MockAppareil: cuissonmocks.NewMockAppareil(suite.mockController),
	}
	fauxAppareilCuisson.EXPECT().Préchauffer(gomock.Any()).AnyTimes().DoAndReturn(func(température recette.Température) chan bool {
		ch := make(chan bool, 1)
		ch <- true
		return ch
	})
	fauxAppareilCuisson.EXPECT().VérifierTempérature().AnyTimes().Return(températureDePréchauffage)

	préparateur := NouveauPagoramix(fauxStock, fauxAppareilCuisson)

	platAttendu := &recette.Plat{
		Nom:      "fraifles bouillies",
		Quantité: 21,
	}

	// Act
	platObtenu, err := préparateur.Préparer(recetteFraifles)

	// Assert
	require.NoError(err)
	require.Equal(platAttendu, platObtenu)
}

func (suite *PréparerTraiflesTestSuite) TestPotionMagiquePréparée() {
	// Arrange
	require := suite.Require()

	températureDePréchauffage := recette.PotionMagique.Préchauffage

	fauxStock := stockmocks.NewMockStock(suite.mockController)
	fauxStock.EXPECT().VérifierDisponibilité(gomock.Any()).AnyTimes().Return(true)
	fauxStock.EXPECT().RécupèrerIngrédient(gomock.Any()).AnyTimes().DoAndReturn(func(ingrédient recette.Ingrédient) (recette.Ingrédient, error) {
		return ingrédient, nil
	})
	fauxStock.EXPECT().StockerIngrédient(gomock.Any()).AnyTimes().Return()

	fauxAppareilCuisson := &fauxAppareilCuisson{
		MockAppareil: cuissonmocks.NewMockAppareil(suite.mockController),
	}
	fauxAppareilCuisson.EXPECT().Préchauffer(gomock.Any()).AnyTimes().DoAndReturn(func(température recette.Température) chan bool {
		ch := make(chan bool, 1)
		ch <- true
		return ch
	})
	fauxAppareilCuisson.EXPECT().VérifierTempérature().AnyTimes().Return(températureDePréchauffage)
	fauxAppareilCuisson.EXPECT().Cuire(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	préparateur := NouveauPagoramix(fauxStock, fauxAppareilCuisson)

	platAttendu := &recette.Plat{
		Nom:      "potion magique",
		Quantité: 28,
	}

	// Act
	platObtenu, err := préparateur.Préparer(recette.PotionMagique)

	// Assert
	require.NoError(err)
	require.Equal(platAttendu, platObtenu)
}

func (suite *PréparerTraiflesTestSuite) TestIngrédientManquant() {
	// Arrange
	require := suite.Require()

	recetteFraifles := recette.Recette{
		IngrédientsDeBase: []recette.Ingrédient{
			{
				Nom:      "fraise",
				Quantité: 8,
			},
			{
				Nom:      "trèfle à 4 feuilles",
				Quantité: 13,
			},
		},
		Déroulé: []recette.Étape{
			{
				Nom: "Mélanger les trèfles et les fraises",
				Base: recette.Ingrédient{
					Nom:      "fraise",
					Quantité: 8,
				},
				Avec: recette.Ingrédient{
					Nom:      "trèfle à 4 feuilles",
					Quantité: 13,
				},
				Action:              recette.Mélanger,
				NomIngrédientObtenu: "fraifles",
			},
		},
	}

	fauxStock := stockmocks.NewMockStock(suite.mockController)
	fauxStock.EXPECT().VérifierDisponibilité(gomock.Any()).AnyTimes().Return(false)

	pagoramix := NouveauPagoramix(fauxStock, nil)

	errAttendue := ErrIngrédientManquant

	// Act
	platObtenu, errObtenue := pagoramix.Préparer(recetteFraifles)

	// Assert
	require.ErrorIs(errObtenue, errAttendue)
	require.Nil(platObtenu)
}

func TestPréparerTraifles(t *testing.T) {
	suite.Run(t, &PréparerTraiflesTestSuite{})
}
