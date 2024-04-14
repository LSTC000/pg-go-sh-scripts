package config

import (
	"net/http"
	"pg-sh-scripts/internal/schema"
	"sync"
)

type HTTPErrors struct {
	// Base Errors
	Internal error
	Validate error
	// Bash Errors
	BashFileExtension error
	BashFileBody      error
	BashCreate        error
	BashGet           error
	BashGetFile       error
	BashGetFilePath   error
	BashGetList       error
	BashExecute       error
	BashExecuteList   error
	BashRemove        error
	// Bash Log Errors
	BashLogGetListByBashId error
}

func setHTTPErrors(errors *HTTPErrors) {
	// Base Errors
	errors.Internal = &schema.HTTPError{
		HTTPCode:    http.StatusInternalServerError,
		ServiceCode: 0,
		Detail:      "Internal Error",
	}
	errors.Validate = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 1,
		Detail:      "Validation Error",
	}
	// Bash Errors
	errors.BashFileExtension = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 100,
		Detail:      "Invalid Bash File Extension",
	}
	errors.BashFileBody = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 101,
		Detail:      "Invalid Bash File Body",
	}
	errors.BashCreate = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 102,
		Detail:      "Creating Bash Error",
	}
	errors.BashGet = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 103,
		Detail:      "Getting Bash Error",
	}
	errors.BashGetFile = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 104,
		Detail:      "Getting Bash File Error",
	}
	errors.BashGetFilePath = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 105,
		Detail:      "Getting Bash File Path Error",
	}
	errors.BashGetList = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 106,
		Detail:      "Getting Bash List Error",
	}
	errors.BashExecute = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 107,
		Detail:      "Executing Bash Error",
	}
	errors.BashExecuteList = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 108,
		Detail:      "Executing Bash List Error",
	}
	errors.BashRemove = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 109,
		Detail:      "Removing Bash Error",
	}
	// Bash Log Errors
	errors.BashLogGetListByBashId = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 108,
		Detail:      "Getting Bash Log List By Bash Id Error",
	}
}

func GetHTTPErrors() *HTTPErrors {
	var (
		errors HTTPErrors
		once   sync.Once
	)

	once.Do(func() {
		setHTTPErrors(&errors)
	})

	return &errors
}
