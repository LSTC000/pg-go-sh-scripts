package config

import (
	"pg-sh-scripts/internal/model"
	"sync"
)

type HTTPErrors struct {
	// Base Errors
	Internal model.HTTPError
	Validate model.HTTPError
	// Bash Errors
	BashFileExtension model.HTTPError
	BashFileBody      model.HTTPError
	BashCreate        model.HTTPError
	BashGet           model.HTTPError
	BashGetFile       model.HTTPError
	BashGetFilePath   model.HTTPError
	BashGetList       model.HTTPError
	BashExecute       model.HTTPError
	BashExecuteList   model.HTTPError
	// Bash Log Errors
	BashLogGetListByBashID model.HTTPError
}

func setHTTPErrors(errors *HTTPErrors) {
	// Base Errors
	errors.Internal = model.HTTPError{
		HTTPCode:    500,
		ServiceCode: 0,
		Detail:      "Internal Error",
	}
	errors.Validate = model.HTTPError{
		HTTPCode:    422,
		ServiceCode: 1,
		Detail:      "Validation Error",
	}
	// Bash Errors
	errors.BashFileExtension = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 100,
		Detail:      "Invalid Bash File Extension",
	}
	errors.BashFileBody = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 101,
		Detail:      "Invalid Bash File Body",
	}
	errors.BashCreate = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 102,
		Detail:      "Creating Bash Error",
	}
	errors.BashGet = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 103,
		Detail:      "Getting Bash Error",
	}
	errors.BashGetFile = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 104,
		Detail:      "Getting Bash File Error",
	}
	errors.BashGetFilePath = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 105,
		Detail:      "Getting Bash File Path Error",
	}
	errors.BashGetList = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 106,
		Detail:      "Getting Bash List Error",
	}
	errors.BashExecute = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 107,
		Detail:      "Executing Bash Error",
	}
	errors.BashExecuteList = model.HTTPError{
		HTTPCode:    400,
		ServiceCode: 108,
		Detail:      "Executing Bash List Error",
	}
	// Bash Log Errors
	errors.BashLogGetListByBashID = model.HTTPError{
		HTTPCode:    400,
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
