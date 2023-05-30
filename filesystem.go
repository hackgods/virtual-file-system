package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileSystem struct {
	BaseDir    string
	Versioning *Versioning // Added Versioning field
}

func NewFileSystem(baseDir string, versioning *Versioning) *FileSystem {
	return &FileSystem{
		BaseDir:    baseDir,
		Versioning: versioning, // Set the provided versioning object
	}
}

// UpdateBaseDir updates the base directory of the FileSystem.
func (fs *FileSystem) UpdateBaseDir(newBaseDir string) {
	fs.BaseDir = newBaseDir
}

func (fs *FileSystem) CreateFile(filename string, data []byte) error {
	filePath := filepath.Join(fs.BaseDir, filename)

	// Check if the file already exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		//return errors.New("file already exists")
	}

	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	// Perform versioning operation
	err = fs.Versioning.AddVersion(filename, data)
	if err != nil {
		fmt.Println("Error adding version:", err)
	}

	return nil
}

func (fs *FileSystem) ReadFile(name string) ([]byte, error) {
	path := filepath.Join(fs.BaseDir, name)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Read file: %s\n", name)
	return content, nil
}

func (fs *FileSystem) UpdateFile(name string, content []byte) error {

	path := filepath.Join(fs.BaseDir, name)
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	latestVersion, err := fs.Versioning.GetLatestVersion(name)
	if err != nil {
		return err
	}
	newVersion := latestVersion + 1

	// Add the new version to the versioning system
	err = fs.Versioning.AddVersion(name, content)
	if err != nil {
		return err
	}

	fmt.Printf("File %s updated successfully with version %d\n", name, newVersion)
	return nil
}

func (fs *FileSystem) DeleteFile(name string) error {
	path := filepath.Join(fs.BaseDir, name)
	err := os.Remove(path)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted file: %s\n", name)
	return nil
}
