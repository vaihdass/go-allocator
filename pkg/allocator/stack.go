package allocator

import (
	"fmt"
	"math"
	"unsafe"
)

const headerSize = 4 // uint32 header size

// StackAllocator uses unaligned pointers
type StackAllocator struct {
	data []byte
}

func NewStackAllocator(capacity int) (*StackAllocator, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("allocator.NewStackAllocator: %w", ErrInvalidCapacity)
	}

	return &StackAllocator{
		make([]byte, 0, capacity),
	}, nil

}

func (a *StackAllocator) Allocate(size int) (Pointer, error) {
	const op = "allocator.StackAllocator.Allocate"
	if size <= 0 {
		return Pointer{}, fmt.Errorf("%s: %w", op, ErrNegativeElemSize)
	}

	// Check is overflows header size
	if size > math.MaxUint32 {
		// TODO: Can increase header size
		return Pointer{}, fmt.Errorf("%s: %w", op, ErrOutOfMemory)
	}

	prevLen := len(a.data)
	newLen := prevLen + headerSize + size

	// Capacity & overflow check
	if newLen > cap(a.data) || newLen < prevLen {
		// TODO: Can increase capacity
		return Pointer{}, fmt.Errorf("%s: %w", op, ErrOutOfMemory)
	}

	a.data = a.data[:newLen]

	header := unsafe.Pointer(&a.data[prevLen])
	Store[uint32](newPointer(header), uint32(size))

	pointer := unsafe.Pointer(&a.data[prevLen+headerSize])

	return newPointer(pointer), nil
}

// Deallocate vulnerability: potentially incorrect pointer.
// Using allocator.Pointer instead of unsafe.Pointer can protect against trivial cases of misuse
// because it doesn't expose a pointer outside the package.
func (a *StackAllocator) Deallocate(p Pointer) error {
	// TODO: Can deallocate without pointer

	if p.pointer == nil {
		return fmt.Errorf("allocator.StackAllocator.Deallocate: %w", ErrInvalidPointer)
	}

	// Getting element size
	header := unsafe.Add(p.pointer, -headerSize)
	size := Load[uint32](newPointer(header))

	// Getting new length (without element & its header)
	prevLen := len(a.data)
	newLen := prevLen - int(size) - headerSize

	a.data = a.data[:newLen]

	return nil
}

func (a *StackAllocator) Free() {
	a.data = a.data[:0]
}
