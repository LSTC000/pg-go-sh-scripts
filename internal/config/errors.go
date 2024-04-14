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
	BashFileExtension     error
	BashFileBody          error
	BashCreate            error
	BashDoesNotExists     error
	BashGetFileBuffer     error
	BashGetPaginationPage error
	BashExecute           error
	BashRemove            error
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
	errors.BashDoesNotExists = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 103,
		Detail:      "Bash Does Not Exists",
	}
	errors.BashGetFileBuffer = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 104,
		Detail:      "Getting Bash File Buffer Error",
	}
	errors.BashGetPaginationPage = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 105,
		Detail:      "Getting Bash Pagination Page Error",
	}
	errors.BashExecute = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 106,
		Detail:      "Executing Bash Error",
	}
	errors.BashRemove = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 107,
		Detail:      "Removing Bash Error",
	}
	// Bash Log Errors
	errors.BashLogGetListByBashId = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 200,
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
