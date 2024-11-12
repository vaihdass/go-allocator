package allocator

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	ErrInvalidCapacity  = errors.New("capacity must be greater than zero")
	ErrNegativeElemSize = errors.New("size of element to allocate must be greater than zero")
	ErrOutOfMemory      = errors.New("allocator is out of memory")
)

// LinearAllocator uses unaligned pointers
type LinearAllocator struct {
	data []byte
}

func NewLinearAllocator(capacity int) (*LinearAllocator, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("allocator.NewLinearAllocator: %w", ErrInvalidCapacity)
	}

	return &LinearAllocator{
		data: make([]byte, capacity),
	}, nil
}

func (a *LinearAllocator) Allocate(size int) (Pointer, error) {
	const op = "allocator.LinearAllocator.Allocate"

	if size <= 0 {
		return Pointer{}, fmt.Errorf("%s: %w", op, ErrNegativeElemSize)
	}

	prevLen := len(a.data)
	newLen := prevLen + size

	// Capacity & overflow check
	if newLen > cap(a.data) || newLen < prevLen {
		// TODO: Can increase capacity
		return Pointer{}, fmt.Errorf("%s: %w", op, ErrOutOfMemory)
	}

	a.data = a.data[:newLen]
	pointer := unsafe.Pointer(&a.data[prevLen])

	return newPointer(pointer), nil
}

func (a *LinearAllocator) Free() {
	a.data = a.data[:0]
}
