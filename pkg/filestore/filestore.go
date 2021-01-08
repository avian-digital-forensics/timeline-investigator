package filestore

import (
	"errors"
	"fmt"
	"os"
)

// Service handles the methods
// for the filestore
type Service interface {
	Upload(identifier, name string, content []byte) (*File, error)
	Delete(identifier, name string) error
}

type svc struct {
	path string
}

// New creates a new Service
func New(path string) (Service, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	return svc{path}, nil
}

// File is the object
// with file-information
type File struct {
	Name string
	Path string
	Size int
}

func (s svc) Upload(identifier, name string, content []byte) (*File, error) {
	if len(name) == 0 {
		return nil, errors.New("specify name for file to upload")
	}

	path := fmt.Sprintf("%s/%s/%s", s.path, identifier, name)
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	size, err := file.Write(content)
	if err != nil {
		os.Remove(path)
		return nil, err
	}

	return &File{Name: name, Path: path, Size: size}, nil
}

func (s svc) Delete(identifier, name string) error {
	return os.Remove(fmt.Sprintf("%s/%s/%s", s.path, identifier, name))
}
