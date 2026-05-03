package utils

type Response[T any, E any] struct {
	Success bool   `json:"success"`
	Message string `json:"message" example:"Operation successful"`
	Path    string `json:"path,omitempty" example:"/api/resource"`
	Error   E      `json:"error,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func BuildResponseSuccess[T any](message string, data T, path ...string) Response[T, any] {
	var p string
	if len(path) > 0 {
		p = path[0]
	}
	return Response[T, any]{
		Success: true,
		Message: message,
		Path:    p,
		Data:    data,
	}
}

func BuildResponseFailed[E any](message string, err E, data any, path ...string) Response[any, E] {
	var p string
	if len(path) > 0 {
		p = path[0]
	}
	return Response[any, E]{
		Success: false,
		Message: message,
		Path:    p,
		Error:   err,
	}
}
