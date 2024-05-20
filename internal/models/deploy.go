package models

type (
	Remote struct {
		IPv4 string
		User string
		Port int
	}

	Remotes struct {
		Remotes map[string]*Remote
	}

	ScriptMove struct {
		From string
		To   string
	}

	Script struct {
		Move *ScriptMove
	}

	Scripts struct {
		Scripts map[string]*Script
	}

	Deploy struct {
		Remotes
		Scripts
	}
)
