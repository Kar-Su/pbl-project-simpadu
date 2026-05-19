package dto

import (
	"errors"
	"web-hosting/internal/database/entities"

	"github.com/google/uuid"
)

const (
	FAILED_CREATE_PENGAMPU = "failed to create pengampu"
	FAILED_UPDATE_PENGAMPU = "failed to update pengampu"
	FAILED_DELETE_PENGAMPU = "failed to delete pengampu"
	FAILED_GET_PENGAMPU    = "failed to get pengampu"

	SUCCESS_CREATE_PENGAMPU = "success to create pengampu"
	SUCCESS_UPDATE_PENGAMPU = "success to update pengampu"
	SUCCESS_DELETE_PENGAMPU = "success to delete pengampu"
	SUCCESS_GET_PENGAMPU    = "success to get pengampu"
)

var (
	ErrPengampuNotFound = errors.New("pengampu not found")
	ErrAlreadyExists    = errors.New("pengampu already exists")
	ErrDuplicatedKey    = errors.New("duplicate key error")
)

type (
	PengampuResponse struct {
		PengampuID uuid.UUID           `json:"pengampu_id"`
		Dosen      *DosenResponse      `json:"dosen"`
		MataKuliah *MataKuliahResponse `json:"mata_kuliah"`
	}

	DosenResponse struct {
		DosenID uuid.UUID `json:"id"`
		Name    string    `json:"name"`
		Email   string    `json:"email"`
	}
	MataKuliahResponse struct {
		MKID   uuid.UUID `json:"id"`
		MKKode string    `json:"kode"`
		Name   string    `json:"name"`
		Sks    uint      `json:"sks"`
	}
)

func ToPengampuResponse(entity entities.Pengampu) PengampuResponse {
	res := PengampuResponse{
		PengampuID: entity.ID,
	}

	res.Dosen = &DosenResponse{
		DosenID: entity.DosenID,
		Name:    entity.Dosen.Name,
		Email:   entity.Dosen.Email,
	}

	res.MataKuliah = &MataKuliahResponse{
		MKID:   entity.MataKuliah.ID,
		MKKode: entity.MataKuliah.Kode,
		Name:   entity.MataKuliah.Name,
		Sks:    entity.MataKuliah.Sks,
	}

	return res
}

type (
	CreatePengampuRequest struct {
		KelasID uuid.UUID `json:"kelas_id" binding:"required"`
		MKKode  string    `json:"mkkode" binding:"required,max=12"`
		DosenID uuid.UUID `json:"dosen_id" binding:"required"`
	}
	UpdatePengampuRequest struct {
		KelasID *uuid.UUID `json:"kelas_id" binding:"omitempty"`
		MKKode  string     `json:"mkkode" binding:"omitempty,max=12"`
		DosenID *uuid.UUID `json:"dosen_id" binding:"omitempty"`
	}
)

type (
	PengampuIdURI struct {
		PengampuID string `uri:"pengampu_id" binding:"required"`
	}

	PengampuKelasIdURI struct {
		KelasID string `uri:"kelas_id" binding:"required"`
	}
)
