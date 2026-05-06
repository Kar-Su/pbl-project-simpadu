package dto

import (
	"errors"
	"web-hosting/internal/database/entities"
)

const (
	MESSAGE_FAILED_GET_ALL        = "failed to get all jurusan"
	MESSAGE_FAILED_GET_JURUSAN    = "failed to get jurusan"
	MESSAGE_FAILED_CREATE_JURUSAN = "failed to create jurusan"
	MESSAGE_FAILED_UPDATE_JURUSAN = "failed to update jurusan"
	MESSAGE_FAILED_DELETE_JURUSAN = "failed to delete jurusan"

	MESSAGE_SUCCESS_GET_ALL        = "successfully get all jurusan"
	MESSAGE_SUCCESS_GET_JURUSAN    = "successfully get jurusan"
	MESSAGE_SUCCESS_CREATE_JURUSAN = "successfully create jurusan"
	MESSAGE_SUCCESS_UPDATE_JURUSAN = "successfully update jurusan"
	MESSAGE_SUCCESS_DELETE_JURUSAN = "successfully delete jurusan"
)

var (
	ErrJurusanNotFound      = errors.New("jurusan not found")
	ErrJurusanAlreadyExists = errors.New("jurusan already exists")
)

type (
	JurusanResponse struct {
		ID   uint   `json:"id" example:"1"`
		Name string `json:"name" example:"teknik-elektro"`
	}

	JurusanRequest struct {
		JurusanName string `json:"jurusan_name" binding:"required" example:"teknik-elektro"`
	}

	JurusanUpdateRequest struct {
		NewName string `json:"new_name" binding:"required" example:"teknik-mesin"`
	}

	JurusanNameQuery struct {
		JurusanName string `form:"name" binding:"required" example:"teknik-elektro"`
	}

	JurusanQuery struct {
		JurusanName string `form:"name" example:"teknik-elektro"`
		JurusanID   uint   `form:"id" example:"1"`
	}
)

func ToJurusanResponse(jurusanEntity entities.Jurusan) JurusanResponse {
	return JurusanResponse{
		ID:   jurusanEntity.ID,
		Name: jurusanEntity.Name,
	}
}
