package unsafe

import (
	"unsafe"
)

type (
	Any struct {
		Type unsafe.Pointer
		Data unsafe.Pointer
	}
)

func Type(source any) unsafe.Pointer {
	return As[Any](&source).Type
}

func Data(source any) unsafe.Pointer {
	return As[Any](&source).Data
}

func As[T any, S any](source *S) *T {
	return (*T)(unsafe.Pointer(source))
}

func New[T any]() Any {
	t := new(T)

	return Any{
		Type: Type(t),
		Data: unsafe.Pointer(t),
	}
}
