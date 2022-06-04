package config

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/home-sol/homectl/pkg/logger"
	"github.com/home-sol/homectl/pkg/utils"
)

var (
	// Default values
	defaultConfig = Configuration{
		BasePath: "",
		Components: Components{
			Terraform: Terraform{
				BasePath:                "components/terraform",
				ApplyAutoApprove:        false,
				DeployRunInit:           true,
				InitRunReconfigure:      true,
				AutoGenerateBackendFile: false,
			},
			Helmfile: Helmfile{
				BasePath:              "components/helmfile",
				KubeconfigPath:        "/dev/shm",
				HelmAwsProfilePattern: "{namespace}-{tenant}-gbl-{stage}-helm",
				ClusterNamePattern:    "{namespace}-{tenant}-{environment}-{stage}-eks-cluster",
			},
		},
		Stacks: Stacks{
			BasePath: "stacks",
			IncludedPaths: []string{
				"**/*",
			},
			ExcludedPaths: []string{
				"globals/**/*",
				"catalog/**/*",
				"**/*globals*",
			},
		},
		Workflows: Workflows{
			BasePath: "workflows",
		},
		Logs: Logs{
			Verbose: false,
			Colors:  true,
		},
	}

	// Config is the CLI configuration structure
	Config Configuration
)

func InitConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	return InitConfigFromDir(wd)
}

// Init config finds and merges CLI configuration in the followinfg order system dir, home dir current dir, ENV vars, command-line args
// https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
func InitConfigFromDir(dir string) error {
	logger.Logger.Debugw("Processing and merging configurations in the following order:")
	logger.Logger.Debugw("system dir, home dir, current dir, ENV vars, command-line arguments")

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetTypeByDefaultValue(true)

	// Add default config
	j, err := json.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(j)
	err = v.MergeConfig(reader)
	if err != nil {
		return err
	}

	// Process config in system folder
	var configFileDirs []string

	// https://pureinfotech.com/list-environment-variables-windows-10/
	// https://docs.microsoft.com/en-us/windows/deployment/usmt/usmt-recognized-environment-variables
	// https://softwareengineering.stackexchange.com/questions/299869/where-is-the-appropriate-place-to-put-application-configuration-files-for-each-p
	// https://stackoverflow.com/questions/37946282/why-does-appdata-in-windows-7-seemingly-points-to-wrong-folder
	if runtime.GOOS == "windows" {
		appDataDir := os.Getenv("LOCALAPPDATA")
		if len(appDataDir) > 0 {
			configFileDirs = append(configFileDirs, appDataDir)
		}
	} else {
		configFileDirs = append(configFileDirs, "/usr/local/etc/homectl")
	}

	// Process config in user's HOME dir
	hd, err := homedir.Dir()
	if err != nil {
		return err
	}
	configFileDirs = append(configFileDirs, path.Join(hd, ".homectl"))

	// Process config in the specified dir
	configFileDirs = append(configFileDirs, dir)

	for _, dir := range configFileDirs {
		configFile := path.Join(dir, "homectl.yaml")
		err = processConfigFile(configFile, v)
		if err != nil {
			return err
		}
	}

	// https://gist.github.com/chazcheadle/45bf85b793dea2b71bd05ebaa3c28644
	// https://sagikazarmark.hu/blog/decoding-custom-formats-with-viper/
	err = v.Unmarshal(&Config)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/NCAR/go-figure
// https://github.com/spf13/viper/issues/181
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
func processConfigFile(path string, v *viper.Viper) error {
	l := logger.Logger.With("config", path)
	if !utils.FileExists(path) {
		l.Debug("No CLI config found")
		return nil
	}

	l.Debug("Found CLI config")

	reader, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			l.Error("Error closing file", err)
		}
	}(reader)

	err = v.MergeConfig(reader)
	if err != nil {
		return err
	}

	l.Debug("Processed CLI config")

	return nil
}
