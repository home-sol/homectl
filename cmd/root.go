package cmd

import (
	"github.com/spf13/cobra"

	"github.com/home-sol/homectl/pkg/config"
	"github.com/home-sol/homectl/pkg/logger"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "homectl",
	Short: "Universal Tool for Home Automation",
	Long:  `'homectl'' is a universal tool for Home automation used for provisioning, managing and orchestrating deployment`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := logger.InitLogger(); err != nil {
			return err
		}
		return config.InitConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
