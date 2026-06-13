package repository

import (
	"context"
	"fmt"
	"time"
	"web-hosting/internal/database/entities"
	"web-hosting/internal/modules/presensi/dto"
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"
	"web-hosting/internal/package/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	PresensiRepository interface {
		CreatePresensiMahasiswa(ctx context.Context, tx *gorm.DB, id uuid.UUID, pengampuID uuid.UUID) (entities.Presensi, error)
		CreatePresensiPegawai(ctx context.Context, tx *gorm.DB) error
		UpdatePresensi(ctx context.Context, tx *gorm.DB, req any) error
		GetPresensiMahasiswa(ctx context.Context, tx *gorm.DB, presensiID uuid.UUID) (entities.Presensi, error)
		GetPresensiPegawai(ctx context.Context, tx *gorm.DB, filter any) (entities.Presensi, error)
		GetAllPresensiPegawai(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entities.Presensi, int64, error)
		CountPresensi(ctx context.Context, tx *gorm.DB, tipe *string) (int64, error)
	}

	presensiRepository struct {
		db *gorm.DB
	}
)

func NewPresensiRepository(db *gorm.DB) PresensiRepository {
	return &presensiRepository{db: db}
}

func (r *presensiRepository) CreatePresensiMahasiswa(ctx context.Context, tx *gorm.DB, id uuid.UUID, pengampuID uuid.UUID) (entities.Presensi, error) {
	if tx == nil {
		tx = r.db
	}

	presensi := entities.Presensi{
		ID:         id,
		Tipe:       "mahasiswa",
		PengampuID: &pengampuID,
	}
	if err := tx.WithContext(ctx).Create(&presensi).Error; err != nil {
		return entities.Presensi{}, err
	}

	var pengampu entities.Pengampu
	if err := tx.WithContext(ctx).
		Select("kelas_id").
		Where("id = ?", pengampuID).
		First(&pengampu).Error; err != nil {
		return entities.Presensi{}, err
	}
	kelasID := pengampu.KelasID

	var kelasMahasiswaList []entities.KelasMahasiswa
	if err := tx.WithContext(ctx).
		Preload("Mahasiswa", helpers.SelectFields("detail_id, email, name")).
		Where("kelas_id = ?", kelasID).Find(&kelasMahasiswaList).Error; err != nil {
		return entities.Presensi{}, err
	}

	presensiMahasiswaBatch := make([]entities.PresensiMahasiswa, len(kelasMahasiswaList))

	for i, mahasiswa := range kelasMahasiswaList {
		presensiMahasiswaBatch[i] = entities.PresensiMahasiswa{
			PresensiID:  id,
			MahasiswaID: *mahasiswa.Mahasiswa.DetailID,
		}
	}

	if err := tx.WithContext(ctx).Create(&presensiMahasiswaBatch).Error; err != nil {
		return entities.Presensi{}, err
	}

	var data entities.Presensi
	if err := tx.WithContext(ctx).Preload("PresensiMahasiswa.Mahasiswa", helpers.SelectFields("detail_id, email, name")).Where("id = ?", id).First(&data).Error; err != nil {
		return entities.Presensi{}, err
	}

	return data, nil
}

func (r *presensiRepository) CreatePresensiPegawai(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	timeNow := types.DateOnly(time.Now())
	var count int64
	if err := tx.WithContext(ctx).Model(&entities.Presensi{}).Where("created_at = ? AND tipe = ?", timeNow, "pegawai").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("Presensi pegawai date (%v) already exist", timeNow.String())
	}

	presensiID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	presensi := entities.Presensi{
		ID:   presensiID,
		Tipe: "pegawai",
	}
	if err := tx.WithContext(ctx).Create(&presensi).Error; err != nil {
		return err
	}

	excludedRoles := []string{
		constants.ROLE_MAHASISWA,
		constants.ROLE_SUPER_ADMIN,
		constants.ROLE_ADMIN_AKADEMIK,
		constants.ROLE_ADMIN_PEGAWAI,
		constants.ROLE_ADMIN_MAHASISWA,
		constants.ROLE_ADMIN_KEUANGAN,
	}

	var pegawaiList []entities.User
	if err := tx.WithContext(ctx).
		Joins("JOIN roles ON roles.id = users.role_id").
		Select("users.detail_id, users.name, users.email").
		Not("roles.name", excludedRoles).
		Find(&pegawaiList).Error; err != nil {
		return err
	}

	PresensiPegawaiBatch := make([]entities.PresensiPegawai, len(pegawaiList))

	for i, pegawai := range pegawaiList {
		PresensiPegawaiBatch[i] = entities.PresensiPegawai{
			PresensiID: presensiID,
			PegawaiID:  *pegawai.DetailID,
		}
	}

	if err := tx.WithContext(ctx).Create(&PresensiPegawaiBatch).Error; err != nil {
		return err
	}

	return nil
}

