package service

import (
	"context"
	"errors"
	"log"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/kelas/dto"
	"web-hosting/internal/modules/kelas/repository"
	userRepository "web-hosting/internal/modules/user/repository"
	"web-hosting/internal/package/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	KelasMahasiswaService interface {
		Create(ctx context.Context, req dto.KelasMahasiswaCreateRequest) error
		Delete(ctx context.Context, mahasiswaId uuid.UUID, kelasId uuid.UUID) error
		GetAllKelasMahasiswa(ctx context.Context, mahasiswaId uuid.UUID) ([]dto.KelasMahasiswaResponse, error)
		GetMahasiswaByKelasId(ctx context.Context, kelasId uuid.UUID) ([]dto.KelasMahasiswaResponse, error)
	}
	kelasMahasiswaService struct {
		db                 *gorm.DB
		userRepo           userRepository.UserRepository
		kelasRepo          repository.KelasRepository
		kelasMahasiswaRepo repository.KelasMahasiswaRepository
	}
)

func NewKelasMahasiswaService(db *gorm.DB, userRepo userRepository.UserRepository, kelasRepo repository.KelasRepository, kelasMahasiswaRepo repository.KelasMahasiswaRepository) KelasMahasiswaService {
	return &kelasMahasiswaService{
		db:                 db,
		userRepo:           userRepo,
		kelasRepo:          kelasRepo,
		kelasMahasiswaRepo: kelasMahasiswaRepo,
	}
}

func (s *kelasMahasiswaService) Create(ctx context.Context, req dto.KelasMahasiswaCreateRequest) error {
	isExist, err := s.kelasRepo.CheckKelasById(ctx, s.db, req.KelasID)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if !isExist {
		return dto.ErrKelasNotFound
	}

	isExist, err = s.userRepo.CheckByMahasiswaRoleAndDetailID(ctx, s.db, req.MahasiswaID)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if !isExist {
		return dto.ErrMahasiswaNotFound
	}

	// LOGIC JIKA MAHASISWA HANYA BOLEH SATU KELAS
	// isAlreadyAssigned, err := s.kelasMahasiswaRepo.CheckMahasiswaAlreadyAssigned(ctx, s.db, req.MahasiswaID)
	// if err != nil {
	// 	log.Printf("Internal Error: %v", err)
	// 	return constants.ErrInternalErr
	// }
	// if isAlreadyAssigned {
	// 	return dto.ErrMahasiswaAlreadyAssigned
	// }

	assignEntity := entities.KelasMahasiswa{
		KelasID:     req.KelasID,
		MahasiswaID: req.MahasiswaID,
	}

	if err := s.kelasMahasiswaRepo.Create(ctx, s.db, assignEntity); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.ErrMahasiswaAlreadyAssigned
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *kelasMahasiswaService) Delete(ctx context.Context, mahasiswaId uuid.UUID, kelasId uuid.UUID) error {
	if err := s.kelasMahasiswaRepo.Delete(ctx, s.db, mahasiswaId, kelasId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrMahasiswaNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	return nil
}

func (s *kelasMahasiswaService) GetAllKelasMahasiswa(ctx context.Context, mahasiswaId uuid.UUID) ([]dto.KelasMahasiswaResponse, error) {
	entities, err := s.kelasMahasiswaRepo.GetAllKelasMahasiswa(ctx, s.db, mahasiswaId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrMahasiswaNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.KelasMahasiswaResponse, len(entities))
	for i, entity := range entities {
		responses[i] = dto.ToKelasMahasiswaResponse(entity)
	}

	return responses, nil
}

func (s *kelasMahasiswaService) GetMahasiswaByKelasId(ctx context.Context, kelasId uuid.UUID) ([]dto.KelasMahasiswaResponse, error) {
	entities, err := s.kelasMahasiswaRepo.GetMahasiswaByKelasId(ctx, s.db, kelasId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrKelasNotFound
		}
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}

	responses := make([]dto.KelasMahasiswaResponse, len(entities))
	for i, entity := range entities {
		responses[i] = dto.ToKelasMahasiswaResponse(entity)
	}

	return responses, nil
}
