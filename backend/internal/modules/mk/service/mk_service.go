package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/mk/dto"
	"web-hosting/internal/modules/mk/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MkService interface {
		CreateMk(ctx context.Context, req dto.MkCreateRequest) error
		UpdateMkByKode(ctx context.Context, kode string, req dto.MkUpdateRequest) (dto.MkResponse, error)
		UpdateMkById(ctx context.Context, id uuid.UUID, req dto.MkUpdateRequest) (dto.MkResponse, error)
		DeleteMkById(ctx context.Context, id uuid.UUID) error
		DeleteMkByKode(ctx context.Context, kode string) error
		GetMkById(ctx context.Context, id uuid.UUID) (dto.MkResponse, error)
		GetMkByKode(ctx context.Context, kode string) (dto.MkResponse, error)
		GetAllMk(ctx context.Context) ([]dto.MkResponse, error)
		GetAllMkPaginated(ctx context.Context, page int) ([]dto.MkResponse, int64, error)
	}

	mkService struct {
		MkRepository repository.MkRepository
		db           *gorm.DB
	}
)

func NewMkService(mkRepository repository.MkRepository, db *gorm.DB) MkService {
	return &mkService{
		MkRepository: mkRepository,
		db:           db,
	}
}

func (s *mkService) CreateMk(ctx context.Context, req dto.MkCreateRequest) error {
	req.Kode = helpers.NormalizeString(req.Kode)

	mkEntity := entities.MataKuliah{
		Kode: req.Kode,
		Name: req.Name,
		Sks:  req.Sks,
	}

	err := s.MkRepository.Create(ctx, s.db, mkEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrMkAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *mkService) UpdateMkByKode(ctx context.Context, kode string, req dto.MkUpdateRequest) (dto.MkResponse, error) {
	req.Kode = helpers.NormalizeString(req.Kode)

	mkEntity := entities.MataKuliah{}

	if req.Kode != "" {
		mkEntity.Kode = req.Kode
	}
	if req.Name != "" {
		mkEntity.Name = req.Name
	}
	if req.Sks != 0 {
		mkEntity.Sks = req.Sks
	}

	updatedMk, err := s.MkRepository.UpdateByKode(ctx, s.db, kode, mkEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.MkResponse{}, dto.ErrMkNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.MkResponse{}, constants.ErrInternalErr
	}

	return dto.ToMkResponse(updatedMk), nil
}

func (s *mkService) UpdateMkById(ctx context.Context, id uuid.UUID, req dto.MkUpdateRequest) (dto.MkResponse, error) {
	req.Kode = helpers.NormalizeString(req.Kode)

	mkEntity := entities.MataKuliah{}

	if req.Kode != "" {
		mkEntity.Kode = req.Kode
	}
	if req.Name != "" {
		mkEntity.Name = req.Name
	}
	if req.Sks != 0 {
		mkEntity.Sks = req.Sks
	}

	updatedMk, err := s.MkRepository.UpdateById(ctx, s.db, id, mkEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.MkResponse{}, dto.ErrMkNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.MkResponse{}, constants.ErrInternalErr
	}

	return dto.ToMkResponse(updatedMk), nil
}

func (s *mkService) DeleteMkById(ctx context.Context, id uuid.UUID) error {
	if err := s.MkRepository.DeleteById(ctx, s.db, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrMkNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *mkService) DeleteMkByKode(ctx context.Context, kode string) error {
	if err := s.MkRepository.DeleteByKode(ctx, s.db, kode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrMkNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *mkService) GetMkById(ctx context.Context, id uuid.UUID) (dto.MkResponse, error) {
	result, err := s.MkRepository.GetById(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.MkResponse{}, dto.ErrMkNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.MkResponse{}, constants.ErrInternalErr
	}

	return dto.ToMkResponse(result), nil
}

func (s *mkService) GetMkByKode(ctx context.Context, kode string) (dto.MkResponse, error) {
	result, err := s.MkRepository.GetByKode(ctx, s.db, kode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.MkResponse{}, dto.ErrMkNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.MkResponse{}, constants.ErrInternalErr
	}

	return dto.ToMkResponse(result), nil
}

func (s *mkService) GetAllMk(ctx context.Context) ([]dto.MkResponse, error) {
	result, err := s.MkRepository.GetAll(ctx, s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrMkNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.MkResponse, len(result))

	for i, r := range result {
		responses[i] = dto.ToMkResponse(r)
	}

	return responses, nil
}

func (s *mkService) GetAllMkPaginated(ctx context.Context, page int) ([]dto.MkResponse, int64, error) {
	offset := (page - 1) * 10
	mks, total, err := s.MkRepository.GetAllPaginated(ctx, s.db, offset, 10)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}
	responses := make([]dto.MkResponse, len(mks))
	for i, r := range mks {
		responses[i] = dto.ToMkResponse(r)
	}
	return responses, total, nil
}
