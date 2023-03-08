package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func runGoCheck() bool {
	goModFile, err := ioutil.ReadFile("go.mod")
	if err != nil {
		fail("No go.mod file found for server plugin.")
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
		fail("No .nvmrc file found, therefore the required Node.js version is unknown.")
		return false
	}
	requiredNodeVersion := strings.TrimSpace(string(nvmrcFile))

	cmd := exec.Command("node", "-v")
	output, err := cmd.Output()
	if err != nil {
		fail("Node.js is not installed or not in PATH. %s", nvmInstallMessage)
		return false
	}
	currentNodeVersion := strings.TrimSpace(string(output))

	if !strings.HasPrefix(requiredNodeVersion, "v") {
		requiredNodeVersion = "v" + requiredNodeVersion
	}

	if semver.MajorMinor(currentNodeVersion) == semver.MajorMinor(requiredNodeVersion) {
		success("Installed Node.js version %s is compatible with required version %s", currentNodeVersion, requiredNodeVersion)
		return true
	}

	fail("Installed Node.js version %s is incompatible with required version %s", currentNodeVersion, requiredNodeVersion)
	fmt.Println(nvmInstallMessage)
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
