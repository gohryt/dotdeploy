package main

import (
	"errors"

	"github.com/melbahja/goph"
	"gopkg.in/yaml.v3"
)

type (
	Connection struct {
		Name string

		Data any

		Client *goph.Client
	}

	ConnectionType struct {
		Type string `yaml:"type"`
		Name string `yaml:"name"`
	}

	Key struct {
		Host     string `yaml:"host" validate:"required"`
		Username string `yaml:"username" validate:"required"`
		File     string `yaml:"file" validate:"required"`
		Password string `yaml:"password" validate:"required"`
	}

	Password struct {
		Host     string `yaml:"host" validate:"required"`
		Username string `yaml:"username" validate:"required"`
		Password string `yaml:"password" validate:"required"`
	}

	Agent struct {
		Host     string `yaml:"host" validate:"required"`
		Username string `yaml:"username" validate:"required"`
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

	switch t.Type {
	case "key":
		connection.Data = new(Key)
	case "password":
		connection.Data = new(Password)
	case "agent":
		connection.Data = new(Agent)
	default:
		return errors.New("unknown connection type")
	}

	return value.Decode(connection.Data)
}

func (deploy *Deploy) Connect(connection *Connection) error {
	switch connection.Data.(type) {
	case *Key:
		return connection.Key()
	case *Password:
		return connection.Password()
	case *Agent:
		return connection.Agent()
	}

	return errors.New("unknown connection type")
}

func (connection *Connection) Key() error {
	key := connection.Data.(*Key)

	authentication, err := goph.Key(key.File, key.Password)
	if err != nil {
		return err
	}

	connection.Client, err = goph.New(key.Username, key.Host, authentication)
	return err
}

func (connection *Connection) Password() error {
	password := connection.Data.(*Password)
	err := error(nil)

	connection.Client, err = goph.New(password.Username, password.Host, goph.Password(password.Password))
	return err
}

func (connection *Connection) Agent() error {
	agent := connection.Data.(*Agent)

	authentication, err := goph.UseAgent()
	if err != nil {
		return err
	}

	connection.Client, err = goph.New(agent.Username, agent.Host, authentication)
	return err
}
