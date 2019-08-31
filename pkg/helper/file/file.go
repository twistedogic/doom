package file

import (
	"os"
)

func CreateIfNotExist(filename string) (*os.File, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else if os.IsNotExist(err) {
		return os.Create(filename)
	}
	return nil, err
}
