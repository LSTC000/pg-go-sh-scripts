package schema

type (
	Message struct {
		Message string `json:"message"`
	}

	HTTPError struct {
		HTTPCode    int    `json:"httpCode"`
		ServiceCode int    `json:"serviceCode"`
		Detail      string `json:"detail"`
	}
)
