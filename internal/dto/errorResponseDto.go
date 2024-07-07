package dto

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorsResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorResponseDto struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
