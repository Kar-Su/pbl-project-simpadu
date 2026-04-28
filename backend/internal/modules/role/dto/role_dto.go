package dto

import "errors"

const (
	MESSAGE_FAILED_GET_ROLE            = "failed to get role"
	MESSAGE_FAILED_GET_ALL_ROLE        = "failed to get all role"
	MESSAGE_ROLE_CREATE_ALREADY_EXISTS = "failed to create role, role already exists"
	MESSAGE_FAILED_DELETE_ROLE         = "failed to delete role"
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
		RoleName string `json:"role_name" binding:"required,role_exist"`
	}

	RoleNameURI struct {
		RoleName string `uri:"role_name" binding:"required, role_exist"`
	}
)
