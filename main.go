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
	nvmInstallMessage = "You can download the required Node.js version using Node Version Manager (https://github.com/nvm-sh/nvm). Once nvm is installed, run `nvm install` in this directory to install the correct Node version."
)

var mmplugCmd = &cobra.Command{
	Use:   "mmplug",
	Short: "mmplug is a command line tool to manage plugin projects",
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if your development environment set up correctly",
	Run:   runDoctor,
}

func main() {
	mmplugCmd.AddCommand(doctorCmd)

	if err := mmplugCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
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

	if !allPassed {
		fail("Not all checks passed. Please fix the above issues and try again.")
	} else {
		success("All checks passed.")
	}
}

func runGoCheck() bool {
	goModFile, err := ioutil.ReadFile("go.mod")
	if err != nil {
		fail("Failed to read go.mod file")
		return false
	}

	modFile, err := modfile.Parse("go.mod", goModFile, nil)
	if err != nil {
		fail("Failed to parse go.mod file")
		return false
	}

	requiredGoVersion := "v" + modFile.Go.Version

	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		fail("Failed to run go version command")
		return false
	}

	// go version go1.19.5 linux/amd64
	currentGoVersion := "v" + strings.Split(string(output), " ")[2][2:]

	if semver.Compare(currentGoVersion, requiredGoVersion) < 0 {
		fail("Go version %s is incompatible with required version %s", currentGoVersion, requiredGoVersion)
		fmt.Println("Please follow the instructions at https://go.dev to download the correct version.")
		return false
	}

	success("Go version %s is compatible with required version %s", currentGoVersion, requiredGoVersion)
	return true
}

func runNodeCheck() bool {
	nvmrcFile, err := ioutil.ReadFile(".nvmrc")
	if err != nil {
		fail("Failed to read .nvmrc file")
		return false
	}

	requiredNodeVersion := strings.TrimSpace(string(nvmrcFile))
	if !strings.HasPrefix(requiredNodeVersion, "v") {
		requiredNodeVersion = "v" + requiredNodeVersion
	}

	cmd := exec.Command("node", "-v")
	output, err := cmd.Output()
	if err != nil {
		fail("Node.js is not installed or not in PATH. %s", nvmInstallMessage)
	}

	currentNodeVersion := strings.TrimSpace(string(output))

	if semver.Major(currentNodeVersion) == semver.Major(requiredNodeVersion) {
		success("Node version %s is compatible with required version %s", currentNodeVersion, requiredNodeVersion)
		return true
	}

	fail("Node version %s is incompatible with required version %s", currentNodeVersion, requiredNodeVersion)

	_, err = exec.LookPath("nvm")
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