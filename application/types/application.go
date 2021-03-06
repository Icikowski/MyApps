package types

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	"icikowski.pl/myapps/common"
)

// Application represents the definition of an external application.
type Application struct {
	Name                string   `json:"name" yaml:"name"`
	Description         string   `json:"description" yaml:"description"`
	InstallScenario     []string `json:"install_scenario" yaml:"install_scenario"`
	UpdateScenario      []string `json:"update_scenario" yaml:"update_scenario"`
	UninstallScenario   []string `json:"uninstall_scenario" yaml:"uninstall_scenario"`
	NewestVersionCheck  []string `json:"newest_version_check" yaml:"newest_version_check"`
	CurrentVersionCheck []string `json:"current_version_check" yaml:"current_version_check"`
}

// executeStep runs single step of any scenario and returns error that contains
// information that user needs to debug issue with the scenario
func executeStep(index int, step string, env []string) (string, error) {
	cmd := exec.Command("sh", "-c", step)
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf(
			"step %d failed\n\ncmd: %s\nerror: %w\noutput:\n%s",
			index, step, err, string(out),
		)
	}
	return string(out), nil
}

// GetCurrentVersion returns the current version of the application.
func (app Application) GetCurrentVersion() (*version.Version, error) {
	var output string
	for i, step := range app.CurrentVersionCheck {
		rawOutput, err := executeStep(i, step, []string{})
		if err != nil {
			return nil, err
		}
		output = strings.TrimSpace(rawOutput)
	}
	return version.NewVersion(output)
}

// GetLatestVersion returns the latest version of the application.
func (app Application) GetLatestVersion() (*version.Version, error) {
	var output string
	for i, step := range app.NewestVersionCheck {
		rawOutput, err := executeStep(i, step, []string{})
		if err != nil {
			return nil, err
		}
		output = strings.TrimSpace(rawOutput)
	}
	return version.NewVersion(output)
}

// IsUpdateAvailable returns true if there is an update available for the
// application.
func (app Application) IsUpdateAvailable() (bool, error) {
	currentVersion, err := app.GetCurrentVersion()
	if err != nil {
		return false, err
	}

	latestVersion, err := app.GetLatestVersion()
	if err != nil {
		return false, err
	}

	return currentVersion.LessThan(latestVersion), nil
}

// Install installs the application by executing the install scenario.
func (app Application) Install() error {
	bar := common.NewProgressBar(len(app.InstallScenario) + 3)
	defer bar.Finish()

	bar.Increment()
	latestVersion, err := app.GetLatestVersion()
	if err != nil {
		return err
	}

	bar.Increment()
	tmpDir, cleanup, err := common.GetTempDir()
	if err != nil {
		return err
	}
	defer cleanup()

	bar.Increment()
	commandEnvironment := append(
		os.Environ(),
		"LATEST_VERSION="+latestVersion.String(),
		"TEMP="+tmpDir,
		"OS="+runtime.GOOS,
		"ARCH="+runtime.GOARCH,
	)

	for i, step := range app.InstallScenario {
		bar.Increment()
		_, err := executeStep(i, step, commandEnvironment)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update updates the application by executing the update scenario.
func (app Application) Update() error {
	bar := common.NewProgressBar(len(app.UpdateScenario) + 3)
	defer bar.Finish()

	bar.Increment()
	currentVersion, err := app.GetLatestVersion()
	if err != nil {
		return err
	}

	bar.Increment()
	latestVersion, err := app.GetLatestVersion()
	if err != nil {
		return err
	}

	bar.Increment()
	tmpDir, cleanup, err := common.GetTempDir()
	if err != nil {
		return err
	}
	defer cleanup()

	bar.Increment()
	commandEnvironment := append(
		os.Environ(),
		"CURRENT_VERSION="+currentVersion.String(),
		"LATEST_VERSION="+latestVersion.String(),
		"TEMP="+tmpDir,
		"OS="+runtime.GOOS,
		"ARCH="+runtime.GOARCH,
	)

	for i, step := range app.UpdateScenario {
		bar.Increment()
		_, err := executeStep(i, step, commandEnvironment)
		if err != nil {
			return err
		}
	}
	return nil
}

// Uninstall removes the application by executing the uninstall scenario.
func (app Application) Uninstall() error {
	bar := common.NewProgressBar(len(app.UninstallScenario) + 3)
	defer bar.Finish()

	bar.Increment()
	currentVersion, err := app.GetLatestVersion()
	if err != nil {
		return err
	}

	bar.Increment()
	commandEnvironment := append(
		os.Environ(),
		"CURRENT_VERSION="+currentVersion.String(),
		"OS="+runtime.GOOS,
		"ARCH="+runtime.GOARCH,
	)

	for i, step := range app.UninstallScenario {
		bar.Increment()
		_, err := executeStep(i, step, commandEnvironment)
		if err != nil {
			return err
		}
	}
	return nil
}

// Applications represents a list of applications.
type Applications []Application

// FindByName returns the application that matches given name.
func (apps Applications) FindByName(name string) (Application, bool) {
	for _, app := range apps {
		if app.Name == name {
			return app, true
		}
	}
	return Application{}, false
}

// FindByNameLike returns the applications which name contains given name.
func (apps Applications) FindByNameLike(name string) (Applications, bool) {
	var result Applications
	for _, app := range apps {
		if strings.Contains(app.Name, name) {
			result = append(result, app)
		}
	}
	return result, len(result) != 0
}
