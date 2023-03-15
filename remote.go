package main

type (
	Remote []*Connection
)

func (remote Remote) Find(name string) (connection *Connection, ok bool) {
	for i, value := range remote {
		if value.Name == name {
			return remote[i], true
		}
	}

	return
}
