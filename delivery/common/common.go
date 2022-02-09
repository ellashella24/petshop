package common

type ResponseSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Code:    200,
		Message: "Successful Operation",
		Data:    data,
	}
}

func ErrorResponse(code int, message string) ResponseError {
	return ResponseError{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
