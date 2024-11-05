package main

import (
	"fmt"
)

const (
	MEMORY_SIZE    = 512
	PAGE_SIZE      = 32
	COUNT_OF_PAGES = MEMORY_SIZE / PAGE_SIZE
)

type Pointer struct {
	DataOffset int
	Size       int
}

type Memory struct {
	Data      [MEMORY_SIZE]rune
	FreeSpace *Pointer
	FreeList  []*Pointer
}

func SetMemory() *Memory {
	return &Memory{
		FreeSpace: &Pointer{DataOffset: 0},
		FreeList:  []*Pointer{},
	}
}

func (memory *Memory) allocate(size int) *Pointer {
	for i, block := range memory.FreeList {
		if block.Size >= size {
			memory.FreeList = append(memory.FreeList[:i], memory.FreeList[i+1:]...)
			return block
		}
	}

	if (memory.FreeSpace.DataOffset / PAGE_SIZE) != ((memory.FreeSpace.DataOffset + size) / PAGE_SIZE) {
		memory.FreeSpace.DataOffset = (memory.FreeSpace.DataOffset/PAGE_SIZE + 1) * PAGE_SIZE
	}

	ptr := &Pointer{DataOffset: memory.FreeSpace.DataOffset, Size: size}
	memory.FreeSpace.DataOffset += size
	return ptr
}

func (memory *Memory) free(ptr *Pointer) {
	for i := ptr.DataOffset; i < ptr.DataOffset+ptr.Size; i++ {
		memory.Data[i] = 0
	}

	memory.FreeList = append(memory.FreeList, ptr)
}

func (memory *Memory) write(ptr *Pointer, data []rune) {
	copy(memory.Data[ptr.DataOffset:ptr.DataOffset+ptr.Size], data)
}

func (memory *Memory) read(ptr *Pointer) []rune {
	startLineIndex := ptr.DataOffset
	endLineIndex := startLineIndex + ptr.Size
	return memory.Data[startLineIndex:endLineIndex]
}

func (memory *Memory) viewMemory() {
	for page := 0; page < COUNT_OF_PAGES; page++ {
		startIndex := page * PAGE_SIZE
		endIndex := startIndex + PAGE_SIZE
		if endIndex > MEMORY_SIZE {
			endIndex = MEMORY_SIZE
		}

		pageData := memory.Data[startIndex:endIndex]
		fmt.Printf("Page №%d: &{%v}\n", page, pageData)
	}
}

func (memory *Memory) InsertData(data string) *Pointer {
	ptr := memory.allocate(len(data))
	memory.write(ptr, []rune(data))
	return ptr
}

func main() {
	memory := SetMemory()

	dataSet := []string{
		"hello", "world", "memory", "allocation", "test", "data", "set",
		"programming", "in", "golang", "is", "fun", "and", "educational",
		"let's", "try", "some", "longer", "strings", "to", "test", "memory",
		"management", "in", "a", "simplified", "memory", "simulator", "using",
		"pages", "and", "allocation", "methods", "this", "should", "help",
		"identify", "any", "potential", "issues", "with", "our", "current",
		"memory", "layout", "and", "free", "space", "tracking", "mechanism",
		"keep", "adding", "more", "strings", "to", "fill", "the", "pages",
		"and observe how allocation behaves when memory is almost full or completely filled",
	}

	var pointers []*Pointer
	for _, data := range dataSet {
		ptr := memory.InsertData(data)
		pointers = append(pointers, ptr)
	}
	memory.viewMemory()

	memory.free(pointers[1])
	memory.free(pointers[3])

	memory.InsertData("reuse")
	memory.InsertData("big reuse")
	fmt.Println()
	memory.viewMemory()
}
