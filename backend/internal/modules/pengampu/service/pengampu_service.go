package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/pengampu/dto"
	"web-hosting/internal/modules/pengampu/repository"
	"web-hosting/internal/package/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PengampuService interface {
		CreatePengampu(ctx context.Context, req dto.CreatePengampuRequest) error
		UpdatePengampuByID(ctx context.Context, id uuid.UUID, req dto.UpdatePengampuRequest) (dto.PengampuResponse, error)
		DeletePengampuByID(ctx context.Context, id uuid.UUID) error
		GetPengampuByID(ctx context.Context, id uuid.UUID) (dto.PengampuResponse, error)
		GetPengampuByKelasID(ctx context.Context, kelasID uuid.UUID) ([]dto.PengampuResponse, error)
	}

	pengampuService struct {
		PengampuRepo repository.PengampuRepository
		db           *gorm.DB
	}
)

func NewPengampuService(db *gorm.DB, pengampuRepo repository.PengampuRepository) PengampuService {
	return &pengampuService{
		db:           db,
		PengampuRepo: pengampuRepo,
	}
}

func (s *pengampuService) CreatePengampu(ctx context.Context, req dto.CreatePengampuRequest) error {
	entity := entities.Pengampu{
		KelasID: req.KelasID,
		MKKode:  req.MKKode,
		DosenID: req.DosenID,
	}

	if err := s.PengampuRepo.Create(ctx, s.db, entity); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrDuplicatedKey
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *pengampuService) UpdatePengampuByID(ctx context.Context, id uuid.UUID, req dto.UpdatePengampuRequest) (dto.PengampuResponse, error) {
	entity, err := s.PengampuRepo.GetByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PengampuResponse{}, dto.ErrPengampuNotFound
		}
	}

	if req.DosenID != nil {
		entity.DosenID = *req.DosenID
	}
	if req.KelasID != nil {
		entity.KelasID = *req.KelasID
	}
	if req.MKKode != "" {
		entity.MKKode = req.MKKode
	}

	updatedEntity, err := s.PengampuRepo.UpdateByID(ctx, s.db, id, entity)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.PengampuResponse{}, dto.ErrDuplicatedKey
		}
		log.Printf("Internal Error: %v", err)
		return dto.PengampuResponse{}, constants.ErrInternalErr
	}

	return dto.ToPengampuResponse(updatedEntity), nil
}

func (s *pengampuService) DeletePengampuByID(ctx context.Context, id uuid.UUID) error {
	if err := s.PengampuRepo.DeleteByID(ctx, s.db, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrPengampuNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *pengampuService) GetPengampuByID(ctx context.Context, id uuid.UUID) (dto.PengampuResponse, error) {
	entity, err := s.PengampuRepo.GetByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PengampuResponse{}, dto.ErrPengampuNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.PengampuResponse{}, constants.ErrInternalErr
	}

	return dto.ToPengampuResponse(entity), nil
}

func (s *pengampuService) GetPengampuByKelasID(ctx context.Context, kelasID uuid.UUID) ([]dto.PengampuResponse, error) {
	entities, err := s.PengampuRepo.GetByKelasID(ctx, s.db, kelasID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrPengampuNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.PengampuResponse, len(entities))
	for i, entity := range entities {
		responses[i] = dto.ToPengampuResponse(entity)
	}
	return responses, nil
}
