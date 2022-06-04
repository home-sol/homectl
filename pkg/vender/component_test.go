package vender_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/home-sol/homectl/pkg/config"
	"github.com/home-sol/homectl/pkg/fs"
	"github.com/home-sol/homectl/pkg/logger"
	"github.com/home-sol/homectl/pkg/vender"
)

const (
	workingDir = "../../examples/complete"
)

func TestVenderComponentPullCommand(t *testing.T) {

	logger.Logger = zap.NewNop().Sugar()

	err := config.InitConfigFromDir(workingDir)
	require.NoError(t, err)

	componentType := "terraform"
	vendorCommand := "pull"

	component := "infra/vpc-flow-logs-bucket"

	fss, err := fs.FromDir(workingDir)
	require.NoError(t, err)

	componentConfig, componentPath, err := config.ReadComponentFile(fss, component, componentType)
	assert.Nil(t, err)

	err = vender.ExecuteComponentVendorCommand(fss, componentConfig.Spec, component, componentPath, false, vendorCommand)
	assert.Nil(t, err)

	// Check if the correct files were pulled and written to the correct folder
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "context.tf")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "default.auto.tfvars")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "introspection.mixin.tf")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "main.tf")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "outputs.tf")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "providers.tf")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "README.md")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "variables.tf")))
	assert.FileExists(t, fss.GetRelativePath(path.Join(componentPath, "versions.tf")))

	// Delete the files
	err = fss.Remove(path.Join(componentPath, "context.tf"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "default.auto.tfvars"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "introspection.mixin.tf"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "main.tf"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "outputs.tf"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "providers.tf"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "README.md"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "variables.tf"))
	assert.Nil(t, err)
	err = fss.Remove(path.Join(componentPath, "versions.tf"))
	assert.Nil(t, err)
}
