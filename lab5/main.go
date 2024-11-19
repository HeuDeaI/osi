package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	BlockSize = 128
	NumBlocks = 32
)

type Block [BlockSize]byte

type File struct {
	Name      string
	BlockNums []int
	Size      int
}

type Directory struct {
	Name     string
	Children map[string]interface{}
	Parent   *Directory
}

type FileSystem struct {
	Root       *Directory
	Blocks     [NumBlocks]Block
	UsedBlocks []bool
	CurrDir    *Directory
}

func NewFileSystem() *FileSystem {
	root := &Directory{
		Name:     "/",
		Children: make(map[string]interface{}),
		Parent:   nil,
	}
	return &FileSystem{
		Root:       root,
		CurrDir:    root,
		UsedBlocks: make([]bool, NumBlocks),
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
	fs.CurrDir.Children[name] = &File{Name: name, BlockNums: []int{}, Size: 0}
}

func (fs *FileSystem) CreateDirectory(name string) {
	fs.CurrDir.Children[name] = &Directory{Name: name, Children: make(map[string]interface{}), Parent: fs.CurrDir}
}

func (fs *FileSystem) WriteToFile(name string, data []byte) {
	file := fs.findFile(name)
	bytesWritten := 0
	for bytesWritten < len(data) {
		blockNum := fs.AllocateBlock()
		copy(fs.Blocks[blockNum][:], data[bytesWritten:])
		file.BlockNums = append(file.BlockNums, blockNum)
		file.Size += BlockSize
		bytesWritten += BlockSize
	}
}

func (fs *FileSystem) ReadFromFile(name string) []byte {
	file := fs.findFile(name)
	data := []byte{}
	for _, blockNum := range file.BlockNums {
		data = append(data, fs.Blocks[blockNum][:]...)
	}
	return data
}

func (fs *FileSystem) ChangeDirectory(path string) {
	if path == "/" {
		fs.CurrDir = fs.Root
		return
	}
	segments := strings.Split(path, "/")
	curr := fs.CurrDir
	for _, seg := range segments {
		if seg == "" || seg == "." {
			continue
		}
		if seg == ".." {
			curr = curr.Parent
			continue
		}
		child := curr.Children[seg]
		if dir, ok := child.(*Directory); ok {
			curr = dir
		}
	}
	fs.CurrDir = curr
}

func (fs *FileSystem) ListDirectory() {
	for name, child := range fs.CurrDir.Children {
		switch child.(type) {
		case *File:
			fmt.Println(name)
		case *Directory:
			fmt.Println(name + "/")
		}
	}
}

func (fs *FileSystem) DeleteFile(name string) {
	file := fs.findFile(name)
	for _, blockNum := range file.BlockNums {
		fs.UsedBlocks[blockNum] = false
	}
	delete(fs.CurrDir.Children, name)
}

func (fs *FileSystem) CopyFile(srcName, destName string) {
	srcFile := fs.findFile(srcName)
	newFile := &File{Name: destName, BlockNums: []int{}, Size: srcFile.Size}
	for _, blockNum := range srcFile.BlockNums {
		newBlockNum := fs.AllocateBlock()
		copy(fs.Blocks[newBlockNum][:], fs.Blocks[blockNum][:])
		newFile.BlockNums = append(newFile.BlockNums, newBlockNum)
	}
	fs.CurrDir.Children[destName] = newFile
}

func (fs *FileSystem) MoveFile(srcName, destName string) {
	fs.CopyFile(srcName, destName)
	fs.DeleteFile(srcName)
}

func (fs *FileSystem) findFile(name string) *File {
	child := fs.CurrDir.Children[name]
	if file, ok := child.(*File); ok {
		return file
	}
	return nil
}

func (fs *FileSystem) GetCurrentPath() string {
	if fs.CurrDir == fs.Root {
		return "/"
	}
	path := ""
	curr := fs.CurrDir
	for curr != nil {
		if curr.Parent == nil {
			break
		}
		path = "/" + curr.Name + path
		curr = curr.Parent
	}
	return path
}

func main() {
	fs := NewFileSystem()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s: ", fs.GetCurrentPath())
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")
		command := args[0]
		switch command {
		case "ls":
			fs.ListDirectory()
		case "mkdir":
			fs.CreateDirectory(args[1])
		case "touch":
			fs.CreateFile(args[1])
		case "cd":
			fs.ChangeDirectory(args[1])
		case "pwd":
			fmt.Println(fs.GetCurrentPath())
		case "cat":
			data := fs.ReadFromFile(args[1])
			fmt.Println(string(data))
		case "echo":
			content := strings.Join(args[2:], " ")
			fs.WriteToFile(args[1], []byte(content))
		case "rm":
			fs.DeleteFile(args[1])
		case "cp":
			fs.CopyFile(args[1], args[2])
		case "mv":
			fs.MoveFile(args[1], args[2])
		case "exit":
			return
		}
	}
}
