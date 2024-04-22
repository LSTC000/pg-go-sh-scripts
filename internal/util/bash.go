package util

import (
	"bufio"
	"bytes"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"pg-sh-scripts/internal/log"
	"strings"
)

const bashFileExtension = ".sh"

type (
	IBashUtil interface {
		ValidateBashFileExtension(string) bool
		GetBashFileExtension(string) string
		GetBashFileTitle(string) string
		GetBashFileBody(*multipart.FileHeader) (string, error)
		GetBashFileBuffer(string) *bytes.Buffer
	}

	BashUtil struct{}
)

func (u *BashUtil) ValidateBashFileExtension(fileExtension string) bool {
	return fileExtension == bashFileExtension
}

func (u *BashUtil) GetBashFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func (u *BashUtil) GetBashFileTitle(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (u *BashUtil) GetBashFileBody(file *multipart.FileHeader) (string, error) {
	var bashFileBody string

	f, err := file.Open()
	if err != nil {
		return bashFileBody, err
	}

	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			logger := log.GetLogger()
			logger.Error(fmt.Sprintf("Close file %s error: %v", file.Filename, err))
		}
	}(f)

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text() + "\n")
	}

	bashFileBody = wr.String()

	return bashFileBody, nil
}

func (u *BashUtil) GetBashFileBuffer(body string) *bytes.Buffer {
	return bytes.NewBufferString(body)
}

func GetBashUtil() IBashUtil {
	return &BashUtil{}
}
