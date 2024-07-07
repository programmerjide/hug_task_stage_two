package dto

type BaseResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type DefaultApiResponse struct {
	BaseResponse[any] `json:",inline"`
}
