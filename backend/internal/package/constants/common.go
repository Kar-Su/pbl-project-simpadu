package constants

const (
	DB         = "db"
	DB_LOG_DIR = "./config/logs/query_log"
	JWTService = "JWTService"

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
