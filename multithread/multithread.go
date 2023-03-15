package multithread

func Go[T any](receiver chan *T, list []*T, process func(object *T) *T) {
	for i := range list {
		go func(object *T) {
			receiver <- process(object)
		}(list[i])
	}
}
