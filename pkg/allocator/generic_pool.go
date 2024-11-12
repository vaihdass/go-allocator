package allocator

import (
	"fmt"
	"unsafe"
)

// GenericPoolAllocator uses unaligned pointers
type GenericPoolAllocator[T any] struct {
	poolAllocator *PoolAllocator
}

func NewGenericPoolAllocator[T any](blockCount int) (*GenericPoolAllocator[T], error) {
	var example T
	a, err := NewPoolAllocator(int(unsafe.Sizeof(example)), blockCount)
	if err != nil {
		return nil, fmt.Errorf("allocator.NewGenericPoolAllocator: %w", err)
	}

	return &GenericPoolAllocator[T]{poolAllocator: a}, nil
}

func (a *GenericPoolAllocator[T]) Allocate() (GenericPointer[T], error) {
	p, err := a.poolAllocator.Allocate()
	if err != nil {
		return GenericPointer[T]{}, fmt.Errorf("allocator.GenericPoolAllocator.Allocate: %w", err)
	}

	return newGenericPointer[T](p), nil
}

func (a *GenericPoolAllocator[T]) Deallocate(p GenericPointer[T]) error {
	err := a.poolAllocator.Deallocate(p.pointer)
	if err != nil {
		return fmt.Errorf("allocator.GenericPoolAllocator.Deallocate: %w", err)
	}

	return nil
}

func (a *GenericPoolAllocator[T]) Free() {
	a.poolAllocator.Free()
}
