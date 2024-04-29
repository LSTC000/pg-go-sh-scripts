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
	BashId                error
	BashFileUpload        error
	BashFileExtension     error
	BashGetFileBody       error
	BashFileTitle         error
	BashFileBody          error
	BashCreate            error
	BashDoesNotExists     error
	BashGetPaginationPage error
	BashExecuteIsSync     error
	BashExecuteDTOList    error
	BashExecute           error
	BashRemove            error

	// Bash Log Errors
	BashLogGetPaginationPageByBashId error

	// Pagination
	PaginationLimitParamMustBeInt  error
	PaginationLimitParamGTEZero    error
	PaginationOffsetParamMustBeInt error
	PaginationOffsetParamGTEZero   error
}

var (
	httpErrorsInstance *HTTPErrors
	httpErrorsOnce     sync.Once
)

func setHTTPErrors(errors *HTTPErrors) {
	// Base Errors
	errors.Internal = &schema.HTTPError{
		HTTPCode:    http.StatusInternalServerError,
		ServiceCode: 0,
		Detail:      "Internal Error",
	}

	// Pagination
	errors.PaginationLimitParamMustBeInt = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 100,
		Detail:      "The limit pagination parameter must be integer",
	}
	errors.PaginationLimitParamGTEZero = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 101,
		Detail:      "The limit pagination parameter must be greater than or equal to zero",
	}
	errors.PaginationOffsetParamMustBeInt = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 102,
		Detail:      "The offset pagination parameter must be integer",
	}
	errors.PaginationOffsetParamGTEZero = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 103,
		Detail:      "The offset pagination parameter must be greater than or equal to zero",
	}

	// Bash Errors
	errors.BashId = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 200,
		Detail:      "The bash id must be of type uuid4 like 151a583c-0ea0-46b8-b8a6-6bdcdd51655a",
	}
	errors.BashFileUpload = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 201,
		Detail:      "The bash script file has not been uploaded",
	}
	errors.BashFileExtension = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 202,
		Detail:      "The file extension should only be .sh",
	}
	errors.BashFileTitle = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 203,
		Detail:      "The file title should not be an empty string",
	}
	errors.BashFileBody = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 204,
		Detail:      "The file body should not be an empty string",
	}
	errors.BashGetFileBody = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 205,
		Detail:      "An error occurred while reading the body of the file",
	}
	errors.BashCreate = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 206,
		Detail:      "An error occurred during the creation of the bash script entity",
	}
	errors.BashDoesNotExists = &schema.HTTPError{
		HTTPCode:    http.StatusNotFound,
		ServiceCode: 207,
		Detail:      "The specified bash script does not exists",
	}
	errors.BashGetPaginationPage = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 208,
		Detail:      "An error occurred while receiving the pagination page of bash scripts",
	}
	errors.BashExecuteIsSync = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 209,
		Detail:      "The isSync executing parameter must be bool",
	}
	errors.BashExecuteDTOList = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 210,
		Detail:      "Invalid body of the request to start executing bash scripts",
	}
	errors.BashExecute = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 211,
		Detail:      "An error occurred while executing the bash script",
	}
	errors.BashRemove = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 212,
		Detail:      "An error occurred while deleting the bash script",
	}

	// Bash Log Errors
	errors.BashLogGetPaginationPageByBashId = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 300,
		Detail:      "An error occurred while receiving the pagination page of bash log scripts",
	}
}

func GetHTTPErrors() *HTTPErrors {
	httpErrorsOnce.Do(func() {
		var httpErrors HTTPErrors

		setHTTPErrors(&httpErrors)

		httpErrorsInstance = &httpErrors
	})

	return httpErrorsInstance
}
