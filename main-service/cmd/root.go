/*
Copyright © 2022 Nooder
*/
package main

import (
	"context"
	"go-service-skeleton/internal/app"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-service-skeleton",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rootCtx := context.TODO()
		cfg := app.Config{}
		err := app.LoadConfig(rootCtx, "/", &cfg)
		if err != nil {
			return err
		}
		return app.StartService(rootCtx, &cfg)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-service-skeleton.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

func main() {
	Execute()
}
