package service

import (
	"context"
	"errors"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/khs/dto"
	"web-hosting/internal/modules/khs/repository"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KhsService interface {
	Create(ctx context.Context, nilai *dto.CreateKhsRequest) error
	GetKHS(ctx context.Context, filter *dto.FilterQuery, page int) ([]dto.KHSResponse, error)
}

type khsService struct {
	khsRepo repository.KhsRepository
	db      *gorm.DB
}

func NewKhsService(khsRepo repository.KhsRepository, db *gorm.DB) KhsService {
	return &khsService{
		khsRepo: khsRepo,
		db:      db,
	}
}

func (s *khsService) Create(ctx context.Context, nilai *dto.CreateKhsRequest) error {
	pengampuID, err := uuid.Parse(nilai.PengampuID)
	if err != nil {
		return errors.New("Failed Parse, Please provide UUIDV7")
	}
	entity := &entities.NilaiMk{
		PengampuID: pengampuID,
		TotalNilai: nilai.TotalNilai,
	}

	mahasiswaID, err := uuid.Parse(nilai.MahasiswaID)
	if err != nil {
		return errors.New("Failed Parse, Please provide UUIDV7")
	}

	if err := s.khsRepo.Create(ctx, s.db, mahasiswaID, entity); err != nil {
		return err
	}
	return nil
}

func (s *khsService) GetKHS(ctx context.Context, filter *dto.FilterQuery, page int) ([]dto.KHSResponse, error) {
	filter.ProdiName = helpers.NormalizeString(filter.ProdiName)
	offset := (page - 1) * 10
	khs, err := s.khsRepo.GetKHS(ctx, s.db, filter, offset, 10)
	if err != nil {
		return nil, err
	}

	return dto.ToKHSResponseList(khs), nil
}
