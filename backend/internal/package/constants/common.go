package constants

import (
	"errors"
	"time"
	"web-hosting/internal/package/env"
)

const (
	DB         = "db"
	DB_LOG_DIR = "./config/logs/query_log"

	JWTService      = "JWTService"
	JWT_ISSUER      = "TIM 1"
	JWT_ACCESS_EXP  = 1 * time.Hour
	JWT_REFRESH_EXP = 24 * time.Hour

	ROLE_ID_SUPER_ADMIN = 1
	ROLE_ID_ADMIN       = 2
	ROLE_ID_MAHASISWA   = 3
	ROLE_ID_DOSEN       = 4

	//! Nama role harus huruf kecil dan tanpa spasi di golang
	ROLE_SUPER_ADMIN = "superadmin"
	ROLE_ADMIN       = "admin"
	ROLE_MAHASISWA   = "mahasiswa"
	ROLE_DOSEN       = "dosen"
)

var (
	JWT_SECRET_KEY = env.GetWithDefault[string]("JWT_SECRET", "")

	ErrInternalErr = errors.New("Internal Error")

	MESAGE_FAILED_GET_DATA_FROM_BODY = "failed to get data from body"
)
