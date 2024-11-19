package main

import (
	"fmt"
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
}

type FileSystem struct {
	Root       *Directory
	Blocks     [NumBlocks]Block
	UsedBlocks []bool
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		Root: &Directory{
			Name:     "/",
			Children: make(map[string]interface{}),
		},
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

func (fs *FileSystem) CreateFile(path string, name string) {
	dir := fs.findDirectory(path)
	dir.Children[name] = &File{Name: name, BlockNums: []int{}, Size: 0}
}

func (fs *FileSystem) CreateDirectory(path string, name string) {
	dir := fs.findDirectory(path)
	dir.Children[name] = &Directory{Name: name, Children: make(map[string]interface{})}
}

func (fs *FileSystem) WriteToFile(path string, data []byte) {
	file := fs.findFile(path)
	bytesWritten := 0
	for bytesWritten < len(data) {
		blockNum := fs.AllocateBlock()
		copy(fs.Blocks[blockNum][:], data[bytesWritten:])
		file.BlockNums = append(file.BlockNums, blockNum)
		file.Size += BlockSize
		bytesWritten += BlockSize
	}
}

func (fs *FileSystem) ReadFromFile(path string) []byte {
	file := fs.findFile(path)
	data := []byte{}
	for _, blockNum := range file.BlockNums {
		data = append(data, fs.Blocks[blockNum][:]...)
	}
	return data
}

func (fs *FileSystem) findDirectory(path string) *Directory {
	segments := strings.Split(path, "/")
	curr := fs.Root
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		curr = curr.Children[seg].(*Directory)
	}
	return curr
}

func (fs *FileSystem) findFile(path string) *File {
	dirPath, fileName := splitPath(path)
	dir := fs.findDirectory(dirPath)
	return dir.Children[fileName].(*File)
}

func splitPath(path string) (dirPath, fileName string) {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return "/", path
	}
	return path[:lastSlash], path[lastSlash+1:]
}

func (fs *FileSystem) CopyFile(srcPath, destPath string) {
	srcFile := fs.findFile(srcPath)
	destDirPath, destFileName := splitPath(destPath)
	destDir := fs.findDirectory(destDirPath)

	newFile := &File{Name: destFileName, BlockNums: []int{}, Size: srcFile.Size}
	for _, blockNum := range srcFile.BlockNums {
		newBlockNum := fs.AllocateBlock()
		copy(fs.Blocks[newBlockNum][:], fs.Blocks[blockNum][:])
		newFile.BlockNums = append(newFile.BlockNums, newBlockNum)
	}
	destDir.Children[destFileName] = newFile
}

func (fs *FileSystem) MoveFile(srcPath, destPath string) {
	fs.CopyFile(srcPath, destPath)
	fs.DeleteFile(srcPath)
}

func (fs *FileSystem) DeleteFile(path string) {
	dirPath, fileName := splitPath(path)
	dir := fs.findDirectory(dirPath)
	file := dir.Children[fileName].(*File)

	for _, blockNum := range file.BlockNums {
		fs.UsedBlocks[blockNum] = false
	}
	delete(dir.Children, fileName)
}

func (fs *FileSystem) Dump() {
	fs.dumpDirectory(fs.Root, "")
}

func (fs *FileSystem) dumpDirectory(dir *Directory, indent string) {
	fmt.Println(indent + dir.Name + "/")
	for name, child := range dir.Children {
		switch child := child.(type) {
		case *File:
			fmt.Printf("%s  %s\n", indent, name)
		case *Directory:
			fs.dumpDirectory(child, indent+"  ")
		}
	}
}

func main() {
	fs := NewFileSystem()

	fs.CreateFile("/", "file1.txt")
	fs.WriteToFile("/file1.txt", []byte("Hello, File System!"))

	fs.CreateFile("/", "file2.txt")
	fs.WriteToFile("/file2.txt", []byte("Another file with some data."))

	fs.CreateDirectory("/", "docs")
	fs.CreateFile("/docs", "report.pdf")
	fs.WriteToFile("/docs/report.pdf", []byte("This is a report document."))

	fs.CreateDirectory("/", "images")
	fs.CreateFile("/images", "photo1.jpg")
	fs.WriteToFile("/images/photo1.jpg", []byte("Image data 1"))

	fs.CreateDirectory("/images", "thumbnails")
	fs.CreateFile("/images/thumbnails", "thumb1.jpg")
	fs.WriteToFile("/images/thumbnails/thumb1.jpg", []byte("Thumbnail image data"))

	fs.CopyFile("/file1.txt", "/docs/file1_copy.txt")
	fs.MoveFile("/images/photo1.jpg", "/docs/photo1.jpg")

	fmt.Println("File System Structure:")
	fs.Dump()

	fmt.Println("\nRead data from /docs/photo1.jpg:")
	data := fs.ReadFromFile("/docs/photo1.jpg")
	fmt.Println(string(data))
}
