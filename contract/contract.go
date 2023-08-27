package contract

type (
	Path struct {
		Path    string `json:"path"`
		Machine string `json:"machine"`
	}
)

type (
	Copy struct {
		From, To Path
	}

	Move struct {
		From, To Path
	}

	Execute struct {
		Path               Path
		Environment, Query []string
	}
)

type (
	Command[T any] struct {
		Type string `json:"type"`
		Data T
	}

	Return struct {
		Type string `json:"type"`
	}
)
