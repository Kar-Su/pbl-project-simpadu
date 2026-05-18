package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/modules/jurusan/dto"
	"web-hosting/internal/modules/jurusan/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

type JurusanService interface {
	CreateJurusan(ctx context.Context, req dto.JurusanRequest) error
	UpdateJurusan(ctx context.Context, req dto.JurusanUpdateRequest, jurusanName string) (dto.JurusanResponse, error)
	DeleteJurusan(ctx context.Context, jurusanName string) error
	GetJurusanByName(ctx context.Context, jurusanName string) (dto.JurusanResponse, error)
	GetAllJurusan(ctx context.Context) ([]dto.JurusanResponse, error)
	GetAllJurusanPaginated(ctx context.Context, page int) ([]dto.JurusanResponse, int64, error)
	GetJurusanById(ctx context.Context, jurusanId uint) (dto.JurusanResponse, error)
}

type jurusanService struct {
	jurusanRepo repository.JurusanRepository
	db          *gorm.DB
}

func NewJurusanService(jurusanRepo repository.JurusanRepository, db *gorm.DB) JurusanService {
	return &jurusanService{jurusanRepo: jurusanRepo, db: db}
}

func (s *jurusanService) CreateJurusan(ctx context.Context, req dto.JurusanRequest) error {
	normName := helpers.NormalizeString(req.JurusanName)

	err := s.jurusanRepo.Create(ctx, s.db, normName)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrJurusanAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *jurusanService) UpdateJurusan(ctx context.Context, req dto.JurusanUpdateRequest, jurusanName string) (dto.JurusanResponse, error) {
	normName := helpers.NormalizeString(req.NewName)

	jurusan, err := s.jurusanRepo.GetByName(ctx, s.db, jurusanName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.JurusanResponse{}, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.JurusanResponse{}, constants.ErrInternalErr
	}

	jurusan.Name = normName

	updatedJurusan, err := s.jurusanRepo.Update(ctx, s.db, jurusanName, jurusan)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.JurusanResponse{}, constants.ErrInternalErr
	}

	return dto.ToJurusanResponse(updatedJurusan), nil
}

func (s *jurusanService) DeleteJurusan(ctx context.Context, jurusanName string) error {
	_, err := s.jurusanRepo.GetByName(ctx, s.db, jurusanName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	err = s.jurusanRepo.Delete(ctx, s.db, jurusanName)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *jurusanService) GetJurusanByName(ctx context.Context, jurusanName string) (dto.JurusanResponse, error) {
	jurusan, err := s.jurusanRepo.GetByName(ctx, s.db, jurusanName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.JurusanResponse{}, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.JurusanResponse{}, constants.ErrInternalErr
	}

	return dto.ToJurusanResponse(jurusan), nil
}

func (s *jurusanService) GetAllJurusan(ctx context.Context) ([]dto.JurusanResponse, error) {
	jurusan, err := s.jurusanRepo.GetAll(ctx, s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.JurusanResponse, len(jurusan))
	for i, j := range jurusan {
		responses[i] = dto.ToJurusanResponse(j)
	}

	return responses, nil
}

func (s *jurusanService) GetJurusanById(ctx context.Context, jurusanId uint) (dto.JurusanResponse, error) {
	jurusan, err := s.jurusanRepo.GetById(ctx, s.db, jurusanId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.JurusanResponse{}, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.JurusanResponse{}, constants.ErrInternalErr
	}

	return dto.ToJurusanResponse(jurusan), nil
}

func (s *jurusanService) GetAllJurusanPaginated(ctx context.Context, page int) ([]dto.JurusanResponse, int64, error) {
	offset := (page - 1) * 10
	jurusans, total, err := s.jurusanRepo.GetAllPaginated(ctx, s.db, offset, 10)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}
	responses := make([]dto.JurusanResponse, len(jurusans))
	for i, j := range jurusans {
		responses[i] = dto.ToJurusanResponse(j)
	}
	return responses, total, nil
}
