package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "homectl",
	Short: "Universal Tool for Home Automation",
	Long:  `'homectl'' is a universal tool for Home automation used for provisioning, managing and orchestrating deployment`,
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
