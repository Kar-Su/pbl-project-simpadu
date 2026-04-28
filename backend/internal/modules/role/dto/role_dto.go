package dto

import "errors"

const (
	MESSAGE_FAILED_GET_ROLE          = "failed to get role"
	MESSAGE_FAILED_DELETE_ROLE       = "failed to delete role"
	MESSAGE_FAILED_GET_REQUEST_BODY  = "failed to get request"
	MESSAGE_FAILED_VALIDATE_ROLE_URI = "failed to validate role uri"
	MESSAGE_FAILED_CREATE_ROLE       = "failed to create role"
	MESSAGE_FAILED_UPDATE_ROLE       = "failed to update role"

	MESSAGE_SUCCESS_CREATE_ROLE = "role created successfully"
	MESSAGE_SUCCESS_UPDATE_ROLE = "role updated successfully"
	MESSAGE_SUCCESS_DELETE_ROLE = "role deleted successfully"
	MESSAGE_SUCCESS_GET_ROLE    = "role retrieved successfully"
)

var (
	ErrRoleNotFound     = errors.New("role not found")
	ErrRoleAlreadyExist = errors.New("role already exists")
)

type (
	RoleResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	RoleCreateRequest struct {
		RoleName string `json:"role_name" binding:"required"`
	}

	RoleUpdateRequest struct {
		RoleName string `json:"role_name" binding:"required"`
	}

	RoleNameURI struct {
		RoleName string `uri:"role_name" binding:"required,"`
	}
)
