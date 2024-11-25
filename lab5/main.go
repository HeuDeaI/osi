package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	BlockSize = 128
	NumBlocks = 32
)

type FileSystem struct {
	Blocks      [NumBlocks][BlockSize]byte
	UsedBlocks  [NumBlocks]bool
	Files       map[string][]int
	Directories map[string][]string
	CurrDir     string
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		Files:       make(map[string][]int),
		Directories: map[string][]string{"/": {}},
		CurrDir:     "/",
	}
}

func (fs *FileSystem) AllocateBlock() int {
	for i := 0; i < NumBlocks; i++ {
		if !fs.UsedBlocks[i] {
			fs.UsedBlocks[i] = true
			return i
		}
	}
	return -1
}

func (fs *FileSystem) CreateFile(name string) {
	fullPath := fs.CurrDir + "/" + name
	fs.Files[fullPath] = []int{}
	fs.Directories[fs.CurrDir] = append(fs.Directories[fs.CurrDir], name)
}

func (fs *FileSystem) CreateDirectory(name string) {
	fullPath := fs.CurrDir + "/" + name
	fs.Directories[fullPath] = []string{}
	fs.Directories[fs.CurrDir] = append(fs.Directories[fs.CurrDir], name+"/")
}

func (fs *FileSystem) WriteToFile(name string, data []byte) {
	fullPath := fs.CurrDir + "/" + name
	if _, exists := fs.Files[fullPath]; !exists {
		fmt.Println("File does not exist.")
		return
	}
	bytesWritten := 0
	for bytesWritten < len(data) {
		blockNum := fs.AllocateBlock()
		if blockNum == -1 {
			fmt.Println("No free blocks available.")
			return
		}
		copy(fs.Blocks[blockNum][:], data[bytesWritten:])
		fs.Files[fullPath] = append(fs.Files[fullPath], blockNum)
		bytesWritten += BlockSize
	}
}

func (fs *FileSystem) ReadFromFile(name string) []byte {
	fullPath := fs.CurrDir + "/" + name
	if _, exists := fs.Files[fullPath]; !exists {
		fmt.Println("File does not exist.")
		return nil
	}
	data := []byte{}
	for _, blockNum := range fs.Files[fullPath] {
		data = append(data, fs.Blocks[blockNum][:]...)
	}
	return data
}

func (fs *FileSystem) ChangeDirectory(path string) {
	if path == "/" {
		fs.CurrDir = "/"
	} else if path == ".." {
		parentDir := fs.CurrDir[:strings.LastIndex(fs.CurrDir, "/")]
		if parentDir == "" {
			parentDir = "/"
		}
		fs.CurrDir = parentDir
	} else {
		newPath := fs.CurrDir + "/" + path
		if _, exists := fs.Directories[newPath]; exists {
			fs.CurrDir = newPath
		} else {
			fmt.Println("Directory does not exist.")
		}
	}
}

func (fs *FileSystem) ListDirectory() {
	for _, name := range fs.Directories[fs.CurrDir] {
		fmt.Println(name)
	}
}

func (fs *FileSystem) SaveToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error saving file system:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, used := range fs.UsedBlocks {
		if used {
			writer.WriteString("1")
		} else {
			writer.WriteString("0")
		}
	}
	writer.WriteString("\n")

	for fileName, blocks := range fs.Files {
		writer.WriteString(fileName + ":" + intSliceToString(blocks) + "\n")
	}

	for dir, children := range fs.Directories {
		writer.WriteString(dir + "=" + strings.Join(children, ",") + "\n")
	}

	writer.Flush()
}

func (fs *FileSystem) LoadFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error loading file system:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	blockUsage := scanner.Text()
	for i, char := range blockUsage {
		fs.UsedBlocks[i] = char == '1'
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			fileName := parts[0]
			blocks := stringToIntSlice(parts[1])
			fs.Files[fileName] = blocks
		} else if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			dirName := parts[0]
			children := strings.Split(parts[1], ",")
			fs.Directories[dirName] = children
		}
	}
}

func intSliceToString(slice []int) string {
	strSlice := []string{}
	for _, num := range slice {
		strSlice = append(strSlice, strconv.Itoa(num))
	}
	return strings.Join(strSlice, ",")
}

func stringToIntSlice(data string) []int {
	strSlice := strings.Split(data, ",")
	intSlice := []int{}
	for _, str := range strSlice {
		num, _ := strconv.Atoi(str)
		intSlice = append(intSlice, num)
	}
	return intSlice
}

func main() {
	fs := NewFileSystem()
	fs.LoadFromFile("filesystem.txt")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s: ", fs.CurrDir)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")
		command := args[0]

		switch command {
		case "ls":
			fs.ListDirectory()
		case "mkdir":
			if len(args) < 2 {
				fmt.Println("Usage: mkdir <directory>")
				continue
			}
			fs.CreateDirectory(args[1])
		case "touch":
			if len(args) < 2 {
				fmt.Println("Usage: touch <file>")
				continue
			}
			fs.CreateFile(args[1])
		case "cd":
			if len(args) < 2 {
				fmt.Println("Usage: cd <path>")
				continue
			}
			fs.ChangeDirectory(args[1])
		case "pwd":
			fmt.Println(fs.CurrDir)
		case "cat":
			if len(args) < 2 {
				fmt.Println("Usage: cat <file>")
				continue
			}
			data := fs.ReadFromFile(args[1])
			if data != nil {
				fmt.Println(string(data))
			}
		case "echo":
			if len(args) < 3 {
				fmt.Println("Usage: echo <file> <content>")
				continue
			}
			content := strings.Join(args[2:], " ")
			fs.WriteToFile(args[1], []byte(content))
		case "exit":
			fs.SaveToFile("filesystem.txt")
			fmt.Println("File system saved. Exiting...")
			return
		default:
			fmt.Println("Unknown command:", command)
		}
	}
}
