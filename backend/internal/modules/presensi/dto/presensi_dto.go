package dto

import (
	"web-hosting/internal/database/entities"
	"web-hosting/internal/package/types"

	"github.com/google/uuid"
)

const (
	FAILED_CREATE_PRESENSI_MAHASISWA = "failed to create presensi mahasiswa"
	FAILED_CREATE_PRESENSI_PEGAWAI   = "failed to create presensi pegawai"
	FAILED_UPDATE_DETAIL_PRESENSI    = "failed to update detail presensi"
	FAILED_GET_PRESENSI              = "failed to get presensi"

	SUCCESS_CREATE_PRESENSI_MAHASISWA = "success to create presensi mahasiswa"
	SUCCESS_CREATE_PRESENSI_PEGAWAI   = "success to create presensi pegawai"
	SUCCESS_UPDATE_DETAIL_PRESENSI    = "success to update detail presensi"
	SUCCESS_GET_PRESENSI              = "success to get presensi"
)

type (
	PresensiMahasiswaCreateRequest struct {
		ID         uuid.UUID `json:"sesi_id" binding:"required"`
		PengampuID uuid.UUID `json:"pengampu_id" binding:"required"`
	}

	PresensiMahasiswaUpdateRequest struct {
		PresensiID uuid.UUID                     `json:"sesi_id" binding:"required" example:"01965a1d-7777-7777-7777-777777777777"`
		Detail     []DetailPresensiUpdateRequest `json:"detail" binding:"required"`
	}

	DetailPresensiUpdateRequest struct {
		DetailID uuid.UUID `json:"detail_id" binding:"required" example:"mahasiswa/pegawai ID"`
		Status   string    `json:"status" binding:"required" example:"hadir/sakit/izin/alpha"`
	}

	PresensiPegawaiUpdateRequest struct {
		Date   types.DateOnly                `json:"date" binding:"required" example:"YYYY-MM-DD"`
		Detail []DetailPresensiUpdateRequest `json:"detail" binding:"required"`
	}
)

type (
	PresensiMahasiswaResponse struct {
		ID         uuid.UUID      `json:"sesi_id"`
		PengampuID uuid.UUID      `json:"pengampu_id"`
		Mahasiswa  []UserResponse `json:"mahasiswa"`
	}

	PresensiPegawaiResponse struct {
		ID      uuid.UUID      `json:"sesi_id"`
		Pegawai []UserResponse `json:"pegawai"`
	}

	UserResponse struct {
		ID     uuid.UUID `json:"detail_id"`
		Name   string    `json:"name"`
		Email  string    `json:"email"`
		Status string    `json:"status"`
	}
)

func ToPresensiResponse(tipePresensi string, entity entities.Presensi) any {
	switch tipePresensi {
	case "mahasiswa":
		response := PresensiMahasiswaResponse{
			ID:         entity.ID,
			PengampuID: *entity.PengampuID,
		}
		for _, p := range entity.PresensiMahasiswa {
			response.Mahasiswa = append(response.Mahasiswa, UserResponse{
				ID:     *p.Mahasiswa.DetailID,
				Name:   p.Mahasiswa.Name,
				Email:  p.Mahasiswa.Email,
				Status: p.Status,
			})
		}
		return response
	case "pegawai":
		response := PresensiPegawaiResponse{
			ID: entity.ID,
		}
		for _, p := range entity.PresensiPegawai {
			response.Pegawai = append(response.Pegawai, UserResponse{
				ID:     *p.Pegawai.DetailID,
				Name:   p.Pegawai.Name,
				Email:  p.Pegawai.Email,
				Status: p.Status,
			})
		}
		return response
	default:
		return nil
	}
}

type (
	GetPresensiMahasiswaQuery struct {
		ID string `form:"sesi_id" binding:"required" example:"SESI_ID"`
	}

	GetPresensiPegawaiQuery struct {
		ID   string `form:"sesi_id" binding:"omitempty" example:"PRESENSI_ID"`
		Date string `form:"date" binding:"omitempty" example:"YYYY-MM-DD"`
	}

	UpdatePresensiByQRQuery struct {
		SesiID string `form:"sesi_id" binding:"required" example:"SESI_ID"`
	}
)
