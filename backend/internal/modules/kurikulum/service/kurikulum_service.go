package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/kurikulum/dto"
	kurikulumRepo "web-hosting/internal/modules/kurikulum/repository"
	prodiService "web-hosting/internal/modules/prodi/service"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	KurikulumService interface {
		CreateKurikulum(ctx context.Context, req dto.KurikulumCreateRequest) error
		UpdateKurikulumById(ctx context.Context, req dto.KurikulumUpdateRequest, id uuid.UUID) (dto.KurikulumResponse, error)
		UpdateKurikulumByKode(ctx context.Context, req dto.KurikulumUpdateRequest, kode string) (dto.KurikulumResponse, error)
		DeleteKurikulumByKode(ctx context.Context, kode string) error
		DeleteKurikulumById(ctx context.Context, id uuid.UUID) error
		GetKurikulumByKode(ctx context.Context, kode string) (dto.KurikulumResponse, error)
		GetKurikulumById(ctx context.Context, id uuid.UUID) (dto.KurikulumResponse, error)
		GetAllKurikulum(ctx context.Context) ([]dto.KurikulumResponse, error)
		GetAllKurikulumPaginated(ctx context.Context, page int) ([]dto.KurikulumResponse, int64, error)
	}

	kurikulumService struct {
		kurikulumRepository kurikulumRepo.KurikulumRepository
		prodiService        prodiService.ProdiService
		db                  *gorm.DB
	}
)

func NewKurikulumService(kurikulumRepository kurikulumRepo.KurikulumRepository, prodiService prodiService.ProdiService, db *gorm.DB) KurikulumService {
	return &kurikulumService{
		kurikulumRepository: kurikulumRepository,
		prodiService:        prodiService,
		db:                  db,
	}
}

func (s *kurikulumService) CreateKurikulum(ctx context.Context, req dto.KurikulumCreateRequest) error {
	req.Kode = helpers.NormalizeString(req.Kode)

	prodi, err := s.prodiService.GetProdiByName(ctx, req.ProdiName)
	if err != nil {
		return err
	}

	kurikulum := entities.Kurikulum{
		Kode:    req.Kode,
		Name:    req.Name,
		ProdiID: prodi.ID,
	}

	err = s.kurikulumRepository.Create(ctx, s.db, kurikulum)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrKurikulumAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *kurikulumService) UpdateKurikulumById(ctx context.Context, req dto.KurikulumUpdateRequest, id uuid.UUID) (dto.KurikulumResponse, error) {
	entity, err := s.kurikulumRepository.GetByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.KurikulumResponse{}, dto.ErrKurikulumNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.KurikulumResponse{}, constants.ErrInternalErr
	}

	if req.Kode != "" {
		entity.Kode = helpers.NormalizeString(req.Kode)
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.ProdiName != "" {
		prodi, err := s.prodiService.GetProdiByName(ctx, req.ProdiName)
		if err != nil {
			return dto.KurikulumResponse{}, err
		}
		entity.ProdiID = prodi.ID
	}

	if err := s.kurikulumRepository.UpdateByID(ctx, s.db, entity, id); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.KurikulumResponse{}, dto.ErrUpdateViolatesUnique
		}
		log.Printf("Internal Error: %v", err)
		return dto.KurikulumResponse{}, constants.ErrInternalErr
	}

	return dto.ToKurikulumResponse(entity), nil
}

func (s *kurikulumService) UpdateKurikulumByKode(ctx context.Context, req dto.KurikulumUpdateRequest, kode string) (dto.KurikulumResponse, error) {
	entity, err := s.kurikulumRepository.GetByKode(ctx, s.db, kode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.KurikulumResponse{}, dto.ErrKurikulumNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.KurikulumResponse{}, constants.ErrInternalErr
	}

	if req.Kode != "" {
		entity.Kode = helpers.NormalizeString(req.Kode)
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.ProdiName != "" {
		prodi, err := s.prodiService.GetProdiByName(ctx, req.ProdiName)
		if err != nil {
			return dto.KurikulumResponse{}, err
		}
		entity.ProdiID = prodi.ID
	}

	if err := s.kurikulumRepository.UpdateByKode(ctx, s.db, entity, kode); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.KurikulumResponse{}, dto.ErrUpdateViolatesUnique
		}
		log.Printf("Internal Error: %v", err)
		return dto.KurikulumResponse{}, constants.ErrInternalErr
	}

	return dto.ToKurikulumResponse(entity), nil
}

func (s *kurikulumService) DeleteKurikulumByKode(ctx context.Context, kode string) error {
	err := s.kurikulumRepository.DeleteByKode(ctx, s.db, kode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrKurikulumNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *kurikulumService) DeleteKurikulumById(ctx context.Context, id uuid.UUID) error {
	err := s.kurikulumRepository.DeleteByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrKurikulumNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *kurikulumService) GetKurikulumByKode(ctx context.Context, kode string) (dto.KurikulumResponse, error) {
	kode = helpers.NormalizeString(kode)

	entity, err := s.kurikulumRepository.GetByKode(ctx, s.db, kode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.KurikulumResponse{}, dto.ErrKurikulumNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.KurikulumResponse{}, constants.ErrInternalErr
	}
	return dto.ToKurikulumResponse(entity), nil
}

func (s *kurikulumService) GetKurikulumById(ctx context.Context, id uuid.UUID) (dto.KurikulumResponse, error) {
	entity, err := s.kurikulumRepository.GetByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.KurikulumResponse{}, dto.ErrKurikulumNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.KurikulumResponse{}, constants.ErrInternalErr
	}
	return dto.ToKurikulumResponse(entity), nil
}

func (s *kurikulumService) GetAllKurikulum(ctx context.Context) ([]dto.KurikulumResponse, error) {
	entities, err := s.kurikulumRepository.GetAll(ctx, s.db)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.KurikulumResponse, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, dto.ToKurikulumResponse(entity))
	}

	return responses, nil
}

func (s *kurikulumService) GetAllKurikulumPaginated(ctx context.Context, page int) ([]dto.KurikulumResponse, int64, error) {
	offset := (page - 1) * 10
	kurikulums, total, err := s.kurikulumRepository.GetAllPaginated(ctx, s.db, offset, 10)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}
	responses := make([]dto.KurikulumResponse, 0, len(kurikulums))
	for _, entity := range kurikulums {
		responses = append(responses, dto.ToKurikulumResponse(entity))
	}
	return responses, total, nil
}
