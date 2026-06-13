package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/database/entities"
	akademikRepo "web-hosting/internal/modules/akademik/repository"
	"web-hosting/internal/modules/kelas/dto"
	kelasRepo "web-hosting/internal/modules/kelas/repository"
	kurikulumRepo "web-hosting/internal/modules/kurikulum/repository"
	prodiRepo "web-hosting/internal/modules/prodi/repository"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	KelasService interface {
		CreateKelas(ctx context.Context, req dto.KelasCreateRequest) error
		UpdateKelas(ctx context.Context, id uuid.UUID, req dto.KelasUpdateRequest) (dto.KelasResponse, error)
		DeleteKelas(ctx context.Context, id uuid.UUID) error
		GetKelasByID(ctx context.Context, id uuid.UUID) (dto.KelasResponse, error)
		GetKelasByProdiName(ctx context.Context, prodiName string) ([]dto.KelasResponse, error)
		GetKelasByProdiNamePaginated(ctx context.Context, prodiName string, page int) ([]dto.KelasResponse, int64, error)
		GetAllKelas(ctx context.Context, page int) ([]dto.KelasResponse, int64, error)
	}

	kelasService struct {
		db            *gorm.DB
		kelasRepo     kelasRepo.KelasRepository
		akademikRepo  akademikRepo.TahunAkademikRepository
		prodiRepo     prodiRepo.ProdiRepository
		kurikulumRepo kurikulumRepo.KurikulumRepository
	}
)

func NewKelasService(db *gorm.DB, kelasRepo kelasRepo.KelasRepository, akademikRepo akademikRepo.TahunAkademikRepository, prodiRepo prodiRepo.ProdiRepository, kurikulumRepo kurikulumRepo.KurikulumRepository) KelasService {
	return &kelasService{
		db:            db,
		kelasRepo:     kelasRepo,
		akademikRepo:  akademikRepo,
		prodiRepo:     prodiRepo,
		kurikulumRepo: kurikulumRepo,
	}
}

func (s *kelasService) CreateKelas(ctx context.Context, req dto.KelasCreateRequest) error {
	req.ProdiName = helpers.NormalizeString(req.ProdiName)
	req.KurikulumKode = helpers.NormalizeString(req.KurikulumKode)

	Prodi, err := s.prodiRepo.GetByName(ctx, s.db, req.ProdiName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrProdiNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	isExist, err := s.akademikRepo.CheckByID(ctx, s.db, req.TahunAkademikID)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if !isExist {
		return dto.ErrTahunAkademikNotFound
	}

	isExist, err = s.kurikulumRepo.CheckKurikulumExistsByKode(ctx, s.db, req.KurikulumKode)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if !isExist {
		return dto.ErrKurikulumNotFound
	}

	entity := entities.Kelas{
		Name:            req.Name,
		Semester:        req.Semester,
		TahunAkademikID: req.TahunAkademikID,
		KurikulumKode:   req.KurikulumKode,
		ProdiID:         Prodi.ID,
	}

	if err := s.kelasRepo.Create(ctx, s.db, entity); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrKelasAlreadyExists
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *kelasService) UpdateKelas(ctx context.Context, id uuid.UUID, req dto.KelasUpdateRequest) (dto.KelasResponse, error) {
	kelasEntity, err := s.kelasRepo.GetNonPreload(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.KelasResponse{}, dto.ErrKelasNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.KelasResponse{}, constants.ErrInternalErr
	}

	if req.Name != "" {
		kelasEntity.Name = req.Name
	}
	if req.Semester != 0 {
		kelasEntity.Semester = req.Semester
	}
	if req.TahunAkademikID != 0 {
		isExist, err := s.akademikRepo.CheckByID(ctx, s.db, req.TahunAkademikID)
		if err != nil {
			log.Printf("Internal Error: %v", err)
			return dto.KelasResponse{}, constants.ErrInternalErr
		}
		if !isExist {
			return dto.KelasResponse{}, dto.ErrTahunAkademikNotFound
		}
		kelasEntity.TahunAkademikID = req.TahunAkademikID
	}
	if req.KurikulumKode != "" {
		req.KurikulumKode = helpers.NormalizeString(req.KurikulumKode)
		isExist, err := s.kurikulumRepo.CheckKurikulumExistsByKode(ctx, s.db, req.KurikulumKode)
		if err != nil {
			log.Printf("Internal Error: %v", err)
			return dto.KelasResponse{}, constants.ErrInternalErr
		}
		if !isExist {
			return dto.KelasResponse{}, dto.ErrKurikulumNotFound
		}
		kelasEntity.KurikulumKode = req.KurikulumKode
	}
	if req.ProdiName != "" {
		req.ProdiName = helpers.NormalizeString(req.ProdiName)
		prodiEntity, err := s.prodiRepo.GetByName(ctx, s.db, req.ProdiName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.KelasResponse{}, dto.ErrProdiNotFound
			}
			log.Printf("Internal Error: %v", err)
			return dto.KelasResponse{}, constants.ErrInternalErr
		}
		kelasEntity.ProdiID = prodiEntity.ID
	}

	if err := s.kelasRepo.Update(ctx, s.db, id, kelasEntity); err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.KelasResponse{}, constants.ErrInternalErr
	}
	return dto.ToKelasResponse(kelasEntity), nil
}

func (s *kelasService) DeleteKelas(ctx context.Context, id uuid.UUID) error {
	err := s.kelasRepo.Delete(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrKelasNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *kelasService) GetKelasByID(ctx context.Context, id uuid.UUID) (dto.KelasResponse, error) {
	kelasEntity, err := s.kelasRepo.GetByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.KelasResponse{}, dto.ErrKelasNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.KelasResponse{}, constants.ErrInternalErr
	}
	return dto.ToKelasResponse(kelasEntity), nil
}

func (s *kelasService) GetKelasByProdiName(ctx context.Context, prodiName string) ([]dto.KelasResponse, error) {
	prodi, err := s.prodiRepo.GetByName(ctx, s.db, prodiName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrProdiNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}
	kelasEntity, err := s.kelasRepo.GetByProdiID(ctx, s.db, prodi.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrKelasNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.KelasResponse, len(kelasEntity))
	for i, entity := range kelasEntity {
		responses[i] = dto.ToKelasResponse(entity)
	}
	return responses, nil
}

func (s *kelasService) GetAllKelas(ctx context.Context, page int) ([]dto.KelasResponse, int64, error) {
	offset := (page - 1) * 10
	kelasEntities, total, err := s.kelasRepo.GetAll(ctx, s.db, offset, 10)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}

	responses := make([]dto.KelasResponse, len(kelasEntities))
	for i, entity := range kelasEntities {
		responses[i] = dto.ToKelasResponse(entity)
	}
	return responses, total, nil
}

func (s *kelasService) GetKelasByProdiNamePaginated(ctx context.Context, prodiName string, page int) ([]dto.KelasResponse, int64, error) {
	prodiName = helpers.NormalizeString(prodiName)

	prodi, err := s.prodiRepo.GetByName(ctx, s.db, prodiName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, dto.ErrProdiNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}

	offset := (page - 1) * 10
	kelasEntities, total, err := s.kelasRepo.GetByProdiIDPaginated(ctx, s.db, prodi.ID, offset, 10)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}

	responses := make([]dto.KelasResponse, len(kelasEntities))
	for i, entity := range kelasEntities {
		responses[i] = dto.ToKelasResponse(entity)
	}
	return responses, total, nil
}
