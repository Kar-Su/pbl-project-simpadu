package dto

import (
	"errors"
	"web-hosting/internal/database/entities"

	"github.com/google/uuid"
)

const (
	MESSAGE_FAILED_ASSIGN_MAHASISWA_TO_KELAS   = "failed to assign mahasiswa to kelas"
	MESSAGE_FAILED_REMOVE_MAHASISWA_FROM_KELAS = "failed to remove mahasiswa from kelas"
	MESSAGE_FAILED_GET_DATA_PIVOT              = "failed to get data pivot"

	MESSAGE_SUCCESS_ASSIGN_MAHASISWA_TO_KELAS   = "success to assign mahasiswa to kelas"
	MESSAGE_SUCCESS_REMOVE_MAHASISWA_FROM_KELAS = "success to remove mahasiswa from kelas"
	MESSAGE_SUCCESS_GET_DATA_PIVOT              = "success to get data pivot"
)

var (
	ErrMahasiswaAlreadyAssigned = errors.New("mahasiswa already assigned to kelas")
	ErrMahasiswaNotFound        = errors.New("mahasiswa not found")
)

type (
	KelasResponseLess struct {
		KelasID uuid.UUID `json:"kelas_id"`
		Name    string    `json:"name"`
	}

	KelasMahasiswaResponse struct {
		Mahasiswa []MahasiswaResponse `json:"mahasiswa,omitempty"`
		Kelas     *KelasResponseLess  `json:"kelas,omitempty"`
	}
)

func ToKelasMahasiswaResponse(entity entities.KelasMahasiswa) KelasMahasiswaResponse {
	res := KelasMahasiswaResponse{}

	if entity.Kelas.ID != uuid.Nil {
		res.Kelas = &KelasResponseLess{
			KelasID: entity.Kelas.ID,
			Name:    entity.Kelas.Name,
		}
	}

	for _, mhs := range entity.Mahasiswa {
		res.Mahasiswa = append(res.Mahasiswa, MahasiswaResponse{
			MahasiswaID: mhs.DetailID,
			Name:        mhs.Name,
			Email:       mhs.Email,
		})
	}

	return res
}

type (
	KelasMahasiswaCreateRequest struct {
		KelasID     uuid.UUID `json:"kelas_id" binding:"required"`
		MahasiswaID uuid.UUID `json:"mahasiswa_id" binding:"required"`
	}
	MahasiswaIdURI struct {
		MahasiswaID uuid.UUID `uri:"mahasiswa_id" binding:"required"`
	}
	KelasIdURI struct {
		KelasID uuid.UUID `uri:"kelas_id" binding:"required"`
	}
	KelasMahasiswURI struct {
		MahasiswaIdURI
		KelasIdURI
	}
)