func (r *presensiRepository) UpdatePresensi(ctx context.Context, tx *gorm.DB, req any) error {
	if tx == nil {
		tx = r.db
	}

	switch v := req.(type) {
	case dto.PresensiPegawaiUpdateRequest:
		var entity entities.Presensi
		if err := tx.WithContext(ctx).Select("id").Where("created_at = ?", v.Date).First(&entity).Error; err != nil {
			return err
		}
		queryBatch := make([]entities.PresensiPegawai, len(v.Detail))
		for i, detail := range v.Detail {
			queryBatch[i] = entities.PresensiPegawai{
				PresensiID: entity.ID,
				PegawaiID:  detail.DetailID,
				Status:     detail.Status,
			}
		}
		if err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "presensi_id"}, {Name: "pegawai_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status"}),
		}).Create(&queryBatch).Error; err != nil {
			return err
		}

	case dto.PresensiMahasiswaUpdateRequest:
		queryBatch := make([]entities.PresensiMahasiswa, len(v.Detail))
		for i, detail := range v.Detail {
			queryBatch[i] = entities.PresensiMahasiswa{
				PresensiID:  v.PresensiID,
				MahasiswaID: detail.DetailID,
				Status:      detail.Status,
			}

		}
		if err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "presensi_id"}, {Name: "mahasiswa_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status"}),
		}).Create(&queryBatch).Error; err != nil {
			return err
		}

	}
	return nil
}

func (r *presensiRepository) GetPresensiMahasiswa(ctx context.Context, tx *gorm.DB, presensiID uuid.UUID) (entities.Presensi, error) {
	if tx == nil {
		tx = r.db
	}

	var presensi entities.Presensi
	if err := tx.WithContext(ctx).Preload("PresensiMahasiswa.Mahasiswa", helpers.SelectFields("detail_id, email, name")).Where("id = ?", presensiID).First(&presensi).Error; err != nil {
		return entities.Presensi{}, err
	}

	return presensi, nil
}

func (r *presensiRepository) GetPresensiPegawai(ctx context.Context, tx *gorm.DB, filter any) (entities.Presensi, error) {
	if tx == nil {
		tx = r.db
	}

	var presensi entities.Presensi
	query := tx.WithContext(ctx).
		Preload("PresensiPegawai.Pegawai", helpers.SelectFields("detail_id, email, name,created_at, updated_at")).Where("tipe = ?", "pegawai")

	switch v := filter.(type) {
	case uuid.UUID:
		query.Where("id = ?", v)
	case types.DateOnly:
		query.Where("created_at = ?", v)
	default:
		return entities.Presensi{}, fmt.Errorf("Filter by ID presensi or created_at")
	}

	if err := query.First(&presensi).Error; err != nil {
		return entities.Presensi{}, err
	}

	return presensi, nil
}

func (r *presensiRepository) GetAllPresensiPegawai(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entities.Presensi, int64, error) {
	if tx == nil {
		tx = r.db
	}

	var total int64
	if err := tx.WithContext(ctx).Model(&entities.Presensi{}).Where("tipe = ?", "pegawai").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var presensi []entities.Presensi
	query := tx.WithContext(ctx).
		Preload("PresensiPegawai.Pegawai", helpers.SelectFields("detail_id, email, name,created_at, updated_at")).Where("tipe = ?", "pegawai").Offset(offset).Limit(limit)

	if err := query.Find(&presensi).Error; err != nil {
		return nil, 0, err
	}

	return presensi, total, nil
}

func (r *presensiRepository) CountPresensi(ctx context.Context, tx *gorm.DB, tipe *string) (int64, error) {
	if tx == nil {
		tx = r.db
	}

	var total int64
	if err := tx.WithContext(ctx).Model(&entities.Presensi{}).Where("tipe = ?", &tipe).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
