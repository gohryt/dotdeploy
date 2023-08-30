package contract

import (
	"log"
	"unsafe"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"gopkg.in/yaml.v3"
)

type (
	TypeRemote int
	TypeScript int

	EnumRemote struct {
		Type TypeRemote
		Data unsafe.Pointer
	}

	EnumScript struct {
		Type TypeScript
		Data unsafe.Pointer
	}

	Item struct {
		Type   string   `yaml:"type" json:"type"`
		Name   string   `yaml:"name" json:"name"`
		Follow []string `yaml:"follow" json:"follow"`
	}

	ItemRemote struct {
		Item
		EnumRemote
	}

	ItemScript struct {
		Item
		EnumScript
	}

	Root struct {
		Remote []ItemRemote `yaml:"remote" json:"remote"`
		Script []ItemScript `yaml:"script" json:"script"`
	}
)

func (remote ItemRemote) MarshalJSON() (buffer []byte, err error) {
	buffer, err = remote.EnumRemote.MarshalJSON()
	if err != nil {
		return
	}

	node := ast.NewRaw(unsafe.String(unsafe.SliceData(buffer), len(buffer)))

	_, err = node.Set("type", ast.NewString(remote.Item.Type))
	if err != nil {
		return
	}

	_, err = node.Set("name", ast.NewString(remote.Name))
	if err != nil {
		return
	}

	follow := []ast.Node(nil)

	for _, value := range remote.Follow {
		follow = append(follow, ast.NewString(value))
	}

	_, err = node.Set("follow", ast.NewArray(follow))
	if err != nil {
		return
	}

	return node.MarshalJSON()
}

func (script ItemScript) MarshalJSON() (buffer []byte, err error) {
	buffer, err = script.EnumScript.MarshalJSON()
	if err != nil {
		log.Println(111)
		return
	}

	node := ast.NewRaw(unsafe.String(unsafe.SliceData(buffer), len(buffer)))

	_, err = node.Set("type", ast.NewString(script.Item.Type))
	if err != nil {
		return
	}

	_, err = node.Set("name", ast.NewString(script.Name))
	if err != nil {
		return
	}

	follow := []ast.Node(nil)

	for _, value := range script.Follow {
		follow = append(follow, ast.NewString(value))
	}

	_, err = node.Set("follow", ast.NewArray(follow))
	if err != nil {
		return
	}

	buffer, err = node.MarshalJSON()
	return
}

func (remote EnumRemote) MarshalJSON() (buffer []byte, err error) {
	switch remote.Type {
	case RemoteAgent:
		return sonic.ConfigStd.Marshal((*Agent)(remote.Data))
	}
	return
}

func (script EnumScript) MarshalJSON() (buffer []byte, err error) {
	switch script.Type {
	case ScriptCopy:
		return sonic.ConfigStd.Marshal((*Copy)(script.Data))
	case ScriptMove:
		return sonic.ConfigStd.Marshal((*Move)(script.Data))
	case ScriptExecute:
		return sonic.ConfigStd.Marshal((*Execute)(script.Data))
	}
	return
}

func (remote *ItemRemote) UnmarshalYAML(node *yaml.Node) error {
	err := node.Decode(&remote.Item)
	if err != nil {
		return err
	}

	switch remote.Item.Type {
	case "agent":
		data := new(Agent)

		err = node.Decode(data)
		if err != nil {
			return err
		}

		remote.EnumRemote = EnumRemote{
			Type: RemoteAgent,
			Data: unsafe.Pointer(data),
		}
	default:
		remote.EnumRemote.Type = RemoteUnknown
	}

	return nil
}

func (script *ItemScript) UnmarshalYAML(node *yaml.Node) error {
	err := node.Decode(&script.Item)
	if err != nil {
		return err
	}

	switch script.Item.Type {
	case "copy":
		data := new(Copy)

		err = node.Decode(data)
		if err != nil {
			return err
		}

		script.EnumScript = EnumScript{
			Type: ScriptCopy,
			Data: unsafe.Pointer(data),
		}
	case "move":
		data := new(Move)

		err = node.Decode(data)
		if err != nil {
			return err
		}

		script.EnumScript = EnumScript{
			Type: ScriptMove,
			Data: unsafe.Pointer(data),
		}
	case "execute":
		data := new(Execute)

		err = node.Decode(data)
		if err != nil {
			return err
		}

		script.EnumScript = EnumScript{
			Type: ScriptExecute,
			Data: unsafe.Pointer(data),
		}
	default:
		script.EnumScript.Type = ScriptUnknown
	}

	return nil
}

func (remote *ItemRemote) UnmarshalJSON(buffer []byte) error {
	err := sonic.ConfigStd.Unmarshal(buffer, &remote.Item)
	if err != nil {
		return err
	}

	switch remote.Item.Type {
	case "agent":
		data := new(Agent)

		err = sonic.ConfigStd.Unmarshal(buffer, data)
		if err != nil {
			return err
		}

		remote.EnumRemote = EnumRemote{
			Type: RemoteAgent,
			Data: unsafe.Pointer(data),
		}
	default:
		remote.EnumRemote.Type = RemoteUnknown
	}

	return nil
}

func (script *ItemScript) UnmarshalJSON(buffer []byte) error {
	err := sonic.ConfigStd.Unmarshal(buffer, &script.Item)
	if err != nil {
		return err
	}

	switch script.Item.Type {
	case "copy":
		data := new(Copy)

		err = sonic.ConfigStd.Unmarshal(buffer, data)
		if err != nil {
			return err
		}

		script.EnumScript = EnumScript{
			Type: ScriptCopy,
			Data: unsafe.Pointer(data),
		}
	case "move":
		data := new(Move)

		err = sonic.ConfigStd.Unmarshal(buffer, data)
		if err != nil {
			return err
		}

		script.EnumScript = EnumScript{
			Type: ScriptMove,
			Data: unsafe.Pointer(data),
		}
	case "execute":
		data := new(Execute)

		err = sonic.ConfigStd.Unmarshal(buffer, data)
		if err != nil {
			return err
		}

		script.EnumScript = EnumScript{
			Type: ScriptExecute,
			Data: unsafe.Pointer(data),
		}
	default:
		script.EnumScript.Type = ScriptUnknown
	}

	return nil
}

type (
	Path struct {
		Path   string `json:"path"`
		Remote string `json:"remote"`
	}
)

type (
	Agent struct {
		Host     string `json:"host"`
		Username string `json:"username"`
	}
)

const (
	RemoteUnknown TypeRemote = iota
	RemoteAgent
)

type (
	Copy struct {
		From Path `yaml:"from" json:"from"`
		To   Path `yaml:"to" json:"to"`
	}

	Move struct {
		From Path `yaml:"from" json:"from"`
		To   Path `yaml:"to" json:"to"`
	}

	Execute struct {
		Path        Path     `yaml:"path" json:"path"`
		Environment []string `yaml:"environment" json:"environment"`
		Query       []string `yaml:"query" json:"query"`
	}
)

const (
	ScriptUnknown TypeScript = iota
	ScriptCopy
	ScriptMove
	ScriptExecute
)
