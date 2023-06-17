package dto

type BaseResponse struct {
	Code int `json:"code"`
}

type BaseResponseWithData[T interface{}] struct {
	BaseResponse
	Data T `json:"data"`
}

type BaseErrorResponse struct {
	BaseResponse
	Message string `json:"message"`
}

type BaseOKResponse struct {
	Code    int  `json:"code"`
	Success bool `json:"success"`
}
type BaseOKResponseWithData[T any] struct {
	BaseOKResponse
	Data T `json:"data"`
}
