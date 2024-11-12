package allocator

import "unsafe"

func Store[T any](p Pointer, v T) {
	*(*T)(p.pointer) = v
}

func Load[T any](p Pointer) T {
	return *(*T)(p.pointer)
}

type Pointer struct {
	pointer unsafe.Pointer
}

func newPointer(p unsafe.Pointer) Pointer {
	return Pointer{pointer: p}
}

type GenericPointer[T any] struct {
	pointer Pointer
}

func newGenericPointer[T any](p Pointer) GenericPointer[T] {
	return GenericPointer[T]{pointer: p}
}

func (p GenericPointer[T]) Load() T {
	return Load[T](p.pointer)
}

func (p GenericPointer[T]) Store(v T) {
	Store[T](p.pointer, v)
}
