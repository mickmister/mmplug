package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	nvmInstallMessage = "You can download the required Node.js version using Node Version Manager (https://github.com/nvm-sh/nvm).\nOnce nvm is installed, run `nvm install` in this directory to install the correct Node version."
)

var mmplugCmd = &cobra.Command{
	Use:          "mmplug",
	Short:        "mmplug is a command line tool to manage plugin projects",
	SilenceUsage: true,
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if your development environment is set up correctly",
	RunE:  runDoctor,
}

func main() {
	mmplugCmd.AddCommand(doctorCmd)

	if err := mmplugCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking project dependencies")

	manifest, err := findManifest()
	if err != nil {
		return errors.Wrap(err, "Failed to find plugin manifest. This is probably not a Mattermost plugin project.")
	}

	if !manifest.HasServer() && !manifest.HasWebapp() {
		return errors.New("This plugin project does not contain server or webapp component in its manifest.")
	}

	allPassed := true

	fmt.Println()
	if manifest.HasServer() {
		fmt.Println("Performing checks for Go")
		if !runGoCheck() {
			allPassed = false
		}
	} else {
		fmt.Println("No server component found in plugin manifest. Now assuming the plugin is webapp-side only.")
	}

	fmt.Println()
	if manifest.HasWebapp() {
		fmt.Println("Performing checks for Node.js")
		if !runNodeCheck() {
			allPassed = false
		}
	} else {
		fmt.Println("No webapp component found in plugin manifest. Now assuming the plugin is server-side only.")
	}

	fmt.Println()
	if allPassed {
		success("All checks passed.")
		return nil
	}

	return errors.New("Not all checks passed. Please fix the above issues and try again.")
}
