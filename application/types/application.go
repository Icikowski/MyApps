package types

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
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

// GetCurrentVersion returns the current version of the application.
func (app Application) GetCurrentVersion() (*version.Version, error) {
	var output string
	for _, step := range app.CurrentVersionCheck {
		cmd := exec.Command("sh", "-c", step)
		rawOutput, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		output = strings.TrimSpace(string(rawOutput))
	}
	return version.NewVersion(output)
}

// runStep runs single step of any scenario and returns error that contains
// information that user needs to debug issue with the scenario
func (app Application) runStep(index int, step string, env []string) (string, error) {
	cmd := exec.Command("sh", "-c", step)
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf(
			"step %d failed\n\ncmd: %s\nerror: %w\noutput:\n%s\n",
			index, step, err, string(out),
		)
	}
	return string(out), nil
}

// GetLatestVersion returns the latest version of the application.
func (app Application) GetLatestVersion() (*version.Version, error) {
	var output string
	for i, step := range app.NewestVersionCheck {
		rawOutput, err := app.runStep(i, step, []string{})
		if err != nil {
			return nil, err
		}
		output = strings.TrimSpace(rawOutput)
	}
	return version.NewVersion(output)
}

// Install installs the application by executing the install scenario.
func (app Application) Install() error {
	bar := progressBar.Start(len(app.InstallScenario) + 2)

	bar.Increment()
	latestVersion, err := app.GetLatestVersion()
	if err != nil {
		return err
	}

	bar.Increment()
	commandEnvironment := append(
		os.Environ(),
		"LATEST_VERSION="+latestVersion.String(),
	)

	for i, step := range app.InstallScenario {
		bar.Increment()
		_, err := app.runStep(i, step, commandEnvironment)
		if err != nil {
			return err
		}
	}
	finishProgressBar(bar)
	return nil
}

// Update updates the application by executing the update scenario.
func (app Application) Update() error {
	bar := progressBar.Start(len(app.UpdateScenario) + 3)

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
	commandEnvironment := append(
		os.Environ(),
		"CURRENT_VERSION="+currentVersion.String(),
		"LATEST_VERSION="+latestVersion.String(),
	)

	for i, step := range app.UpdateScenario {
		bar.Increment()
		_, err := app.runStep(i, step, commandEnvironment)
		if err != nil {
			return err
		}
	}
	finishProgressBar(bar)
	return nil
}

// Uninstall removes the application by executing the uninstall scenario.
func (app Application) Uninstall() error {
	bar := progressBar.Start(len(app.UninstallScenario) + 2)

	bar.Increment()
	currentVersion, err := app.GetLatestVersion()
	if err != nil {
		return err
	}

	bar.Increment()
	commandEnvironment := append(
		os.Environ(),
		"CURRENT_VERSION="+currentVersion.String(),
	)

	for i, step := range app.UninstallScenario {
		bar.Increment()
		_, err := app.runStep(i, step, commandEnvironment)
		if err != nil {
			return err
		}
	}
	finishProgressBar(bar)
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
