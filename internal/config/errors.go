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
	BashFileBody          error
	BashCreate            error
	BashDoesNotExists     error
	BashGetFileBuffer     error
	BashGetPaginationPage error
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
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 201,
		Detail:      "The bash script file has not been uploaded",
	}
	errors.BashFileExtension = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 202,
		Detail:      "The file extension should only be .sh",
	}
	errors.BashFileBody = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 203,
		Detail:      "An error occurred while reading the body of the file",
	}
	errors.BashCreate = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 204,
		Detail:      "An error occurred during the creation of the bash script entity",
	}
	errors.BashDoesNotExists = &schema.HTTPError{
		HTTPCode:    http.StatusNotFound,
		ServiceCode: 205,
		Detail:      "The specified bash script does not exists",
	}
	errors.BashGetFileBuffer = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 206,
		Detail:      "An error occurred while retrieving the content buffer of the bash script",
	}
	errors.BashGetPaginationPage = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 207,
		Detail:      "An error occurred while receiving the pagination page of bash scripts",
	}
	errors.BashExecuteDTOList = &schema.HTTPError{
		HTTPCode:    http.StatusUnprocessableEntity,
		ServiceCode: 208,
		Detail:      "Invalid body of the request to start executing bash scripts",
	}
	errors.BashExecute = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 209,
		Detail:      "An error occurred while executing the bash script",
	}
	errors.BashRemove = &schema.HTTPError{
		HTTPCode:    http.StatusBadRequest,
		ServiceCode: 210,
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
