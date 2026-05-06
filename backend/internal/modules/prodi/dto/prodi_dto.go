package dto

import (
	"errors"
	"web-hosting/internal/database/entities"
)

const (
	MESSAGE_FAILED_CREATE_PRODI = "Failed to create prodi"
	MESSAGE_FAILED_UPDATE_PRODI = "Failed to update prodi"
	MESSAGE_FAILED_DELETE_PRODI = "Failed to delete prodi"
	MESSAGE_FAILED_GET_PRODI    = "Failed to get prodi"

	MESSAGE_SUCCESS_CREATE_PRODI = "Successfully created prodi"
	MESSAGE_SUCCESS_UPDATE_PRODI = "Successfully updated prodi"
	MESSAGE_SUCCESS_DELETE_PRODI = "Successfully deleted prodi"
	MESSAGE_SUCCESS_GET_PRODI    = "Successfully got prodi"
)

var (
	ErrProdiNotFound      = errors.New("prodi not found")
	ErrProdiAlreadyExists = errors.New("prodi already exists")
)

type (
	ProdiResponse struct {
		ID      uint             `json:"id" example:"1"`
		Name    string           `json:"name" example:"teknik-elektro"`
		Jenjang string           `json:"jenjang" example:"D3"`
		Jurusan entities.Jurusan `json:"jurusan"`
	}

	ProdiCreateRequest struct {
		Name        string `json:"name" binding:"required" example:"teknik-elektro"`
		Jenjang     string `json:"jenjang" binding:"required,enum_jenjang" example:"D3" enums:"D3,D4"`
		JurusanName string `json:"jurusan_name" binding:"required" example:"teknik-elektro"`
	}

	ProdiUpdateRequest struct {
		Name      string `json:"name" binding:"omitempty" example:"teknik-elektro"`
		Jenjang   string `json:"jenjang" binding:"omitempty,enum_jenjang" example:"D3"`
		JurusanID uint   `json:"jurusan_id" binding:"omitempty" example:"1"`
	}

	ProdiQuery struct {
		ID   uint   `form:"id" example:"1"`
		Name string `form:"name" example:"teknik-elektro"`
	}

	ProdiNameQuery struct {
		Name string `form:"name" binding:"required" example:"teknik-elektro"`
	}

	JurusanURI struct {
		JurusanName string `uri:"jurusan_name" binding:"required" example:"teknik-elektro"`
	}
)

func ToProdiResponse(prodi entities.Prodi) ProdiResponse {
	return ProdiResponse{
		ID:      prodi.ID,
		Name:    prodi.Name,
		Jenjang: prodi.Jenjang,
		Jurusan: prodi.Jurusan,
	}
}
