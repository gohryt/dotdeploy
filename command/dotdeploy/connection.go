package main

import (
	"gopkg.in/yaml.v3"

	"github.com/gohryt/dotdeploy"
	"github.com/gohryt/dotdeploy/unsafe"
)

type (
	Connection struct {
		Name string
		Data any
	}

	ConnectionType struct {
		Type string `yaml:"type"`
		Name string `yaml:"name"`
	}

	Key struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		File     string `yaml:"file"`
		Password string `yaml:"password"`
	}

	Password struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}

	Agent struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
	}
)

func (connection *Connection) UnmarshalYAML(value *yaml.Node) error {
	t := new(ConnectionType)

	err := value.Decode(t)
	if err != nil {
		return err
	}

	if t.Name != "" {
		connection.Name = t.Name
	} else {
		connection.Name = t.Type
	}

	mask := unsafe.As[unsafe.Any](&connection.Data)

	switch t.Type {
	case "key":
		connection.Data = new(Key)

		err = value.Decode(connection.Data)
		mask.Type = dotdeploy.ConnectionTypeKey
	case "password":
		connection.Data = new(Password)

		err = value.Decode(connection.Data)
		mask.Type = dotdeploy.ConnectionTypePassword
	case "agent":
		connection.Data = new(Agent)

		err = value.Decode(connection.Data)
		mask.Type = dotdeploy.ConnectionTypeAgent
	default:
		return dotdeploy.ErrUnknowConnectionType
	}

	return err
}

func (connection *Connection) Connection() *dotdeploy.Connection {
	return &dotdeploy.Connection{
		Name: connection.Name,
		Data: connection.Data,
	}
}
