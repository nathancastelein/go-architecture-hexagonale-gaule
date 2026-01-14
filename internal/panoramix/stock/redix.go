package stock

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
)

const stockKey = "stock:ingredients"

type Redix struct {
	client *redis.Client
}

func NouveauRedix(client *redis.Client) Stock {
	return &Redix{
		client: client,
	}
}

// VérifierDisponibilité implements Stock.
func (s *Redix) VérifierDisponibilité(ingrédient recette.Ingrédient) bool {
	ctx := context.Background()

	quantitéStr, err := s.client.HGet(ctx, stockKey, ingrédient.Nom).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		return false
	}

	quantité, err := strconv.Atoi(quantitéStr)
	if err != nil {
		return false
	}

	return quantité >= ingrédient.Quantité
}

// RécupèrerIngrédient implements Stock.
func (s *Redix) RécupèrerIngrédient(ingrédientDemandé recette.Ingrédient) (recette.Ingrédient, error) {
	ctx := context.Background()

	// Use a transaction to ensure atomicity
	txf := func(tx *redis.Tx) error {
		quantitéStr, err := tx.HGet(ctx, stockKey, ingrédientDemandé.Nom).Result()
		if err == redis.Nil {
			return ErrIngrédientNonTrouvé
		}
		if err != nil {
			return err
		}

		quantité, err := strconv.Atoi(quantitéStr)
		if err != nil {
			return err
		}

		if quantité < ingrédientDemandé.Quantité {
			return ErrIngrédientNonTrouvé
		}

		nouvelleQuantité := quantité - ingrédientDemandé.Quantité

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if nouvelleQuantité == 0 {
				pipe.HDel(ctx, stockKey, ingrédientDemandé.Nom)
			} else {
				pipe.HSet(ctx, stockKey, ingrédientDemandé.Nom, strconv.Itoa(nouvelleQuantité))
			}
			return nil
		})

		return err
	}

	err := s.client.Watch(ctx, txf, stockKey)
	if err != nil {
		return recette.Ingrédient{}, err
	}

	return ingrédientDemandé, nil
}

// StockerIngrédient implements Stock.
func (s *Redix) StockerIngrédient(ingrédient recette.Ingrédient) {
	ctx := context.Background()
	s.client.HIncrBy(ctx, stockKey, ingrédient.Nom, int64(ingrédient.Quantité))
}
