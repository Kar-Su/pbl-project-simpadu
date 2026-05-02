package constants

import (
	"errors"
	"time"
	"web-hosting/internal/package/env"
)

const (
	DB         = "db"
	DB_TEST    = "db_test"
	DB_LOG_DIR = "internal/configs/logs/query_log"

	JWTService      = "JWTService"
	JWT_ISSUER      = "TIM 1"
	JWT_ACCESS_EXP  = 1 * time.Hour
	JWT_REFRESH_EXP = 24 * time.Hour

	//! Nama role harus huruf kecil dan tanpa spasi di golang
	ROLE_SUPER_ADMIN     = "super-admin"
	ROLE_ADMIN_AKADEMIK  = "admin-akademik"
	ROLE_ADMIN_PEGAWAI   = "admin-pegawai"
	ROLE_ADMIN_MAHASISWA = "admin-mahasiswa"
	ROLE_ADMIN_KEUANGAN  = "admin-keuangan"
	ROLE_DOSEN           = "dosen"

	EMAIL_SUPER_ADMIN     = "tim1@poliban.ac.id"
	EMAIL_ADMIN_PEGAWAI   = "tim2@poliban.ac.id"
	EMAIL_ADMIN_MAHASISWA = "tim3@poliban.ac.id"
	EMAIL_ADMIN_KEUANGAN  = "tim4@poliban.ac.id"
)

var (
	JWT_SECRET_KEY = env.GetWithDefault[string]("JWT_SECRET", "")

	ErrInternalErr = errors.New("Internal Error")

	MESAGE_FAILED_GET_DATA_FROM_BODY = "failed to get data from body"
)
