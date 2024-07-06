package dto

type BaseResponse[T any] struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    T                      `json:"data"`
	Meta    map[string]interface{} `json:"_meta,omitempty"`
	Links   map[string]interface{} `json:"_link,omitempty"`
}

type DefaultApiResponse struct {
	BaseResponse[any] `json:",inline"`
}
