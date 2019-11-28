package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileSystemAliasService map[string]string

func (f FileSystemAliasService) GetAlias(email string) string {
	if alias, found := f[email]; found {
		return alias
	}
	return email
}

func NewFileSystemAliasService(path string) (FileSystemAliasService, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("couldn't open %s, %s", path, err.Error())
	}

	decoder := json.NewDecoder(file)

	service := make(map[string]string)

	err = decoder.Decode(&service)

	if err != nil {
		return nil, fmt.Errorf("couldn't decode contents of %s, %s", path, err.Error())
	}

	return service, nil
}

type NoOpAliasService struct {

}

func (n NoOpAliasService) GetAlias(email string) string {
	return email
}
