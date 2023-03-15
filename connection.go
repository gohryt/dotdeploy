package dotdeploy

import (
	"errors"

	"github.com/melbahja/goph"

	"github.com/gohryt/dotdeploy/unsafe"
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

	ErrKeyHostEmpty          = errors.New(`key.Host == ""`)
	ErrKeyUsernameEmpty      = errors.New(`key.Username == ""`)
	ErrKeyFileEmpty          = errors.New(`key.File == ""`)
	ErrKeyPasswordEmpty      = errors.New(`key.Password == ""`)
	ErrPasswordHostEmpty     = errors.New(`password.Host == ""`)
	ErrPasswordUsernameEmpty = errors.New(`password.Username == ""`)
	ErrPasswordPasswordEmpty = errors.New(`password.Password == ""`)
	ErrAgentHostEmpty        = errors.New(`agent.Host == ""`)
	ErrAgentUsernameEmpty    = errors.New(`agent.Username == ""`)
)

var (
	ConnectionTypeKey      = unsafe.Type(new(Key))
	ConnectionTypePassword = unsafe.Type(new(Password))
	ConnectionTypeAgent    = unsafe.Type(new(Agent))
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

func (key *Key) Validate() error {
	join := []error(nil)

	if key.Host == "" {
		join = append(join, ErrKeyHostEmpty)
	}

	if key.Username == "" {
		join = append(join, ErrKeyUsernameEmpty)
	}

	if key.File == "" {
		join = append(join, ErrKeyFileEmpty)
	}

	if key.Password == "" {
		join = append(join, ErrKeyPasswordEmpty)
	}

	return errors.Join(join...)
}

func (password *Password) Validate() error {
	join := []error(nil)

	if password.Host == "" {
		join = append(join, ErrPasswordHostEmpty)
	}

	if password.Username == "" {
		join = append(join, ErrPasswordUsernameEmpty)
	}

	if password.Password == "" {
		join = append(join, ErrPasswordPasswordEmpty)
	}

	return errors.Join(join...)
}

func (agent *Agent) Validate() error {
	join := []error(nil)

	if agent.Host == "" {
		join = append(join, ErrAgentHostEmpty)
	}

	if agent.Username == "" {
		join = append(join, ErrAgentUsernameEmpty)
	}

	return errors.Join(join...)
}
