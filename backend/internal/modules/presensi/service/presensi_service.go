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
		presensiID := uuid.MustParse(filter.(string))
		data, err := s.repo.GetPresensiMahasiswa(ctx, s.db, presensiID)
		if err != nil {
			return nil, err
		}
		return dto.ToPresensiResponse(tipePresensi, data), nil
	case "pegawai":
		data, err := s.repo.GetPresensiPegawai(ctx, s.db, filter)
		if err != nil {
			return nil, err
		}
		return dto.ToPresensiResponse(tipePresensi, data), nil

	default:
		return nil, fmt.Errorf("invalid tipe presensi: %s", tipePresensi)
	}
}
