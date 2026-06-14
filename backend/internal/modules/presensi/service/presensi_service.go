package service

import (
	"context"
	"fmt"
	"web-hosting/internal/modules/presensi/dto"
	"web-hosting/internal/modules/presensi/repository"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PresensiService interface {
		CreatePresensi(ctx context.Context, req any) (any, error)
		UpdatePresensi(ctx context.Context, tipePresensi string, req any) error
		GetPresensi(ctx context.Context, tipePresensi string, filter any) (any, error)
		GetAllPresensiPaginated(ctx context.Context, tipePresensi string, filter any, page int) (any, int64, error)
		CountPresensi(ctx context.Context, tipe string) (int64, error)
		GetStatusPresensiPegawaiMe(ctx context.Context, tx *gorm.DB, pegawaiID uuid.UUID) (string, error)
	}
	presensiService struct {
		db   *gorm.DB
		repo repository.PresensiRepository
	}
)

func NewPresensiService(db *gorm.DB, repo repository.PresensiRepository) PresensiService {
	return &presensiService{
		db:   db,
		repo: repo,
	}
}

func (s *presensiService) CreatePresensi(ctx context.Context, req any) (any, error) {
	switch v := req.(type) {
	case dto.PresensiMahasiswaCreateRequest:
		data, err := s.repo.CreatePresensiMahasiswa(ctx, s.db, v.ID, v.PengampuID)
		if err != nil {
			return nil, err
		}
		return dto.ToPresensiResponse("mahasiswa", data), nil
	default:
		err := s.repo.CreatePresensiPegawai(ctx, s.db)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func (s *presensiService) UpdatePresensi(ctx context.Context, tipePresensi string, req any) error {
	tipePresensi = helpers.NormalizeString(tipePresensi)
	if err := s.repo.UpdatePresensi(ctx, s.db, req); err != nil {
		return err
	}
	return nil
}

func (s *presensiService) GetPresensi(ctx context.Context, tipePresensi string, filter any) (any, error) {
	tipePresensi = helpers.NormalizeString(tipePresensi)

	switch tipePresensi {
	case "mahasiswa":
		presensiID, err := uuid.Parse(filter.(string))
		if err != nil {
			return nil, err
		}
		data, err := s.repo.GetPresensiMahasiswa(ctx, s.db, presensiID)
		if err != nil {
			return nil, err
		}
		return dto.ToPresensiResponse(tipePresensi, data), nil
	case "pegawai":
		data, err := s.repo.GetPresensiPegawai(ctx, s.db, filter)
		if err != nil {
			return nil, err
		} else if data.ID == uuid.Nil {
			return nil, fmt.Errorf("presensi pegawai not found")
		}
		return dto.ToPresensiResponse(tipePresensi, data), nil

	default:
		return nil, fmt.Errorf("invalid tipe presensi: %s", tipePresensi)
	}
}

func (s *presensiService) GetAllPresensiPaginated(ctx context.Context, tipePresensi string, filter any, page int) (any, int64, error) {
	tipePresensi = helpers.NormalizeString(tipePresensi)
	switch tipePresensi {
	case "pegawai":
		offset := (page - 1) * 10
		data, total, err := s.repo.GetAllPresensiPegawai(ctx, s.db, offset, 10)
		if err != nil {
			return nil, 0, err
		}
		responses := make([]dto.PresensiPegawaiResponse, 0, len(data))
		for _, d := range data {
			responses = append(responses, dto.ToPresensiResponse(tipePresensi, d).(dto.PresensiPegawaiResponse))
		}
		return responses, total, nil
	default:
		return nil, 0, fmt.Errorf("invalid tipe presensi: %s", tipePresensi)
	}

}

func (s *presensiService) CountPresensi(ctx context.Context, tipe string) (int64, error) {
	tipe = helpers.NormalizeString(tipe)
	switch tipe {
	case "mahasiswa", "pegawai":
		return s.repo.CountPresensi(ctx, s.db, &tipe)
	default:
		return 0, fmt.Errorf("invalid tipe presensi: %s", tipe)
	}
}

func (s *presensiService) GetStatusPresensiPegawaiMe(ctx context.Context, tx *gorm.DB, pegawaiID uuid.UUID) (string, error) {
	status, err := s.repo.GetStatusPresensiPegawaiMe(ctx, tx, pegawaiID)
	if err != nil {
		return "", err
	}
	return status, nil
}
