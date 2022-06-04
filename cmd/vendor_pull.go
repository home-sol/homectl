package cmd

import (
	"github.com/spf13/cobra"
)

// vendorPullCmd executes 'vendor pull' CLI commands
var vendorPullCmd = &cobra.Command{
	Use:                "pull",
	Short:              "Execute 'vendor pull' commands",
	Long:               `This command executes 'atmos vendor pull' CLI commands`,
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: false},
	RunE: func(cmd *cobra.Command, args []string) error {
		return execVendorCommand(cmd, args, "pull")
	},
}

func init() {
	vendorCmd.AddCommand(vendorPullCmd)
	vendorPullCmd.PersistentFlags().StringP("component", "c", "", "atmos vendor pull --component <component>")
	vendorPullCmd.PersistentFlags().StringP("stack", "s", "", "atmos vendor pull --stack <stack>")
	vendorPullCmd.PersistentFlags().StringP("type", "t", "terraform", "atmos vendor pull --component <component> type=terraform/helmfile")
	vendorPullCmd.PersistentFlags().Bool("dry-run", false, "atmos vendor pull --component <component> --dry-run")
}
