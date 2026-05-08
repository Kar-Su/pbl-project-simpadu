package dto

import (
	"errors"
	"web-hosting/internal/database/entities"
)

const (
	MESSAGE_FAILED_CREATE_PIVOT_KURIKULUM_MK = "failed to create pivot kurikulum-mk"
	MESSAGE_FAILED_UPDATE_PIVOT_KURIKULUM_MK = "failed to update pivot kurikulum-mk"
	MESSAGE_FAILED_DELETE_PIVOT_KURIKULUM_MK = "failed to delete pivot kurikulum-mk"

	MESSAGE_SUCCESS_CREATE_PIVOT_KURIKULUM_MK = "success to create pivot kurikulum-mk"
	MESSAGE_SUCCESS_UPDATE_PIVOT_KURIKULUM_MK = "success to update pivot kurikulum-mk"
	MESSAGE_SUCCESS_DELETE_PIVOT_KURIKULUM_MK = "success to delete pivot kurikulum-mk"
)

var (
	ErrMkNotFound         = errors.New("mk not found")
	ErrPivotNotFound      = errors.New("pivot not found")
	ErrPivotAlreadyExists = errors.New("pivot already exists")
)

type (
	PivotResponse struct {
		KurikulumKode string              `json:"kurikulum_kode"`
		Semester      int                 `json:"semester"`
		Wajib         bool                `json:"wajib"`
		MataKuliah    entities.MataKuliah `json:"mata_kuliah"`
	}

	PivotCreateRequest struct {
		KurikulumKode string `json:"kurikulum_kode" binding:"required,max=12"`
		MkKode        string `json:"mk_kode" binding:"required,max=12"`
		Semester      int    `json:"semester" binding:"required,gte=0"`
		Wajib         bool   `json:"wajib" binding:"required"`
	}

	PivotUpdateRequest struct {
		Semester int  `json:"semester" binding:"omitempty,gte=0"`
		Wajib    bool `json:"wajib" binding:"omitempty"`
	}

	PivotURI struct {
		KurikulumKode string `uri:"kurikulum_kode" binding:"required,max=12 example:kur-2024"`
		MkKode        string `uri:"mk_kode" binding:"required,max=12 example:MK001"`
	}
)

func ToPivotResponse(pivot entities.KurikulumMK) PivotResponse {
	return PivotResponse{
		KurikulumKode: pivot.KurikulumKode,
		Semester:      pivot.Semester,
		Wajib:         pivot.Wajib,
		MataKuliah:    pivot.MataKuliah,
	}
}
