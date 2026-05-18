package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/kurikulum/dto"
	kurikulumRepo "web-hosting/internal/modules/kurikulum/repository"
	mkRepo "web-hosting/internal/modules/mk/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

type (
	KurikulumMKService interface {
		CreateKurikulumMK(ctx context.Context, req dto.PivotCreateRequest) error
		UpdateKurikulumMK(ctx context.Context, kurikulumKode string, mkKode string, req dto.PivotUpdateRequest) (dto.PivotResponse, error)
		DeleteKurikulumMK(ctx context.Context, kurikulumKode string, mkKode string) error
	}

	kurikulumMKService struct {
		kurikulumRepo kurikulumRepo.KurikulumRepository
		pivotRepo     kurikulumRepo.KurikulumMKRepo
		mkRepo        mkRepo.MkRepository
		db            *gorm.DB
	}
)

func NewKurikulumMKService(db *gorm.DB, kurikulumRepo kurikulumRepo.KurikulumRepository, pivotRepo kurikulumRepo.KurikulumMKRepo, mkRepo mkRepo.MkRepository) KurikulumMKService {
	return &kurikulumMKService{
		kurikulumRepo: kurikulumRepo,
		pivotRepo:     pivotRepo,
		mkRepo:        mkRepo,
		db:            db,
	}
}

func (s *kurikulumMKService) CreateKurikulumMK(ctx context.Context, req dto.PivotCreateRequest) error {
	req.MkKode = helpers.NormalizeString(req.MkKode)
	isExist, err := s.mkRepo.CheckMkExists(ctx, s.db, req.MkKode)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if !isExist {
		return dto.ErrMkNotFound
	}

	req.KurikulumKode = helpers.NormalizeString(req.KurikulumKode)
	isExist, err = s.kurikulumRepo.CheckKurikulumExistsByKode(ctx, s.db, req.KurikulumKode)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if !isExist {
		return dto.ErrKurikulumNotFound
	}

	entity := entities.KurikulumMK{
		KurikulumKode: req.KurikulumKode,
		MkKode:        req.MkKode,
		Semester:      req.Semester,
		Wajib:         req.Wajib,
	}

	if err := s.pivotRepo.Create(ctx, s.db, entity); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrPivotAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *kurikulumMKService) UpdateKurikulumMK(ctx context.Context, kurikulumKode string, mkKode string, req dto.PivotUpdateRequest) (dto.PivotResponse, error) {
	kurikulumKode = helpers.NormalizeString(kurikulumKode)
	mkKode = helpers.NormalizeString(mkKode)

	entity, err := s.pivotRepo.Get(ctx, s.db, kurikulumKode, mkKode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PivotResponse{}, dto.ErrPivotNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.PivotResponse{}, constants.ErrInternalErr
	}

	if req.Semester != nil {
		entity.Semester = *req.Semester
	}
	if req.Wajib != nil {
		entity.Wajib = *req.Wajib
	}

	if err := s.pivotRepo.Update(ctx, s.db, kurikulumKode, mkKode, entity); err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.PivotResponse{}, constants.ErrInternalErr
	}

	return dto.ToPivotResponse(entity), nil
}

func (s *kurikulumMKService) DeleteKurikulumMK(ctx context.Context, kurikulumKode string, mkKode string) error {
	kurikulumKode = helpers.NormalizeString(kurikulumKode)
	mkKode = helpers.NormalizeString(mkKode)

	if err := s.pivotRepo.Delete(ctx, s.db, kurikulumKode, mkKode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrPivotNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}
