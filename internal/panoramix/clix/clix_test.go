package clix

import (
	"bytes"
	"errors"
	"testing"

	druidemocks "github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/druide/mocks"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CLIxTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
}

func (suite *CLIxTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
}

func (suite *CLIxTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func TestCLIx(t *testing.T) {
	suite.Run(t, &CLIxTestSuite{})
}

func (suite *CLIxTestSuite) TestNouveauCLIx_ContientLaCommandePotionMagique() {
	// Arrange
	require := suite.Require()
	fauxDruide := druidemocks.NewMockDruide(suite.mockController)

	// Act
	cli := NouveauCLIx(fauxDruide)

	// Assert
	commands := cli.rootCmd.Commands()
	require.Len(commands, 1)
	require.Equal("potion-magique", commands[0].Use)
}

func (suite *CLIxTestSuite) TestExécuter_PotionMagique_Succès() {
	// Arrange
	require := suite.Require()
	fauxDruide := druidemocks.NewMockDruide(suite.mockController)

	platAttendu := &recette.Plat{
		Nom:      "potion magique",
		Quantité: 28,
	}
	fauxDruide.EXPECT().
		Préparer(recette.PotionMagique).
		Return(platAttendu, nil)

	cli := NouveauCLIx(fauxDruide)

	stdout := &bytes.Buffer{}
	cli.rootCmd.SetOut(stdout)
	cli.rootCmd.SetArgs([]string{"potion-magique"})

	// Act
	err := cli.Exécuter()

	// Assert
	require.NoError(err)
	output := stdout.String()
	require.Contains(output, "Préparation réussie!")
	require.Contains(output, "potion magique")
	require.Contains(output, "28")
}

func (suite *CLIxTestSuite) TestExécuter_PotionMagique_Échec() {
	// Arrange
	require := suite.Require()
	fauxDruide := druidemocks.NewMockDruide(suite.mockController)

	errPréparation := errors.New("ingrédient manquant: gui")
	fauxDruide.EXPECT().
		Préparer(recette.PotionMagique).
		Return(nil, errPréparation)

	cli := NouveauCLIx(fauxDruide)

	stderr := &bytes.Buffer{}
	cli.rootCmd.SetErr(stderr)
	cli.rootCmd.SetArgs([]string{"potion-magique"})

	// Act
	err := cli.Exécuter()

	// Assert
	require.Error(err)
	require.ErrorIs(err, errPréparation)
	errOutput := stderr.String()
	require.Contains(errOutput, "Échec de la préparation")
}

func (suite *CLIxTestSuite) TestExécuter_CommandeInconnue() {
	// Arrange
	require := suite.Require()
	fauxDruide := druidemocks.NewMockDruide(suite.mockController)

	cli := NouveauCLIx(fauxDruide)

	stderr := &bytes.Buffer{}
	cli.rootCmd.SetErr(stderr)
	cli.rootCmd.SetArgs([]string{"commande-inexistante"})

	// Act
	err := cli.Exécuter()

	// Assert
	require.Error(err)
	require.Contains(err.Error(), "unknown command")
}

func (suite *CLIxTestSuite) TestExécuter_SansArgument_AfficheLAide() {
	// Arrange
	require := suite.Require()
	fauxDruide := druidemocks.NewMockDruide(suite.mockController)

	cli := NouveauCLIx(fauxDruide)

	stdout := &bytes.Buffer{}
	cli.rootCmd.SetOut(stdout)
	cli.rootCmd.SetArgs([]string{})

	// Act
	err := cli.Exécuter()

	// Assert
	require.NoError(err)
	output := stdout.String()
	require.Contains(output, "Druide permet de préparer des potions magiques")
}

func (suite *CLIxTestSuite) TestPréparerRecette_AfficheLeNomDeLaRecette() {
	// Arrange
	require := suite.Require()
	fauxDruide := druidemocks.NewMockDruide(suite.mockController)

	platAttendu := &recette.Plat{
		Nom:      "potion magique",
		Quantité: 28,
	}
	fauxDruide.EXPECT().
		Préparer(recette.PotionMagique).
		Return(platAttendu, nil)

	cli := NouveauCLIx(fauxDruide)

	stdout := &bytes.Buffer{}
	cli.rootCmd.SetOut(stdout)
	cli.rootCmd.SetArgs([]string{"potion-magique"})

	// Act
	err := cli.Exécuter()

	// Assert
	require.NoError(err)
	output := stdout.String()
	require.Contains(output, "Préparation: Potion Magique")
}
