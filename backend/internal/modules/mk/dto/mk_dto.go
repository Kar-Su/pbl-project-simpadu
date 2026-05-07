package dto

import (
	"errors"
	"web-hosting/internal/database/entities"
)

const (
	MESSAGE_FAILED_CREATE_MK = "failed to create mk"
	MESSAGE_FAILED_UPDATE_MK = "failed to update mk"
	MESSAGE_FAILED_DELETE_MK = "failed to delete mk"
	MESSAGE_FAILED_GET_MK    = "failed to get mk"

	MESSAGE_SUCCESS_CREATE_MK = "successfully created mk"
	MESSAGE_SUCCESS_UPDATE_MK = "successfully updated mk"
	MESSAGE_SUCCESS_DELETE_MK = "successfully deleted mk"
	MESSAGE_SUCCESS_GET_MK    = "successfully got mk"
)

var (
	ErrMkNotFound      = errors.New("Mata Kuliah Not Found")
	ErrMkAlreadyExists = errors.New("Mata Kuliah Already Exists")
	ErrQueryParams     = errors.New("no query params provided")
)

type (
	MkResponse struct {
		ID   string `json:"id" example:"12345678-1234-1234-1234-123456789012"`
		Kode string `json:"kode" example:"MK001"`
		Name string `json:"name" example:"VM DOCKER"`
		Sks  uint   `json:"sks" example:"3"`
	}

	MkCreateRequest struct {
		Kode string `json:"kode" binding:"required,max=12" example:"MK001"`
		Name string `json:"name" binding:"required,max=255" example:"VM DOCKER"`
		Sks  uint   `json:"sks" binding:"required,gt=0" example:"3"`
	}

	MkUpdateRequest struct {
		Kode string `json:"kode" binding:"omitempty,max=12" example:"MK001"`
		Name string `json:"name" binding:"omitempty,max=255" example:"VM DOCKER"`
		Sks  uint   `json:"sks" binding:"omitempty,gte=0" example:"3"`
	}

	MkQuery struct {
		ID   string `form:"id" example:"12345678-1234-1234-1234-123456789012"`
		Kode string `form:"kode" binding:"max=12" example:"MK001"`
	}
)

func ToMkResponse(entity entities.MataKuliah) MkResponse {
	return MkResponse{
		ID:   entity.ID.String(),
		Kode: entity.Kode,
		Name: entity.Name,
		Sks:  entity.Sks,
	}
}
