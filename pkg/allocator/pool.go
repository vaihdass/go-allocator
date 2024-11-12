package allocator

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	ErrNegativeBlockMetric = errors.New("negative block metric: size or count")
	ErrInvalidPointer      = errors.New("invalid pointer")
)

const (
	Int32Size = 4
	Int64Size = 8
)

// PoolAllocator uses unaligned pointers
type PoolAllocator struct {
	pool       []byte
	freeBlocks map[unsafe.Pointer]struct{}
	blockSize  int
}

func NewPoolAllocator(blockSize, blockCount int) (*PoolAllocator, error) {
	const op = "allocator.NewPoolAllocator"
	if blockSize <= 0 || blockCount <= 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrNegativeBlockMetric)
	}

	capacity := blockSize * blockCount
	if capacity < blockCount { // Overflow check
		return nil, fmt.Errorf("%s: %w", op, ErrOutOfMemory)
	}

	allocator := PoolAllocator{
		pool:       make([]byte, capacity),
		freeBlocks: make(map[unsafe.Pointer]struct{}, blockCount),
		blockSize:  blockSize,
	}
	allocator.resetMemoryState()

	return &allocator, nil
}

func (a *PoolAllocator) Allocate() (Pointer, error) {
	if len(a.freeBlocks) == 0 {
		// TODO: Can increase capacity
		return Pointer{}, fmt.Errorf("allocator.PoolAllocator.Allocate: %w", ErrOutOfMemory)
	}

	var pointer unsafe.Pointer
	for freePointer := range a.freeBlocks {
		pointer = freePointer
		break
	}
	delete(a.freeBlocks, pointer)

	return newPointer(pointer), nil
}

// Deallocate vulnerability: potentially incorrect pointer.
// Using allocator.Pointer instead of unsafe.Pointer can protect against trivial cases of misuse
// because it doesn't expose a pointer outside the package.
func (a *PoolAllocator) Deallocate(p Pointer) error {
	if p.pointer == nil {
		return fmt.Errorf("allocator.PoolAllocator.Deallocate: %w", ErrInvalidPointer)
	}

	a.freeBlocks[p.pointer] = struct{}{}

	return nil
}

func (a *PoolAllocator) Free() {
	a.resetMemoryState()
}

func (a *PoolAllocator) resetMemoryState() {
	for offset := 0; offset < len(a.pool); offset += a.blockSize {
		p := unsafe.Pointer(&a.pool[offset])
		a.freeBlocks[p] = struct{}{}
	}
}
