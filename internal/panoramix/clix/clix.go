package clix

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/druide"
	"github.com/nathancastelein/go-architecture-hexagonale-gaule/internal/panoramix/recette"
)

type CLIx struct {
	druide  druide.Druide
	rootCmd *cobra.Command
}

func NouveauCLIx(druide druide.Druide) *CLIx {
	cli := &CLIx{
		druide: druide,
	}

	cli.rootCmd = &cobra.Command{
		Use:   "druide",
		Short: "Druide est un outil de préparation de potions magiques",
		Long:  "Druide permet de préparer des potions magiques en utilisant l'architecture hexagonale.",
	}

	cli.rootCmd.AddCommand(cli.créerCommandePotionMagique())

	return cli
}

func (cli *CLIx) Exécuter() error {
	return cli.rootCmd.Execute()
}

func (cli *CLIx) créerCommandePotionMagique() *cobra.Command {
	return &cobra.Command{
		Use:   "potion-magique",
		Short: "Prépare la potion magique d'Astérix",
		Long:  "Prépare la célèbre potion magique qui donne une force surhumaine.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.préparerRecette(cmd, recette.PotionMagique)
		},
	}
}

func (cli *CLIx) préparerRecette(cmd *cobra.Command, recetteÀPréparer recette.Recette) error {
	fmt.Fprintf(cmd.OutOrStdout(), "Préparation: %s\n", recetteÀPréparer.Nom)

	plat, err := cli.druide.Préparer(recetteÀPréparer)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Échec de la préparation: %v\n", err)
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Préparation réussie!\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Plat obtenu: %s (quantité: %d)\n", plat.Nom, plat.Quantité)

	return nil
}
