package dotdeploy

import (
	"errors"

	"github.com/melbahja/goph"
)

type (
	Remote []*Connection

	Connection struct {
		Name string
		Data any

		Client *goph.Client

		Error error
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

func Connect(connection *Connection) *Connection {
	switch connection.Data.(type) {
	case *Key:
		connection.Error = connection.Key()
	case *Password:
		connection.Error = connection.Password()
	case *Agent:
		connection.Error = connection.Agent()
	default:
		connection.Error = ErrUnknowConnectionType
	}

	return connection
}

func (remote Remote) Find(name string) (connection *Connection, ok bool) {
	for i, value := range remote {
		if value.Name == name {
			return remote[i], true
		}
	}

	return
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
