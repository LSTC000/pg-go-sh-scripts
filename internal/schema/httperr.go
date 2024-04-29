package schema

import "fmt"

type HTTPError struct {
	HTTPCode    int    `json:"httpCode"`
	ServiceCode int    `json:"serviceCode"`
	Detail      string `json:"detail"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf(
		"HTTP Code: %d, Service code: %d, Detail: %s",
		e.HTTPCode,
		e.ServiceCode,
		e.Detail,
	)
}
