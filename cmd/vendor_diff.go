package cmd

import (
	"github.com/spf13/cobra"
)

// vendorDiffCmd executes 'vendor diff' CLI commands
var vendorDiffCmd = &cobra.Command{
	Use:                "diff",
	Short:              "Execute 'vendor diff' commands",
	Long:               `This command executes 'homectl vendor diff' CLI commands`,
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: false},
	RunE: func(cmd *cobra.Command, args []string) error {
		return execVendorCommand(cmd, args, "diff")
	},
}

func init() {
	vendorCmd.AddCommand(vendorDiffCmd)
	vendorDiffCmd.PersistentFlags().StringP("component", "c", "", "homectl vendor diff --component <component>")
	vendorDiffCmd.PersistentFlags().StringP("stack", "s", "", "homectl vendor diff --stack <stack>")
	vendorDiffCmd.PersistentFlags().StringP("type", "t", "terraform", "homectl vendor diff --component <component> --type (terraform|helmfile)")
	vendorDiffCmd.PersistentFlags().Bool("dry-run", false, "homectl vendor diff --component <component> --dry-run")
}
