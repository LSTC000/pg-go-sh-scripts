package gosha

import (
	"fmt"
	"os"
	"path"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	tmpDir        = "tmp"
	bashExtension = ".sh"
)

func GetTmpFile(body string) (*os.File, error) {
	fileName := fmt.Sprintf("%v-%d%s", uuid.NewV4(), time.Now().Unix(), bashExtension)
	tmpPath := path.Join(tmpDir, fileName)

	if _, err := os.Stat(tmpDir); err != nil {
		if err := os.Mkdir(tmpDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	_, err := os.Create(tmpPath)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(tmpPath, os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.Write([]byte(body)); err != nil {
		return nil, err
	}

	return f, nil
}

func RemoveTmpFile(f *os.File) error {
	if err := os.Remove(f.Name()); err != nil {
		return err
	}
	return nil
}
