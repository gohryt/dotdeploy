package main

import (
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/melbahja/goph"
)

type (
	Connection struct {
		Data Checkable

		Name   string
		Client *goph.Client
	}

	ConnectionType struct {
		Type string `json:"type"`
	}

	Key struct {
		Name string

		Host     string
		File     string
		Username string
		Password string
	}

	Password struct {
		Name string

		Host     string
		Username string
		Password string
	}

	Agent struct {
		Name string

		Host     string
		Username string
	}
)

var (
	empty = Empty("empty")
)

func (connection *Connection) UnmarshalJSON(source []byte) error {
	t := new(ConnectionType)

	err := json.Unmarshal(source, t)
	if err != nil {
		return err
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

		return nil
	}

	return json.Unmarshal(source, connection.Data)
}

func (key *Key) Check() error {
	err := error(nil)

	if key.Name == "" {
		err = errors.Join(err, errors.New("'name' can't be empty"))
	}

	if key.Host == "" {
		err = errors.Join(err, errors.New("'host' can't be empty"))
	}

	if key.File == "" {
		err = errors.Join(err, errors.New("'file' can't be empty"))
	}

	if key.Username == "" {
		err = errors.Join(err, errors.New("'username' can't be empty"))
	}

	if key.Password == "" {
		err = errors.Join(err, errors.New("'password' can't be empty"))
	}

	return err
}

func (key *Key) String() string {
	return fmt.Sprintf("connect by key %s with host %s file %s username %s password %s", key.Name, key.Host, key.File, key.Username, key.Password)
}

func (password *Password) Check() error {
	err := error(nil)

	if password.Name == "" {
		err = errors.Join(err, errors.New("'name' can't be empty"))
	}

	if password.Host == "" {
		err = errors.Join(err, errors.New("'host' can't be empty"))
	}

	if password.Username == "" {
		err = errors.Join(err, errors.New("'username' can't be empty"))
	}

	if password.Password == "" {
		err = errors.Join(err, errors.New("'password' can't be empty"))
	}

	return err
}

func (password *Password) String() string {
	return fmt.Sprintf("connect by password %s with host %s username %s password %s", password.Name, password.Host, password.Username, password.Password)
}

func (agent *Agent) Check() error {
	err := error(nil)

	if agent.Name == "" {
		err = errors.Join(err, errors.New("'name' can't be empty"))
	}

	if agent.Host == "" {
		err = errors.Join(err, errors.New("'host' can't be empty"))
	}

	if agent.Username == "" {
		err = errors.Join(err, errors.New("'username' can't be empty"))
	}

	return err
}

func (agent *Agent) String() string {
	return fmt.Sprintf("connect by agent %s with host %s username %s", agent.Name, agent.Host, agent.Username)
}

func (deploy *Deploy) Connect(connection *Connection) error {
	data := connection.Data

	switch data.(type) {
	case *Key:
		key := data.(*Key)

		client, err := deploy.Key(key)

		connection.Name = key.Name
		connection.Client = client

		return err
	case *Password:
		password := data.(*Password)

		client, err := deploy.Password(password)

		connection.Name = password.Name
		connection.Client = client

		return err
	case *Agent:
		agent := data.(*Agent)

		client, err := deploy.Agent(agent)

		connection.Name = agent.Name
		connection.Client = client

		return err
	default:
		return errors.New("unknown connection type")
	}
}

func (deploy *Deploy) Key(key *Key) (client *goph.Client, err error) {
	authentication, err := goph.Key(key.File, key.Password)
	if err != nil {
		return
	}

	client, err = goph.New(key.Username, key.Host, authentication)
	return
}

func (deploy *Deploy) Password(password *Password) (client *goph.Client, err error) {
	client, err = goph.New(password.Username, password.Host, goph.Password(password.Password))
	return
}

func (deploy *Deploy) Agent(agent *Agent) (client *goph.Client, err error) {
	authentication, err := goph.UseAgent()
	if err != nil {
		return
	}

	client, err = goph.New(agent.Username, agent.Host, authentication)
	return
}
