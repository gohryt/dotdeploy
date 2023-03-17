package main

import (
	"github.com/gohryt/dotdeploy"
)

type (
	Deploy struct {
		Folder string `yaml:"folder"`
		Keep   bool   `yaml:"keep"`

		Remote []Connection `yaml:"Remote"`
		Do     []Action     `yaml:"Do"`
	}

	File struct {
		Path string `yaml:"path"`
	}

	Path struct {
		Connection string `yaml:"connection"`
		Path       string `yaml:"path"`
	}
)

func (deploy *Deploy) Deploy() *dotdeploy.Deploy {
	new := &dotdeploy.Deploy{
		Folder: deploy.Folder,
		Keep:   deploy.Keep,

		Remote: make(dotdeploy.Remote, len(deploy.Remote)),
		Do:     make(dotdeploy.Do, len(deploy.Do)),
	}

	for i, connection := range deploy.Remote {
		new.Remote[i] = connection.Connection()
	}

	for i, action := range deploy.Do {
		new.Do[i] = action.Action()
	}

	return new
}
