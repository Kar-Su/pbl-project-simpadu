package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"web-hosting/internal/database/entities"
	jurusanDto "web-hosting/internal/modules/jurusan/dto"
	jurusanService "web-hosting/internal/modules/jurusan/service"
	"web-hosting/internal/modules/prodi/dto"
	"web-hosting/internal/modules/prodi/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"gorm.io/gorm"
)

type (
	ProdiService interface {
		CreateProdi(ctx context.Context, req dto.ProdiCreateRequest) error
		UpdateProdi(ctx context.Context, prodiName string, req dto.ProdiUpdateRequest) (dto.ProdiResponse, error)
		DeleteProdi(ctx context.Context, prodiName string) error
		GetProdiByName(ctx context.Context, prodiName string) (dto.ProdiResponse, error)
		GetProdiById(ctx context.Context, prodiId uint) (dto.ProdiResponse, error)
		GetAllProdi(ctx context.Context) ([]dto.ProdiResponse, error)
		GetProdiByJurusanName(ctx context.Context, jurusanName string) ([]dto.ProdiResponse, error)
	}

	prodiService struct {
		prodiRepo      repository.ProdiRepository
		jurusanService jurusanService.JurusanService
		db             *gorm.DB
	}
)

func NewProdiService(prodiRepo repository.ProdiRepository, jurusanService jurusanService.JurusanService, db *gorm.DB) ProdiService {
	return &prodiService{
		prodiRepo:      prodiRepo,
		jurusanService: jurusanService,
		db:             db,
	}
}

func (s *prodiService) CreateProdi(ctx context.Context, req dto.ProdiCreateRequest) error {
	req.JurusanName = helpers.NormalizeString(req.JurusanName)
	jurusan, err := s.jurusanService.GetJurusanByName(ctx, req.JurusanName)
	if err != nil {
		return err
	}
	req.Name = helpers.NormalizeString(req.Name)

	prodiEntity := entities.Prodi{
		Name:      req.Name,
		Jenjang:   strings.ToUpper(req.Jenjang),
		JurusanID: jurusan.ID,
	}

	if err := s.prodiRepo.Create(ctx, s.db, prodiEntity); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrProdiAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *prodiService) UpdateProdi(ctx context.Context, prodiName string, req dto.ProdiUpdateRequest) (dto.ProdiResponse, error) {
	prodiName = helpers.NormalizeString(prodiName)

	var prodiEntity entities.Prodi
	if req.Name != "" {
		prodiEntity.Name = helpers.NormalizeString(req.Name)
	}
	if req.Jenjang != "" {
		prodiEntity.Jenjang = strings.ToUpper(req.Jenjang)
	}
	if req.JurusanID != 0 {
		prodiEntity.JurusanID = req.JurusanID
	}

	updatedProdi, err := s.prodiRepo.Update(ctx, s.db, prodiName, prodiEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ProdiResponse{}, dto.ErrProdiNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.ProdiResponse{}, constants.ErrInternalErr
	}

	return dto.ToProdiResponse(updatedProdi), nil
}

func (s *prodiService) DeleteProdi(ctx context.Context, prodiName string) error {
	prodiName = helpers.NormalizeString(prodiName)

	if err := s.prodiRepo.Delete(ctx, s.db, prodiName); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrProdiNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *prodiService) GetProdiByName(ctx context.Context, prodiName string) (dto.ProdiResponse, error) {
	prodiName = helpers.NormalizeString(prodiName)

	prodi, err := s.prodiRepo.GetByName(ctx, s.db, prodiName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ProdiResponse{}, dto.ErrProdiNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.ProdiResponse{}, constants.ErrInternalErr
	}

	return dto.ToProdiResponse(prodi), nil
}

func (s *prodiService) GetProdiById(ctx context.Context, prodiId uint) (dto.ProdiResponse, error) {
	prodi, err := s.prodiRepo.GetByID(ctx, s.db, prodiId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ProdiResponse{}, dto.ErrProdiNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.ProdiResponse{}, constants.ErrInternalErr
	}

	return dto.ToProdiResponse(prodi), nil
}

func (s *prodiService) GetAllProdi(ctx context.Context) ([]dto.ProdiResponse, error) {
	prodis, err := s.prodiRepo.GetAll(ctx, s.db)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.ProdiResponse, len(prodis))
	for i, prodi := range prodis {
		responses[i] = dto.ToProdiResponse(prodi)
	}

	return responses, nil
}

func (s *prodiService) GetProdiByJurusanName(ctx context.Context, jurusanName string) ([]dto.ProdiResponse, error) {
	jurusanName = helpers.NormalizeString(jurusanName)
	jurusan, err := s.jurusanService.GetJurusanByName(ctx, jurusanName)
	if err != nil {
		return nil, err
	}

	prodis, err := s.prodiRepo.GetByJurusanID(ctx, s.db, jurusan.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, jurusanDto.ErrJurusanNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.ProdiResponse, len(prodis))
	for i, prodi := range prodis {
		responses[i] = dto.ToProdiResponse(prodi)
	}

	return responses, nil
}
