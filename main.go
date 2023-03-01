package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"

	"github.com/spf13/cobra"
)

const (
	nvmInstallMessage = "You can download the required Node.js version using Node Version Manager (https://github.com/nvm-sh/nvm).\nOnce nvm is installed, run `nvm install` in this directory to install the correct Node version."
)

var mmplugCmd = &cobra.Command{
	Use:   "mmplug",
	Short: "mmplug is a command line tool to manage plugin projects",
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if your development environment is set up correctly",
	Run:   runDoctor,
}

func main() {
	mmplugCmd.AddCommand(doctorCmd)

	if err := mmplugCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDoctor(cmd *cobra.Command, args []string) {
	fmt.Printf("Checking project dependencies\n")

	allPassed := true

	fmt.Println()
	if !runGoCheck() {
		allPassed = false
	}

	fmt.Println()
	if !runNodeCheck() {
		allPassed = false
	}

	fmt.Println()
	if allPassed {
		success("All checks passed.")
		return
	}

	fmt.Println("Not all checks passed. Please fix the above issues and try again.")
}

func runGoCheck() bool {
	goModFile, err := ioutil.ReadFile("go.mod")
	if err != nil {
		success("No go.mod file found. Now assuming the plugin is webapp-side only.")
		return false
	}

	modFile, err := modfile.Parse("go.mod", goModFile, nil)
	if err != nil {
		fail("Failed to parse go.mod file. Error: %v", err)
		return false
	}
	requiredGoVersion := modFile.Go.Version

	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		fail("Failed to run `go version` command. Error: %v", err)
		return false
	}

	// go version go1.19.5 linux/amd64
	words := strings.Split(string(output), " ")
	if len(words) < 3 || !strings.HasPrefix(words[2], "go") {
		fail("Failed to parse installed Go version from `go version` command")
		return false
	}
	currentGoVersion := words[2][2:]

	if semver.Compare("v"+currentGoVersion, "v"+requiredGoVersion) >= 0 {
		success("Go version %s is compatible with required version %s", currentGoVersion, requiredGoVersion)
		return true
	}

	fail("Go version %s is incompatible with required version %s", currentGoVersion, requiredGoVersion)
	fmt.Println("Please follow the instructions at https://go.dev to download the correct version.")
	return false
}

func runNodeCheck() bool {
	nvmrcFile, err := ioutil.ReadFile(".nvmrc")
	if err != nil {
		success("No .nvmrc file found. Now assuming the plugin is server-side only.")
		return true
	}

	requiredNodeVersion := strings.TrimSpace(string(nvmrcFile))
	if !strings.HasPrefix(requiredNodeVersion, "v") {
		requiredNodeVersion = "v" + requiredNodeVersion
	}

	cmd := exec.Command("node", "-v")
	output, err := cmd.Output()
	if err != nil {
		fail("Node.js is not installed or not in PATH. %s", nvmInstallMessage)
		return false
	}
	currentNodeVersion := strings.TrimSpace(string(output))

	if semver.MajorMinor(currentNodeVersion) == semver.MajorMinor(requiredNodeVersion) {
		success("Node version %s is compatible with required version %s", currentNodeVersion, requiredNodeVersion)
		return true
	}

	fail("Node version %s is incompatible with required version %s", currentNodeVersion, requiredNodeVersion)

	cmd = exec.Command("nvm")
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(nvmInstallMessage)
		return false
	}

	fmt.Println("Run `nvm install` in this directory to install the correct version.")

	return false
}

func success(format string, args ...any) {
	formatted := fmt.Sprintf(format, args...)
	fmt.Printf("✅ %s\n", formatted)
}

func fail(format string, args ...any) {
	formatted := fmt.Sprintf(format, args...)
	fmt.Printf("❌ %s\n", formatted)
}
