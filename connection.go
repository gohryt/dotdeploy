package dotdeploy

import (
	"errors"

	"github.com/melbahja/goph"
)

type (
	Connection struct {
		Name string
		Data any

		Client *goph.Client
	}

	ConnectionType struct {
		Type string
		Name string
	}

	Key struct {
		Host     string `validate:"required"`
		Username string `validate:"required"`
		File     string `validate:"required"`
		Password string `validate:"required"`
	}

	Password struct {
		Host     string `validate:"required"`
		Username string `validate:"required"`
		Password string `validate:"required"`
	}

	Agent struct {
		Host     string `validate:"required"`
		Username string `validate:"required"`
	}
)

var (
	ErrUnknowConnectionType = errors.New("unknown connection type")
)

func (deploy *Deploy) Connect(connection *Connection) error {
	switch connection.Data.(type) {
	case *Key:
		return connection.Key()
	case *Password:
		return connection.Password()
	case *Agent:
		return connection.Agent()
	}

	return ErrUnknowConnectionType
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
