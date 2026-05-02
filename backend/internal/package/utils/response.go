package utils

type Response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"POKOKNYA SUKSES"`
	Error   any    `json:"error,omitempty" example:"null"`
	Data    any    `json:"data,omitempty"`
}

type ResponseErr struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Internal | Request | Unauthorized error"`
	Error   any    `json:"error,omitempty" example:"Detail errornya"`
	Data    any    `json:"data,omitempty" example:"null"`
}

type EmptyObj struct{}

func BuildResponseSuccess(message string, data any) Response {
	res := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildResponseFailed(message string, err string, data any) Response {
	res := Response{
		Success: false,
		Message: message,
		Error:   err,
		Data:    data,
	}
	return res
}
