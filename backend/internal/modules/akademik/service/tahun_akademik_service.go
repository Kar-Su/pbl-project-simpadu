package service

import (
	"context"
	"errors"
	"log"
	"time"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/akademik/dto"
	"web-hosting/internal/modules/akademik/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

type (
	TahunAkademikService interface {
		CreateTahunAkademik(ctx context.Context, req dto.AkademikCreateRequest) error
		UpdateTahunAkademik(ctx context.Context, id uint, req dto.AkademikUpdateRequest) (dto.AkademikResponse, error)
		DeleteTahunAkademik(ctx context.Context, id uint) error
		GetTahunAkademikByID(ctx context.Context, id uint) (dto.AkademikResponse, error)
		GetAllTahunAkademik(ctx context.Context) ([]dto.AkademikResponse, error)
		GetTahunAkademikByStatus(ctx context.Context, status string) ([]dto.AkademikResponse, error)
	}

	tahunAkademikService struct {
		TahunAkademikRepo repository.TahunAkademikRepository
		db                *gorm.DB
	}
)

func NewTahunAkademikService(repo repository.TahunAkademikRepository, db *gorm.DB) TahunAkademikService {
	return &tahunAkademikService{
		TahunAkademikRepo: repo,
		db:                db,
	}
}

func (s *tahunAkademikService) CreateTahunAkademik(ctx context.Context, req dto.AkademikCreateRequest) error {
	req.TipeSemester = helpers.NormalizeString(req.TipeSemester)

	enyity := entities.TahunAkademik{
		ID:           req.ID,
		TipeSemester: req.TipeSemester,
		TahunAwal:    req.TahunAwal,
		TahunAkhir:   req.TahunAkhir,
	}

	tAwal := time.Time(req.TahunAwal)
	tAkhir := time.Time(req.TahunAkhir)

	if tAwal.After(tAkhir) {
		return dto.ErrInvalidTahunAkademik
	}

	err := s.TahunAkademikRepo.Create(ctx, s.db, enyity)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrTahunAkademikAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *tahunAkademikService) UpdateTahunAkademik(ctx context.Context, id uint, req dto.AkademikUpdateRequest) (dto.AkademikResponse, error) {
	entity, err := s.TahunAkademikRepo.GetByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AkademikResponse{}, dto.ErrTahunAkademikNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.AkademikResponse{}, constants.ErrInternalErr
	}

	if req.ID != 0 {
		entity.ID = req.ID
	}
	if req.TipeSemester != "" {
		entity.TipeSemester = helpers.NormalizeString(req.TipeSemester)
	}
	if !req.TahunAwal.IsZero() {
		entity.TahunAwal = req.TahunAwal
	}
	if !req.TahunAkhir.IsZero() {
		entity.TahunAkhir = req.TahunAkhir
	}
	if req.Status != "" {
		entity.Status = helpers.NormalizeString(req.Status)
	}

	tAwal := time.Time(entity.TahunAwal)
	tAkhir := time.Time(entity.TahunAkhir)
	if tAwal.After(tAkhir) {
		return dto.AkademikResponse{}, dto.ErrInvalidTahunAkademik
	}

	if err := s.TahunAkademikRepo.Update(ctx, s.db, id, entity); err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.AkademikResponse{}, constants.ErrInternalErr
	}

	return dto.ToAkademikResponse(entity), nil
}

func (s *tahunAkademikService) DeleteTahunAkademik(ctx context.Context, id uint) error {
	err := s.TahunAkademikRepo.Delete(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrTahunAkademikNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *tahunAkademikService) GetTahunAkademikByID(ctx context.Context, id uint) (dto.AkademikResponse, error) {
	entity, err := s.TahunAkademikRepo.GetByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AkademikResponse{}, dto.ErrTahunAkademikNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.AkademikResponse{}, constants.ErrInternalErr
	}

	return dto.ToAkademikResponse(entity), err
}

func (s *tahunAkademikService) GetAllTahunAkademik(ctx context.Context) ([]dto.AkademikResponse, error) {
	entities, err := s.TahunAkademikRepo.GetAll(ctx, s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrTahunAkademikNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.AkademikResponse, len(entities))
	for i, entity := range entities {
		responses[i] = dto.ToAkademikResponse(entity)
	}

	return responses, nil
}

func (s *tahunAkademikService) GetTahunAkademikByStatus(ctx context.Context, status string) ([]dto.AkademikResponse, error) {
	entities, err := s.TahunAkademikRepo.GetByStatus(ctx, s.db, status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrTahunAkademikNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.AkademikResponse, len(entities))
	for i, entity := range entities {
		responses[i] = dto.ToAkademikResponse(entity)
	}

	return responses, nil
}
