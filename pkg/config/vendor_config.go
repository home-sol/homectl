package config

import (
	"fmt"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/home-sol/homectl/pkg/fs"
)

// ReadComponentFile reads and processes `component.yaml` vendor config file
func ReadComponentFile(fss *fs.FileSystem, component string, componentType string) (VendorComponentConfig, string, error) {
	var componentBasePath string
	var componentConfig VendorComponentConfig

	if componentType == "terraform" {
		componentBasePath = Config.Components.Terraform.BasePath
	} else if componentType == "helmfile" {
		componentBasePath = Config.Components.Helmfile.BasePath
	} else {
		return componentConfig, "", fmt.Errorf("type '%s' is not supported. Valid types are 'terraform' and 'helmfile'", componentType)
	}

	componentPath := path.Join(Config.BasePath, componentBasePath, component)

	dirExists, err := fss.IsDirectory(componentPath)
	if err != nil {
		return componentConfig, "", err
	}

	if !dirExists {
		return componentConfig, "", fmt.Errorf("folder '%s' does not exist", componentPath)
	}

	componentConfigFile := path.Join(componentPath, "component.yaml")
	if !fss.FileExists(componentConfigFile) {
		return componentConfig, "", fmt.Errorf("vendor config file 'component.yaml' does not exist in the '%s' folder", componentPath)
	}

	componentConfigFileContent, err := fss.ReadFile(componentConfigFile)
	if err != nil {
		return componentConfig, "", err
	}

	if err = yaml.Unmarshal(componentConfigFileContent, &componentConfig); err != nil {
		return componentConfig, "", err
	}

	if componentConfig.Kind != "ComponentVendorConfig" {
		return componentConfig, "", fmt.Errorf("invalid 'kind: %s' in the vendor config file 'component.yaml'. Supported kinds: 'ComponentVendorConfig'",
			componentConfig.Kind,
		)
	}

	return componentConfig, componentPath, nil
}
