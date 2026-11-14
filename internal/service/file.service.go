package service

import (
	"bufio"
	"io"
	"os"
	"path/filepath"

	"github.com/c12s/metrics/internal/errors"
)

type LocalFileService struct{}

func NewLocalFileService() *LocalFileService {
	return &LocalFileService{}
}

func C12sStateDir() string {
	dir := os.Getenv("NODE_ID_DIR_PATH")
	if dir == "" {
		dir = "/var/lib/c12s"
	}
	return dir
}

func C12sFilePath(fileName string) string {
	return filepath.Join(C12sStateDir(), fileName)
}
func (fs *LocalFileService) WriteToFile(filename string, data []byte) error {
	return fs.writeFile(filename, data, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
}

func (fs *LocalFileService) AppendToFile(filename string, data []byte) error {
	return fs.writeFile(filename, data, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

func (fs *LocalFileService) ReadFromFile(filename string) ([]byte, *errors.ErrorStruct) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return bytes, nil
}

func (fs *LocalFileService) writeFile(filename string, data []byte, flag int) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	return writer.Flush()
}
