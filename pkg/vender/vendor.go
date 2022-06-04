package vender

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/go-getter"
	"github.com/otiai10/copy"

	"github.com/home-sol/homectl/pkg/config"
	"github.com/home-sol/homectl/pkg/fs"
	"github.com/home-sol/homectl/pkg/logger"
	"github.com/home-sol/homectl/pkg/utils"
)

// executeComponentVendorCommandInternal executes a component vendor command
// Supports all protocols (local files, Git, Mercurial, HTTP, HTTPS, Amazon S3, Google GCP),
// URL and archive formats described in https://github.com/hashicorp/go-getter
// https://www.allee.xyz/en/posts/getting-started-with-go-getter
// https://github.com/otiai10/copy
func ExecuteComponentVendorCommand(
	fss *fs.FileSystem,
	vendorComponentSpec config.VendorComponentSpec,
	component string,
	componentPath string,
	dryRun bool,
	vendorCommand string,
) error {

	l := logger.Logger.With("component", component, "componentPath", componentPath)

	var tempDir string
	var err error
	var t *template.Template
	var uri string

	if vendorCommand == "pull" {
		if vendorComponentSpec.Source.Uri == "" {
			return errors.New("'uri' must be specified in 'source.uri' in the 'component.yaml' file")
		}

		// Parse 'uri' template
		if vendorComponentSpec.Source.Version != "" {
			t, err = template.New(fmt.Sprintf("source-uri-%s", vendorComponentSpec.Source.Version)).Parse(vendorComponentSpec.Source.Uri)
			if err != nil {
				return err
			}

			var tpl bytes.Buffer
			err = t.Execute(&tpl, vendorComponentSpec.Source)
			if err != nil {
				return err
			}

			uri = tpl.String()
		} else {
			uri = vendorComponentSpec.Source.Uri
		}

		l.Infof("Pulling sources for the component from '%s'", uri)

		if !dryRun {
			// Create temp folder
			// We are using a temp folder for the following reasons:
			// 1. 'git' does not clone into an existing folder (and we have the existing component folder with `component.yaml` in it)
			// 2. We have the option to skip some files we don't need and include only the files we need when copying from the temp folder to the destination folder
			tempDir, err = ioutil.TempDir("", strconv.FormatInt(time.Now().Unix(), 10))
			if err != nil {
				return err
			}

			defer func(path string) {
				err := os.RemoveAll(path)
				if err != nil {
					l.Error(err)
				}
			}(tempDir)

			// Download the source into the temp folder
			client := &getter.Client{
				Ctx: context.Background(),
				// Define the destination to where the files will be stored. This will create the directory if it doesn't exist
				Dst: tempDir,
				Dir: true,
				// Source
				Src:  uri,
				Mode: getter.ClientModeDir,
			}

			if err = client.Get(); err != nil {
				return err
			}

			// Copy from the temp folder to the destination folder with skipping of some files
			copyOptions := copy.Options{
				// Skip specifies which files should be skipped
				Skip: func(src string) (bool, error) {
					if strings.HasSuffix(src, ".git") {
						return true, nil
					}

					trimmedSrc := utils.TrimBasePathFromPath(tempDir+"/", src)

					// Exclude the files that match the 'excluded_paths' patterns
					// It supports POSIX-style Globs for file names/paths (double-star `**` is supported)
					// https://en.wikipedia.org/wiki/Glob_(programming)
					// https://github.com/bmatcuk/doublestar#patterns
					for _, excludePath := range vendorComponentSpec.Source.ExcludedPaths {
						excludeMatch, err := doublestar.PathMatch(excludePath, src)
						if err != nil {
							return true, err
						} else if excludeMatch {
							// If the file matches ANY of the 'excluded_paths' patterns, exclude the file
							l.Infow("Excluding the file", "src", trimmedSrc, "pattern", excludePath)
							return true, nil
						}
					}

					// Only include the files that match the 'included_paths' patterns (if any pattern is specified)
					if len(vendorComponentSpec.Source.IncludedPaths) > 0 {
						for _, includePath := range vendorComponentSpec.Source.IncludedPaths {
							includeMatch, err := doublestar.PathMatch(includePath, src)
							if err != nil {
								return true, err
							} else if includeMatch {
								// If the file matches ANY of the 'included_paths' patterns, include the file
								l.Infow("Including the file", "src", trimmedSrc, "pattern", includePath)
								return false, nil
							}
						}

						l.Infof("Excluding since it does not match any pattern from 'included_paths'", "src", trimmedSrc)
						return true, nil
					}

					// If 'included_paths' is not provided, include all files that were not excluded
					l.Infof("Including the file", "src", trimmedSrc)
					return false, nil
				},

				// Preserve the atime and the mtime of the entries
				// On linux we can preserve only up to 1 millisecond accuracy
				PreserveTimes: false,

				// Preserve the uid and the gid of all entries
				PreserveOwner: false,
			}

			if err = copy.Copy(tempDir, fss.GetRelativePath(componentPath), copyOptions); err != nil {
				return err
			}
		}

		// Process mixins
		if len(vendorComponentSpec.Mixins) > 0 {
			for _, mixin := range vendorComponentSpec.Mixins {
				if mixin.Uri == "" {
					return errors.New("'uri' must be specified for each 'mixin' in the 'component.yaml' file")
				}

				if mixin.Filename == "" {
					return errors.New("'filename' must be specified for each 'mixin' in the 'component.yaml' file")
				}

				// Parse 'uri' template
				if mixin.Version != "" {
					t, err = template.New(fmt.Sprintf("mixin-uri-%s", mixin.Version)).Parse(mixin.Uri)
					if err != nil {
						return err
					}

					var tpl bytes.Buffer
					err = t.Execute(&tpl, mixin)
					if err != nil {
						return err
					}

					uri = tpl.String()
				} else {
					uri = mixin.Uri
				}

				l.With("componentPath", path.Join(componentPath, mixin.Filename)).Infof("Pulling the mixin '%s'", uri)

				if !dryRun {
					err = os.RemoveAll(tempDir)
					if err != nil {
						return err
					}

					// Download the mixin into the temp file
					client := &getter.Client{
						Ctx:  context.Background(),
						Dst:  path.Join(tempDir, mixin.Filename),
						Dir:  false,
						Src:  uri,
						Mode: getter.ClientModeFile,
					}

					if err = client.Get(); err != nil {
						return err
					}

					// Copy from the temp folder to the destination folder
					copyOptions := copy.Options{
						// Preserve the atime and the mtime of the entries
						PreserveTimes: false,

						// Preserve the uid and the gid of all entries
						PreserveOwner: false,
					}

					if err = copy.Copy(tempDir, fss.GetRelativePath(componentPath), copyOptions); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// executeStackVendorCommandInternal executes a stack vendor command
// TODO: implement this
func ExecuteStackVendorCommand(
	fss *fs.FileSystem,
	stack string,
	dryRun bool,
	vendorCommand string,
) error {
	return fmt.Errorf("command 'homectl vendor %s --stack <stack>' is not implemented yet", vendorCommand)
}
