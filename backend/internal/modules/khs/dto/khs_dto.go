package dto

import (
	"fmt"
	"strings"
	"web-hosting/internal/database/entities"

	"github.com/google/uuid"
)

type CreateKhsRequest struct {
	MahasiswaID string  `json:"mahasiswa_id" binding:"required" example:"miaco-3525sdvs-3vfdv"`
	PengampuID  string  `json:"pengampu_id" binding:"required" example:"miaco-3525sdvs-3vfdv"`
	TotalNilai  float32 `json:"total_nilai" binding:"required,gte=0,lte=100" example:"4"`
}

type FilterQuery struct {
	Semester  uint   `form:"semester" binding:"omitempty" example:"1"`
	ProdiName string `form:"prodi_name" binding:"omitempty" example:"teknik informatika"`
}

type NilaiResponse struct {
	KodeMK string  `json:"kode_mk"`
	NamaMK string  `json:"nama_mk"`
	Sks    uint    `json:"sks"`
	Nilai  float32 `json:"nilai"`
	Grade  string  `json:"grade"`
}

type KHSResponse struct {
	MahasiswaName string          `json:"mahasiswa_name"`
	Semester      uint            `json:"semester"`
	ProdiName     string          `json:"prodi_name"`
	KelasName     string          `json:"kelas_name"`
	TahunAkademik uint            `json:"tahun_akademik"`
	KurikulumName string          `json:"kurikulum_name"`
	Nilai         []NilaiResponse `json:"nilai"`
}

func formatNamaKelasDinamis(name string, semester uint) string {
	if strings.Contains(name, "-") {
		return strings.Replace(name, "-", fmt.Sprintf("-%d", semester), 1)
	}
	if len(name) > 1 {
		posisiSisip := len(name) - 1
		return fmt.Sprintf("%s%d%s", name[:posisiSisip], semester, name[posisiSisip:])
	}
	return fmt.Sprintf("%s-%d", name, semester)
}

func ToKHSResponse(entity entities.Khs) KHSResponse {
	var nilaiResponses []NilaiResponse

	var prodiName, kelasName, kurikulumName string
	var tahunAkademik uint

	for i, nilai := range entity.NilaiMk {
		nilaiResponses = append(nilaiResponses, NilaiResponse{
			KodeMK: nilai.Pengampu.MataKuliah.Kode,
			NamaMK: nilai.Pengampu.MataKuliah.Name,
			Sks:    uint(nilai.Pengampu.MataKuliah.Sks),
			Nilai:  nilai.TotalNilai,
			Grade:  nilai.GradeNilai,
		})

		if i == 0 && nilai.Pengampu.Kelas.ID != uuid.Nil {
			prodiName = nilai.Pengampu.Kelas.Prodi.Name
			kurikulumName = nilai.Pengampu.Kelas.Kurikulum.Name

			kelasName = formatNamaKelasDinamis(nilai.Pengampu.Kelas.Name, entity.Semester)

			tahunAkademik = nilai.Pengampu.Kelas.TahunAkademikID
		}
	}

	return KHSResponse{
		MahasiswaName: entity.Mahasiswa.Name,
		Semester:      entity.Semester,
		ProdiName:     prodiName,
		KelasName:     kelasName,
		TahunAkademik: tahunAkademik,
		KurikulumName: kurikulumName,
		Nilai:         nilaiResponses,
	}
}

func ToKHSResponseList(entities []entities.Khs) []KHSResponse {
	var responses []KHSResponse
	for _, entity := range entities {
		responses = append(responses, ToKHSResponse(entity))
	}
	return responses
}
