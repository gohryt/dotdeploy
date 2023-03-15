package main

type (
	Do []*Action
)

func (do Do) Find(name string) (action *Action, ok bool) {
	for i, value := range do {
		if value.Name == name {
			return do[i], true
		}
	}

	return
}
