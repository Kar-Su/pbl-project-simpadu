package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/jurusan/dto"
	"web-hosting/internal/modules/jurusan/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

type JurusanService interface {
	CreateJurusan(ctx context.Context, req dto.JurusanRequest) (entities.Jurusan, error)
	UpdateJurusan(ctx context.Context, req dto.JurusanUpdateRequest, jurusanName string) (entities.Jurusan, error)
	DeleteJurusan(ctx context.Context, jurusanName string) error
	GetJurusanByName(ctx context.Context, jurusanName string) (entities.Jurusan, error)
	GetAllJurusan(ctx context.Context) ([]entities.Jurusan, error)
	GetJurusanById(ctx context.Context, jurusanId uint) (entities.Jurusan, error)
}

type jurusanService struct {
	jurusanRepo repository.JurusanRepository
	db          *gorm.DB
}

func NewJurusanService(jurusanRepo repository.JurusanRepository, db *gorm.DB) JurusanService {
	return &jurusanService{jurusanRepo: jurusanRepo, db: db}
}

func (s *jurusanService) CreateJurusan(ctx context.Context, req dto.JurusanRequest) (entities.Jurusan, error) {
	normName := helpers.NormalizeString(req.JurusanName)

	jurusan, err := s.jurusanRepo.Create(ctx, s.db, normName)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.Jurusan{}, dto.ErrJurusanAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return entities.Jurusan{}, constants.ErrInternalErr
	}

	return jurusan, nil
}

func (s *jurusanService) UpdateJurusan(ctx context.Context, req dto.JurusanUpdateRequest, jurusanName string) (entities.Jurusan, error) {
	normName := helpers.NormalizeString(req.NewName)

	jurusan, err := s.jurusanRepo.GetByName(ctx, s.db, jurusanName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Jurusan{}, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return entities.Jurusan{}, constants.ErrInternalErr
	}

	jurusan.Name = normName

	updatedJurusan, err := s.jurusanRepo.Update(ctx, s.db, jurusanName, jurusan)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return entities.Jurusan{}, constants.ErrInternalErr
	}

	return updatedJurusan, nil
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

func (s *jurusanService) GetJurusanByName(ctx context.Context, jurusanName string) (entities.Jurusan, error) {
	jurusan, err := s.jurusanRepo.GetByName(ctx, s.db, jurusanName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Jurusan{}, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return entities.Jurusan{}, constants.ErrInternalErr
	}

	return jurusan, nil
}

func (s *jurusanService) GetAllJurusan(ctx context.Context) ([]entities.Jurusan, error) {
	jurusan, err := s.jurusanRepo.GetAll(ctx, s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	return jurusan, nil
}

func (s *jurusanService) GetJurusanById(ctx context.Context, jurusanId uint) (entities.Jurusan, error) {
	jurusan, err := s.jurusanRepo.GetById(ctx, s.db, jurusanId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Jurusan{}, dto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return entities.Jurusan{}, constants.ErrInternalErr
	}

	return jurusan, nil
}
