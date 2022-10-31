package main

import (
	"errors"
	"os"
)

func createDirIfNotExists(path string) (err error) {
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0755)
	}
	return
}

func exists(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
