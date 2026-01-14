package stock

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
)

type RedixTestSuite struct {
	suite.Suite
	ctx            context.Context
	redisContainer *tcredis.RedisContainer
	redisClient    *redis.Client
}

func (suite *RedixTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	redisContainer, err := tcredis.Run(suite.ctx, "redis:8-alpine")
	if err != nil {
		suite.T().Skipf("skipping Redix tests: Fail to run Docker: %v", err)
	}
	suite.redisContainer = redisContainer

	connectionString, err := redisContainer.ConnectionString(suite.ctx)
	suite.Require().NoError(err)

	opts, err := redis.ParseURL(connectionString)
	suite.Require().NoError(err)

	suite.redisClient = redis.NewClient(opts)
}

func (suite *RedixTestSuite) TearDownSuite() {
	if suite.redisClient != nil {
		suite.redisClient.Close()
	}
	if suite.redisContainer != nil {
		suite.redisContainer.Terminate(suite.ctx)
	}
}

func (suite *RedixTestSuite) SetupTest() {
	// Clean Redis before each test
	suite.redisClient.FlushAll(suite.ctx)
}

func TestRedix(t *testing.T) {
	suite.Run(t, &RedixTestSuite{})
}

func (suite *RedixTestSuite) TestVérifierDisponibilité_IngrédientEnStock() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	suite.redisClient.HSet(suite.ctx, stockKey, "sanglier", "3")
	suite.redisClient.HSet(suite.ctx, stockKey, "gui", "15")

	// Act
	ingrédientDisponible := redix.VérifierDisponibilité(recette.Ingrédient{
		Nom:      "sanglier",
		Quantité: 2,
	})

	// Assert
	require.True(ingrédientDisponible)
}

func (suite *RedixTestSuite) TestVérifierDisponibilité_IngrédientManquant() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	suite.redisClient.HSet(suite.ctx, stockKey, "sanglier", "3")
	suite.redisClient.HSet(suite.ctx, stockKey, "gui", "15")

	// Act
	ingrédientDisponible := redix.VérifierDisponibilité(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	require.False(ingrédientDisponible)
}

func (suite *RedixTestSuite) TestVérifierDisponibilité_QuantitéInsuffisante() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	suite.redisClient.HSet(suite.ctx, stockKey, "sanglier", "3")

	// Act
	ingrédientDisponible := redix.VérifierDisponibilité(recette.Ingrédient{
		Nom:      "sanglier",
		Quantité: 5,
	})

	// Assert
	require.False(ingrédientDisponible)
}

func (suite *RedixTestSuite) TestRécupérerIngrédient_IngrédientEnStock() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	suite.redisClient.HSet(suite.ctx, stockKey, "sanglier", "3")
	suite.redisClient.HSet(suite.ctx, stockKey, "gui", "15")

	// Act
	ingrédientObtenu, err := redix.RécupèrerIngrédient(recette.Ingrédient{
		Nom:      "sanglier",
		Quantité: 2,
	})

	// Assert
	require.NoError(err)
	require.Equal(recette.Ingrédient{Nom: "sanglier", Quantité: 2}, ingrédientObtenu)

	quantitéRestante, err := suite.redisClient.HGet(suite.ctx, stockKey, "sanglier").Result()
	require.NoError(err)
	require.Equal("1", quantitéRestante)
}

func (suite *RedixTestSuite) TestRécupérerIngrédient_IngrédientEnStock_ToutRetirer() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	suite.redisClient.HSet(suite.ctx, stockKey, "sanglier", "3")
	suite.redisClient.HSet(suite.ctx, stockKey, "gui", "15")

	// Act
	ingrédientObtenu, err := redix.RécupèrerIngrédient(recette.Ingrédient{
		Nom:      "sanglier",
		Quantité: 3,
	})

	// Assert
	require.NoError(err)
	require.Equal(recette.Ingrédient{Nom: "sanglier", Quantité: 3}, ingrédientObtenu)

	// Verify ingredient is removed from Redis
	exists, err := suite.redisClient.HExists(suite.ctx, stockKey, "sanglier").Result()
	require.NoError(err)
	require.False(exists)

	// Verify gui is still there
	guiQuantité, err := suite.redisClient.HGet(suite.ctx, stockKey, "gui").Result()
	require.NoError(err)
	require.Equal("15", guiQuantité)
}

func (suite *RedixTestSuite) TestRécupérerIngrédient_IngrédientManquant() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	suite.redisClient.HSet(suite.ctx, stockKey, "sanglier", "3")

	// Act
	ingrédientObtenu, err := redix.RécupèrerIngrédient(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	require.ErrorIs(err, ErrIngrédientNonTrouvé)
	require.Equal(recette.Ingrédient{}, ingrédientObtenu)
}

func (suite *RedixTestSuite) TestStockerIngrédient_StockVide() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	// Act
	redix.StockerIngrédient(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	quantité, err := suite.redisClient.HGet(suite.ctx, stockKey, "cervoise").Result()
	require.NoError(err)
	require.Equal("5", quantité)
}

func (suite *RedixTestSuite) TestStockerIngrédient_StockDéjàRempli() {
	// Arrange
	require := suite.Require()
	redix := NouveauRedix(suite.redisClient)

	suite.redisClient.HSet(suite.ctx, stockKey, "cervoise", "5")

	// Act
	redix.StockerIngrédient(recette.Ingrédient{
		Nom:      "cervoise",
		Quantité: 5,
	})

	// Assert
	quantité, err := suite.redisClient.HGet(suite.ctx, stockKey, "cervoise").Result()
	require.NoError(err)
	require.Equal("10", quantité)
}
