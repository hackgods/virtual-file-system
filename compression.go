package main

import (
	//"bytes"
	"compress/gzip"
	"os"
	"path/filepath"

	//"fmt"
	"io"
)

func Compress(data []byte, filename string) error {
	filePath := filepath.Join("storageData", filename)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	_, err = gzWriter.Write(data)
	if err != nil {
		return err
	}

	err = gzWriter.Flush()
	if err != nil {
		return err
	}

	err = gzWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func Decompress(filename string) ([]byte, error) {
	filePath := filepath.Join("storageData", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	data, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
