package utils

import (
	"bufio"
	"io"
	"os"
)

var UtilFile = NewFile()

func NewFile() *File {
	return &File{}
}

type File struct {
}

// create file
func (f *File) CreateFile(filename string) error {
	newFile, err := os.Create(filename)
	defer newFile.Close()
	return err
}

// file or path is exists
func (f *File) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//get file lines
//params : filename
//return : fileLine, error
func (f *File) GetFileLines(filename string) (fileLine int64, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0766)
	if err != nil {
		return fileLine, err
	}
	defer file.Close()

	fileLine = 1
	r := bufio.NewReader(file)
	for {
		_, err := r.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		fileLine += 1
	}
	return fileLine, nil
}
