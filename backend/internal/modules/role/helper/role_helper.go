package helper

import (
	"web-hosting/internal/package/constants"
	"web-hosting/internal/package/helpers"
)

func RoleNameToRoleID(roleName string) uint {
	roleName = helpers.NormalizeString(roleName)
	switch roleName {
	case constants.ROLE_SUPER_ADMIN:
		return constants.ROLE_ID_SUPER_ADMIN
	case constants.ROLE_ADMIN:
		return constants.ROLE_ID_ADMIN
	case constants.ROLE_MAHASISWA:
		return constants.ROLE_ID_MAHASISWA
	case constants.ROLE_DOSEN:
		return constants.ROLE_ID_DOSEN
	default:
		return 0
	}
}
