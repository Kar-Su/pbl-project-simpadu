package dto

import (
	"errors"
	"web-hosting/internal/database/entities"
)

const (
	MESSAGE_FAILED_CREATE_KURIKULUM = "failed to create kurikulum"
	MESSAGE_FAILED_UPDATE_KURIKULUM = "failed to update kurikulum"
	MESSAGE_FAILED_DELETE_KURIKULUM = "failed to delete kurikulum"
	MESSAGE_FAILED_GET_KURIKULUM    = "failed to get kurikulum"

	MESSAGE_SUCCESS_CREATE_KURIKULUM = "success to create kurikulum"
	MESSAGE_SUCCESS_UPDATE_KURIKULUM = "success to update kurikulum"
	MESSAGE_SUCCESS_DELETE_KURIKULUM = "success to delete kurikulum"
	MESSAGE_SUCCESS_GET_KURIKULUM    = "success to get kurikulum"
)

var (
	ErrKurikulumNotFound      = errors.New("kurikulum not found")
	ErrKurikulumAlreadyExists = errors.New("kurikulum already exists")
	ErrUpdateViolatesUnique   = errors.New("update violates unique constraint")
	ErrQueryParams            = errors.New("query params are invalid, only choice id or kode is provided")
)

type (
	KurikulumResponse struct {
		ID          string                `json:"id" example:"0000-0000"`
		Kode        string                `json:"kode" example:"my-hutao"`
		Name        string                `json:"name" example:"My Hutao"`
		Prodi       ProdiResponse         `json:"prodi"`
		KurikulumMK []KurikulumMKResponse `json:"kurikulum_mk"`
	}

	ProdiResponse struct {
		ID      uint            `json:"id" example:"1"`
		Name    string          `json:"name" example:"teknik informatika"`
		Jenjang string          `json:"jenjang" example:"D3"`
		Jurusan JurusanResponse `json:"jurusan"`
	}

	JurusanResponse struct {
		ID   uint   `json:"id" example:"1"`
		Name string `json:"name" example:"teknik-elektro"`
	}

	KurikulumMKResponse struct {
		Semester   int                `json:"semester" example:"3"`
		Wajib      bool               `json:"wajib" example:"true"`
		MataKuliah MataKuliahResponse `json:"mata_kuliah"`
	}

	MataKuliahResponse struct {
		ID   string `json:"id" example:"MK001"`
		Kode string `json:"kode" example:"MK001"`
		Name string `json:"name" example:"Pemrograman Berorientasi Objek"`
		Sks  uint   `json:"sks" example:"3"`
	}

	KurikulumCreateRequest struct {
		Kode      string `json:"kode" binding:"required,max=12" example:"myhutao-2024"`
		Name      string `json:"name" binding:"required,max=255" example:"Merdeka Belajar"`
		ProdiName string `json:"prodi_name" binding:"required,max=255" example:"teknik-informatika"`
	}

	KurikulumUpdateRequest struct {
		Kode      string `json:"kode" binding:"omitempty,max=12" example:"myhutao-2024"`
		Name      string `json:"name" binding:"omitempty,max=255" example:"Merdeka Belajar"`
		ProdiName string `json:"prodi_name" binding:"omitempty,max=255" example:"teknik-informatika"`
	}

	KurikulumQuery struct {
		ID   string `form:"id" binding:"omitempty,max=36" example:"000-000-000-000"`
		Kode string `form:"kode" binding:"omitempty,max=12" example:"myhutao-2024"`
		// ProdiName string `form:"prodi_name" binding:"omitempty,max=255" example:"teknik-informatika"`
	}

	KurikulumKodeURI struct {
		Kode string `uri:"kode" binding:"required,max=12" example:"myhutao-2024"`
	}
)

func ToKurikulumResponse(k entities.Kurikulum) KurikulumResponse {
	var mkResponses []KurikulumMKResponse
	for _, mk := range k.KurikulumMK {
		mkResponses = append(mkResponses, KurikulumMKResponse{
			Semester: mk.Semester,
			Wajib:    mk.Wajib,
			MataKuliah: MataKuliahResponse{
				ID:   mk.MataKuliah.ID.String(),
				Kode: mk.MataKuliah.Kode,
				Name: mk.MataKuliah.Name,
				Sks:  mk.MataKuliah.Sks,
			},
		})
	}

	prodiResponse := ProdiResponse{
		ID:      k.Prodi.ID,
		Name:    k.Prodi.Name,
		Jenjang: k.Prodi.Jenjang,
		Jurusan: JurusanResponse{
			ID:   k.Prodi.Jurusan.ID,
			Name: k.Prodi.Jurusan.Name,
		},
	}

	return KurikulumResponse{
		ID:          k.ID.String(),
		Kode:        k.Kode,
		Name:        k.Name,
		Prodi:       prodiResponse,
		KurikulumMK: mkResponses,
	}
}
