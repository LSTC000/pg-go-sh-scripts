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

func GetBashFileBuffer(title, body string) (*bytes.Buffer, error) {
	if title == "" || body == "" {
		return nil, fmt.Errorf("title or body bash script cannot be empty")
	}

	fileBuffer := bytes.NewBufferString(body)

	return fileBuffer, nil
}

func GetBashFileTitle(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func GetBashFileBody(file *multipart.FileHeader) (string, error) {
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

func GetBashFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func ValidateBashFileExtension(fileExtension string) error {
	if fileExtension != bashFileExtension {
		return fmt.Errorf("invalid bash file extension: %s", fileExtension)
	}
	return nil
}
