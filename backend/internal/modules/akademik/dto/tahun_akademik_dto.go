package dto

import (
	"errors"
	"web-hosting/internal/database/entities"
	karTypes "web-hosting/internal/package/types"
)

const (
	MESSAGE_FAILED_CREATE_TAHUN_AKADEMIK = "failed to create tahun akademik"
	MESSAGE_FAILED_UPDATE_TAHUN_AKADEMIK = "failed to update tahun akademik"
	MESSAGE_FAILED_DELETE_TAHUN_AKADEMIK = "failed to delete tahun akademik"
	MESSAGE_FAILED_GET_TAHUN_AKADEMIK    = "failed to get tahun akademik"

	MESSAGE_SUCCESS_CREATE_TAHUN_AKADEMIK = "success to create tahun akademik"
	MESSAGE_SUCCESS_UPDATE_TAHUN_AKADEMIK = "success to update tahun akademik"
	MESSAGE_SUCCESS_DELETE_TAHUN_AKADEMIK = "success to delete tahun akademik"
	MESSAGE_SUCCESS_GET_TAHUN_AKADEMIK    = "success to get tahun akademik"
)

var (
	ErrTahunAkademikNotFound      = errors.New("tahun akademik not found")
	ErrTahunAkademikAlreadyExists = errors.New("tahun akademik already exists")
	ErrInvalidTahunAkademik       = errors.New("invalid tahun awal/akhir akademik")
)

type (
	AkademikResponse struct {
		ID           uint              `json:"id" example:"20241"`
		TipeSemester string            `json:"tipe_semester" example:"semester"`
		TahunAwal    karTypes.DateOnly `json:"tahun_awal" example:"2024-01-01"`
		TahunAkhir   karTypes.DateOnly `json:"tahun_akhir" example:"2025-01-01"`
		Status       string            `json:"status" example:"aktif"`
	}

	AkademikCreateRequest struct {
		ID           uint              `json:"id" binding:"required,gte=0" example:"20241"`
		TipeSemester string            `json:"tipe_semester" binding:"required,enumTipeSemester" example:"semester"`
		TahunAwal    karTypes.DateOnly `json:"tahun_awal" binding:"required" example:"2024-01-01"`
		TahunAkhir   karTypes.DateOnly `json:"tahun_akhir" binding:"required" example:"2025-01-01"`
	}

	AkademikUpdateRequest struct {
		ID           uint              `json:"id" binding:"omitempty,gte=0" example:"20241"`
		TipeSemester string            `json:"tipe_semester" binding:"omitempty,enumTipeSemester" example:"semester"`
		TahunAwal    karTypes.DateOnly `json:"tahun_awal" binding:"omitempty" example:"2024-01-01"`
		TahunAkhir   karTypes.DateOnly `json:"tahun_akhir" binding:"omitempty" example:"2025-01-01"`
		Status       string            `json:"status" binding:"omitempty,enumStatus" example:"aktif"`
	}
	AkademikIdUri struct {
		ID uint `uri:"id" binding:"required,gte=0" example:"20241"`
	}

	AkademikStatusUri struct {
		Status string `uri:"status" binding:"required,enumStatus" example:"aktif"`
	}
)

func ToAkademikResponse(entity entities.TahunAkademik) AkademikResponse {
	return AkademikResponse{
		ID:           entity.ID,
		TipeSemester: entity.TipeSemester,
		TahunAwal:    entity.TahunAwal,
		TahunAkhir:   entity.TahunAkhir,
		Status:       entity.Status,
	}
}
