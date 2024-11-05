package main

import "fmt"

const (
	MEMORY_SIZE    = 512
	PAGE_SIZE      = 32
	COUNT_OF_PAGES = MEMORY_SIZE / PAGE_SIZE
)

type Page struct {
	data      [PAGE_SIZE]rune
	freeSpace int
}

type Pointer struct {
	pageNumber   int
	dataLocation int
}

type Memory struct {
	pages [COUNT_OF_PAGES]*Page
}

func setMemory() *Memory {
	memory := &Memory{}
	for i := 0; i < COUNT_OF_PAGES; i++ {
		memory.pages[i] = &Page{freeSpace: PAGE_SIZE}
	}
	return memory
}

func (memory *Memory) allocate(size int) *Pointer {
	for i := 0; i < COUNT_OF_PAGES; i++ {
		if freeSpace := &memory.pages[i].freeSpace; *freeSpace > size {
			ptr := &Pointer{pageNumber: i, dataLocation: PAGE_SIZE - *freeSpace}
			*freeSpace -= size
			return ptr
		}
	}
	return nil
}

func (memory *Memory) free(size int) {
}

func (memory *Memory) write(ptr *Pointer, data []rune) {
	pageSpace := memory.pages[ptr.pageNumber].data[ptr.dataLocation:]
	copy(pageSpace, data)
}

func (memory *Memory) read(ptr *Pointer) []rune {
	page := memory.pages[ptr.pageNumber]
	startLineIndex := ptr.dataLocation
	endLineIndex := 0
	for i := ptr.dataLocation; i < PAGE_SIZE; i++ {
		if page.data[i] == 0 {
			endLineIndex = i
			break
		}
	}
	return page.data[startLineIndex:endLineIndex]
}

func (memory *Memory) viewMemory() {
	for i := 0; i < COUNT_OF_PAGES; i++ {
		fmt.Printf("Page №%v: %v\n", i, memory.pages[i])
	}
}

func (memory *Memory) InsertData(data string) {
	ptr := memory.allocate(len(data) + 1)
	memory.write(ptr, []rune(data))
}

func main() {
	memory := setMemory()

	dataSet := []string{
		"hello", "world", "memory", "allocation", "test", "data", "set",
		"programming", "in", "golang", "is", "fun", "and", "educational",
		"let's", "try", "some", "longer", "strings", "to", "test", "memory",
		"management", "in", "a", "simplified", "memory", "simulator", "using",
		"pages", "and", "allocation", "methods", "this", "should", "help",
		"identify", "any", "potential", "issues", "with", "our", "current",
		"memory", "layout", "and", "free", "space", "tracking", "mechanism",
		"keep", "adding", "more", "strings", "to", "fill", "the", "pages",
		"and", "observe", "how", "allocation", "behaves", "when", "memory",
		"is", "almost", "full", "or", "completely", "filled",
	}
	for _, data := range dataSet {
		memory.InsertData(data)
	}

	memory.viewMemory()
}
