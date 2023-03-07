package main

import (
	"encoding/json"
	"os"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

func findManifest() (*model.Manifest, error) {
	_, manifestFilePath, err := model.FindManifest(".")
	if err != nil {
		return nil, errors.Wrap(err, "failed to find manifest in current working directory")
	}
	manifestFile, err := os.Open(manifestFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", manifestFilePath)
	}
	defer manifestFile.Close()

	var manifest model.Manifest
	decoder := json.NewDecoder(manifestFile)
	if err = decoder.Decode(&manifest); err != nil {
		return nil, errors.Wrap(err, "failed to parse manifest")
	}

	return &manifest, nil
}
