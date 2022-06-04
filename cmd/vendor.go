package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/home-sol/homectl/pkg/config"
	"github.com/home-sol/homectl/pkg/fs"
	"github.com/home-sol/homectl/pkg/vender"
)

// vendorCmd executes 'atmos vendor' CLI commands
var vendorCmd = &cobra.Command{
	Use:                "vendor",
	Short:              "Execute 'vendor' commands",
	Long:               `This command executes 'homectl vendor' CLI commands`,
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: false},
}

func init() {
	RootCmd.AddCommand(vendorCmd)
}

func execVendorCommand(cmd *cobra.Command, args []string, vendorCommand string) error {

	flags := cmd.Flags()

	dryRun, err := flags.GetBool("dry-run")
	if err != nil {
		return err
	}

	component, err := flags.GetString("component")
	if err != nil {
		return err
	}

	stack, err := flags.GetString("stack")
	if err != nil {
		return err
	}

	if component != "" && stack != "" {
		return errors.New("either '--component' or '--stack' parameter needs to be provided, but not both")
	}

	fss, err := fs.Cwd()
	if err != nil {
		return err
	}

	if component != "" {
		// Process component vendoring
		componentType, err := flags.GetString("type")
		if err != nil {
			return err
		}

		if componentType == "" {
			componentType = "terraform"
		}

		componentConfig, componentPath, err := config.ReadComponentFile(fss, component, componentType)
		if err != nil {
			return err
		}

		return vender.ExecuteComponentVendorCommand(fss, componentConfig.Spec, component, componentPath, dryRun, vendorCommand)
	} else {
		// Process stack vendoring
		return vender.ExecuteStackVendorCommand(fss, stack, dryRun, vendorCommand)
	}
}
