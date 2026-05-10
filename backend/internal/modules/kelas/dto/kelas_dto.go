package dto

import (
	"errors"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/package/types"

	"github.com/google/uuid"
)

const (
	MESSAGE_FAILED_CREATE_KELAS = "failed to create kelas"
	MESSAGE_FAILED_UPDATE_KELAS = "failed to update kelas"
	MESSAGE_FAILED_DELETE_KELAS = "failed to delete kelas"
	MESSAGE_FAILED_GET_KELAS    = "failed to get kelas"

	MESSAGE_SUCCESS_CREATE_KELAS = "successfully created kelas"
	MESSAGE_SUCCESS_UPDATE_KELAS = "successfully updated kelas"
	MESSAGE_SUCCESS_DELETE_KELAS = "successfully deleted kelas"
	MESSAGE_SUCCESS_GET_KELAS    = "successfully get kelas"
)

var (
	ErrKelasNotFound         = errors.New("kelas not found")
	ErrKelasAlreadyExists    = errors.New("kelas already exists")
	ErrTahunAkademikNotFound = errors.New("tahun akademik not found")
	ErrProdiNotFound         = errors.New("prodi not found")
	ErrKurikulumNotFound     = errors.New("kurikulum not found")
)

type (
	KelasResponse struct {
		ID            uuid.UUID              `json:"id"`
		Name          string                 `json:"name"`
		Semester      uint                   `json:"semester"`
		TahunAkademik *TahunAkademikResponse `json:"tahun_akademik,omitempty"`
		Prodi         *ProdiResponse         `json:"prodi,omitempty"`
		Kurikulum     *KurikulumResponse     `json:"kurikulum,omitempty"`
		KurikulumMK   []KurikulumMKResponse  `json:"kurikulum_mk,omitempty"`
		Mahasiswa     []MahasiswaResponse    `json:"mahasiswa,omitempty"`
	}

	ProdiResponse struct {
		ID      uint             `json:"id"`
		Name    string           `json:"name"`
		Jenjang string           `json:"jenjang"`
		Jurusan *JurusanResponse `json:"jurusan,omitempty"`
	}

	JurusanResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	TahunAkademikResponse struct {
		ID           uint           `json:"id"`
		TipeSemester string         `json:"tipe_semester"`
		TahunAwal    types.DateOnly `json:"tahun_awal"`
		TahunAkhir   types.DateOnly `json:"tahun_akhir"`
		Status       string         `json:"status"`
	}

	KurikulumResponse struct {
		Kode string `json:"kode"`
		Name string `json:"name"`
	}

	KurikulumMKResponse struct {
		Semester   int                `json:"semester"`
		Wajib      bool               `json:"wajib"`
		MataKuliah MataKuliahResponse `json:"mata_kuliah"`
	}

	MataKuliahResponse struct {
		Kode string `json:"kode"`
		Name string `json:"name"`
		SKS  uint   `json:"sks"`
	}

	MahasiswaResponse struct {
		MahasiswaID *uuid.UUID `json:"mahasiswa_id"`
		Name        string     `json:"name"`
		Email       string     `json:"email"`
	}
)

func ToKelasResponse(entity entities.Kelas) KelasResponse {
	res := KelasResponse{
		ID:       entity.ID,
		Name:     entity.Name,
		Semester: entity.Semester,
	}

	if entity.Prodi.ID != 0 {
		res.Prodi = &ProdiResponse{
			ID:      entity.Prodi.ID,
			Name:    entity.Prodi.Name,
			Jenjang: entity.Prodi.Jenjang,
		}
		if entity.Prodi.Jurusan.ID != 0 {
			res.Prodi.Jurusan = &JurusanResponse{
				ID:   entity.Prodi.Jurusan.ID,
				Name: entity.Prodi.Jurusan.Name,
			}
		}
	}

	if entity.TahunAkademik.ID != 0 {
		res.TahunAkademik = &TahunAkademikResponse{
			ID:           entity.TahunAkademik.ID,
			TipeSemester: string(entity.TahunAkademik.TipeSemester),
			TahunAwal:    entity.TahunAkademik.TahunAwal,
			TahunAkhir:   entity.TahunAkademik.TahunAkhir,
			Status:       string(entity.TahunAkademik.Status),
		}
	}

	if entity.Kurikulum.Kode != "" {
		res.Kurikulum = &KurikulumResponse{
			Kode: entity.Kurikulum.Kode,
			Name: entity.Kurikulum.Name,
		}
	}

	for _, mk := range entity.Kurikulum.KurikulumMK {
		res.KurikulumMK = append(res.KurikulumMK, KurikulumMKResponse{
			Semester: mk.Semester,
			Wajib:    mk.Wajib,
			MataKuliah: MataKuliahResponse{
				Kode: mk.MataKuliah.Kode,
				Name: mk.MataKuliah.Name,
				SKS:  mk.MataKuliah.Sks,
			},
		})
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
	KelasCreateRequest struct {
		Name            string `json:"name" binding:"required,max=255"`
		Semester        uint   `json:"semester" binding:"omitempty,gte=0"`
		ProdiName       string `json:"prodi_name" binding:"required,max=255"`
		TahunAkademikID uint   `json:"tahun_akademik_id" binding:"required"`
		KurikulumKode   string `json:"kurikulum_kode" binding:"required,max=255"`
	}

	KelasUpdateRequest struct {
		Name            string `json:"name" binding:"omitempty,max=255"`
		Semester        uint   `json:"semester" binding:"omitempty,gte=0"`
		ProdiName       string `json:"prodi_name" binding:"omitempty,max=255"`
		TahunAkademikID uint   `json:"tahun_akademik_id" binding:"omitempty"`
		KurikulumKode   string `json:"kurikulum_kode" binding:"omitempty,max=255"`
	}

	ProdiNameURI struct {
		Name string `uri:"prodi_name" binding:"required,max=255"`
	}

	KelasIDURI struct {
		ID uuid.UUID `uri:"id" binding:"required"`
	}
)
