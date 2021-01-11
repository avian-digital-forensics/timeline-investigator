package filestore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Service handles the methods
// for the filestore
type Service interface {
	Upload(identifier, name string, content []byte) (*File, error)
	Delete(identifier, name string) error
	GetContent(filename string) ([]byte, error)
}

type svc struct {
	path string
}

// New creates a new Service
func New(path string) (Service, error) {
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	if strings.HasSuffix(path, "\\") {
		path = strings.TrimSuffix(path, "\\")
	}

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

	// Create the file and handle the errors
	path := fmt.Sprintf("%s/%s/%s", s.path, identifier, name)
	file, err := os.Create(path)
	if err != nil {
		// Return the error if it isn't a path error
		if _, ok := err.(*os.PathError); !ok {
			return nil, err
		}

		// Create a directory for the identifier
		if err := os.Mkdir(fmt.Sprintf("%s/%s", s.path, identifier), os.ModePerm); err != nil {
			return nil, fmt.Errorf("Failed to create new dir: %s", err.Error())
		}

		// Try to create the file again
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	size, err := file.Write(content)
	if err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("Failed to write data to file: %s", err.Error())
	}

	return &File{Name: name, Path: path, Size: size}, nil
}

func (s svc) Delete(identifier, name string) error {
	return os.Remove(fmt.Sprintf("%s/%s/%s", s.path, identifier, name))
}

func (svc) GetContent(filename string) ([]byte, error) { return ioutil.ReadFile(filename) }
