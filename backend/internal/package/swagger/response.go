// Package swagger menyediakan tipe-tipe response contoh yang digunakan
// khusus untuk dokumentasi Swagger / OpenAPI.
// Tipe-tipe di sini TIDAK digunakan dalam logika bisnis, hanya sebagai
// referensi schema pada anotasi @Success / @Failure controller.
package swagger

// TIPE RESPONSE UMUM

// SuccessResponseNoData digunakan untuk endpoint yang sukses tanpa data tambahan.
// Contoh: Logout, Delete, Reset Password.
type SuccessResponseNoData struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"success"`
	Path    string `json:"path,omitempty" example:"/api/auth/logout"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// 401 - UNAUTHORIZED (Auth Middleware)

// ErrUnauthorizedMissingHeader dikembalikan ketika header Authorization tidak disertakan.
// Pesan: message="failed_auth", error="Authorization header missing"
type ErrUnauthorizedMissingHeader struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed_auth"`
	Error   string `json:"error" example:"Authorization header missing"`
	Path    string `json:"path,omitempty" example:"/api/auth/logout"`
}

// ErrUnauthorizedInvalidHeader dikembalikan ketika format header Authorization tidak valid (bukan "Bearer ...").
// Pesan: message="failed_auth", error="invalid authentication header"
type ErrUnauthorizedInvalidHeader struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed_auth"`
	Error   string `json:"error" example:"invalid authentication header"`
	Path    string `json:"path,omitempty" example:"/api/auth/logout"`
}

// ErrUnauthorizedInvalidToken dikembalikan ketika JWT token tidak valid atau sudah expired.
// Pesan: message="failed_auth", error="invalid token"
type ErrUnauthorizedInvalidToken struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed_auth"`
	Error   string `json:"error" example:"invalid token"`
	Path    string `json:"path,omitempty" example:"/api/auth/logout"`
}

// 403 - FORBIDDEN (Role Middleware)

// ErrForbiddenAccess dikembalikan ketika user tidak memiliki role yang diizinkan untuk mengakses endpoint.
// Pesan: message="Role anda tidak diizinkan", error="Forbidden"
type ErrForbiddenAccess struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Role anda tidak diizinkan"`
	Error   string `json:"error" example:"Forbidden"`
	Path    string `json:"path,omitempty" example:"/api/super/role"`
}

// 500 - INTERNAL SERVER ERROR

// ErrInternalServer dikembalikan ketika terjadi kesalahan yang tidak terduga di server.
// Pesan: message=(sesuai konteks), error="Internal Error"
type ErrInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Internal Error"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/auth/login"`
}

// 400 - BAD REQUEST (Parsing Body)

// ErrBadRequestBody dikembalikan ketika body request gagal di-parse atau validasi field gagal.
// Pesan: message="failed to get data from body", error=(detail validasi)
type ErrBadRequestBody struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to get data from body"`
	Error   string `json:"error" example:"Key: 'Email' Error:Field validation for 'Email' failed on the 'required' tag"`
	Path    string `json:"path,omitempty" example:"/api/auth/login"`
}

// ErrBadRequestURI dikembalikan ketika parameter URI gagal di-parse atau validasi gagal.
// Pesan: message="bad request", error=(detail validasi)
type ErrBadRequestURI struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"bad request"`
	Error   string `json:"error" example:"Key: 'RoleName' Error:Field validation for 'RoleName' failed on the 'is_non_admin' tag"`
	Path    string `json:"path,omitempty" example:"/api/user/sync/admin-mahasiswa/10"`
}

// ErrBadRequestRoleURI dikembalikan ketika parameter role_name di URI gagal di-parse.
// Pesan: message="failed to validate role uri", error=(detail validasi)
type ErrBadRequestRoleURI struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to validate role uri"`
	Error   string `json:"error" example:"Key: 'RoleName' Error:Field validation for 'RoleName' failed on the 'required' tag"`
	Path    string `json:"path,omitempty" example:"/api/super/role/"`
}

// AUTH MODULE - ERROR RESPONSES

// ErrLoginFailed dikembalikan ketika login gagal akibat email tidak ditemukan atau password salah.
//
// Kemungkinan error:
//   - message="failed to login user", error="user not found"         -> email tidak terdaftar
//   - message="failed to login user", error="crypto/bcrypt: ..."    -> password tidak cocok
type ErrLoginFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to login user"`
	Error   string `json:"error" example:"user not found"`
	Path    string `json:"path,omitempty" example:"/api/auth/login"`
}

// ErrLoginInternalServer dikembalikan ketika login gagal akibat error internal server.
// Pesan: message="failed to login user", error="Internal Error"
type ErrLoginInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to login user"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/auth/login"`
}

// ErrLogoutFailed dikembalikan ketika proses logout gagal.
// Pesan: message="failed logout", error="Internal Error"
type ErrLogoutFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed logout"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/auth/logout"`
}

// ErrRefreshTokenExpired dikembalikan ketika refresh token sudah kedaluwarsa.
// Pesan: message="failed refresh token", error="refresh token expired"
type ErrRefreshTokenExpired struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed refresh token"`
	Error   string `json:"error" example:"refresh token expired"`
	Path    string `json:"path,omitempty" example:"/api/auth/refresh-token"`
}

// ErrRefreshTokenNotFound dikembalikan ketika refresh token tidak ditemukan di database.
// Pesan: message="failed refresh token", error="refresh token not found"
type ErrRefreshTokenNotFound struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed refresh token"`
	Error   string `json:"error" example:"refresh token not found"`
	Path    string `json:"path,omitempty" example:"/api/auth/refresh-token"`
}

// ErrRefreshTokenInternalServer dikembalikan ketika refresh token gagal akibat error server.
// Pesan: message="failed refresh token", error="Internal Error"
type ErrRefreshTokenInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed refresh token"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/auth/refresh-token"`
}

// ErrFindRefreshTokenNotFound dikembalikan ketika data refresh token tidak ditemukan berdasarkan string token.
// Pesan: message="failed find refresh token", error="refresh token not found"
type ErrFindRefreshTokenNotFound struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed find refresh token"`
	Error   string `json:"error" example:"refresh token not found"`
	Path    string `json:"path,omitempty" example:"/api/auth/refresh-token/MBG-JAYA67"`
}

// ErrFindRefreshTokenInternal dikembalikan ketika pencarian refresh token gagal akibat error server.
// Pesan: message="failed find refresh token", error="Internal Error"
type ErrFindRefreshTokenInternal struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed find refresh token"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/auth/refresh-token/MBG-JAYA67"`
}

// ErrUnauthorizedResetPassword dikembalikan ketika user mencoba reset password milik orang lain
// dan bukan Super Admin.
// Pesan: message="User unauthorized", error="You are not authorized to reset this password"
type ErrUnauthorizedResetPassword struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"User unauthorized"`
	Error   string `json:"error" example:"You are not authorized to reset this password"`
	Path    string `json:"path,omitempty" example:"/api/auth/reset-password"`
}

// ErrResetPasswordFailed dikembalikan ketika reset password gagal (misal: email tidak terdaftar).
//
// Kemungkinan error:
//   - message="failed send password reset", error="user not found"  -> email tidak terdaftar
type ErrResetPasswordFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed send password reset"`
	Error   string `json:"error" example:"user not found"`
	Path    string `json:"path,omitempty" example:"/api/auth/reset-password"`
}

// ErrResetPasswordInternalServer dikembalikan ketika reset password gagal akibat error server.
// Pesan: message="failed send password reset", error="Internal Error"
type ErrResetPasswordInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed send password reset"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/auth/reset-password"`
}

// ROLE MODULE - ERROR RESPONSES

// ErrGetRoleFailed dikembalikan ketika pengambilan data role gagal.
// Pesan: message="failed to get role", error="Internal Error"
type ErrGetRoleFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to get role"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/role"`
}

// ErrCreateRoleBody dikembalikan ketika body request pembuatan role gagal di-parse.
// Pesan: message="failed to get request", error=(detail validasi)
type ErrCreateRoleBody struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to get request"`
	Error   string `json:"error" example:"Key: 'RoleName' Error:Field validation for 'RoleName' failed on the 'required' tag"`
	Path    string `json:"path,omitempty" example:"/api/super/role"`
}

// ErrCreateRoleFailed dikembalikan ketika pembuatan role gagal (misal: role sudah ada).
//
// Kemungkinan error:
//   - message="failed to create role", error="role already exists"  -> role dengan nama sama sudah ada
type ErrCreateRoleFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to create role"`
	Error   string `json:"error" example:"role already exists"`
	Path    string `json:"path,omitempty" example:"/api/super/role"`
}

// ErrCreateRoleInternalServer dikembalikan ketika pembuatan role gagal akibat error server.
// Pesan: message="failed to create role", error="Internal Error"
type ErrCreateRoleInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to create role"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/super/role"`
}

// ErrUpdateRoleFailed dikembalikan ketika update role gagal (misal: role tidak ditemukan).
//
// Kemungkinan error:
//   - message="failed to update role", error="role not found"  -> role_name tidak ada di database
//   - message="failed to get request", error=(validasi)        -> body request tidak valid
//   - message="failed to validate role uri", error=(validasi)  -> URI parameter tidak valid
type ErrUpdateRoleFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to update role"`
	Error   string `json:"error" example:"role not found"`
	Path    string `json:"path,omitempty" example:"/api/super/role/mahasiswa"`
}

// ErrUpdateRoleInternalServer dikembalikan ketika update role gagal akibat error server.
// Pesan: message="failed to update role", error="Internal Error"
type ErrUpdateRoleInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to update role"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/super/role/mahasiswa"`
}

// ErrDeleteRoleFailed dikembalikan ketika penghapusan role gagal (misal: role tidak ditemukan).
//
// Kemungkinan error:
//   - message="failed to delete role", error="role not found"   -> role_name tidak ada di database
//   - message="failed to validate role uri", error=(validasi)   -> URI parameter tidak valid
type ErrDeleteRoleFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to delete role"`
	Error   string `json:"error" example:"role not found"`
	Path    string `json:"path,omitempty" example:"/api/super/role/mahasiswa"`
}

// ErrDeleteRoleInternalServer dikembalikan ketika penghapusan role gagal akibat error server.
// Pesan: message="failed to delete role", error="Internal Error"
type ErrDeleteRoleInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to delete role"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/super/role/mahasiswa"`
}

// USER MODULE - ERROR RESPONSES

// ErrGetUserFailed dikembalikan ketika pengambilan data user gagal.
//
// Kemungkinan error:
//   - message="failed to get user", error="user not found"  -> user tidak ada di database
//   - message="failed to get user", error="role not found"  -> role_name tidak valid
type ErrGetUserFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to get user"`
	Error   string `json:"error" example:"user not found"`
	Path    string `json:"path,omitempty" example:"/api/me"`
}

// ErrGetUserInternalServer dikembalikan ketika pengambilan data user gagal akibat error server.
// Pesan: message="failed to get user", error="Internal Error"
type ErrGetUserInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to get user"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/me"`
}

// ErrRegisterUserFailed dikembalikan ketika registrasi user gagal.
//
// Kemungkinan error:
//   - message="failed to register user", error="email already exists"  -> email sudah terdaftar
//   - message="failed to register user", error="role not found"        -> role_name tidak valid
//   - message="failed to get data from body", error=(validasi)         -> body request tidak valid
type ErrRegisterUserFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to register user"`
	Error   string `json:"error" example:"email already exists"`
	Path    string `json:"path,omitempty" example:"/api/super/user"`
}

// ErrRegisterUserInternalServer dikembalikan ketika registrasi user gagal akibat error server.
// Pesan: message="failed to register user", error="Internal Error"
type ErrRegisterUserInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to register user"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/super/user"`
}

// ErrUpdateUserFailed dikembalikan ketika update user gagal.
//
// Kemungkinan error:
//   - message="failed to update user", error="user not found"          -> user tidak ada
//   - message="failed to update user", error="email already exists"    -> email sudah digunakan user lain
//   - message="failed to get data from body", error=(validasi)         -> body request tidak valid
type ErrUpdateUserFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to update user"`
	Error   string `json:"error" example:"user not found"`
	Path    string `json:"path,omitempty" example:"/api/super/user/019748ae-beef-7abc-b123-abcdef012345"`
}

// ErrUpdateUserInternalServer dikembalikan ketika update user gagal akibat error server.
// Pesan: message="failed to update user", error="Internal Error"
type ErrUpdateUserInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to update user"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/super/user/019748ae-beef-7abc-b123-abcdef012345"`
}

// ErrUnauthorizedUpdateNonAdmin dikembalikan ketika user mencoba update profil user lain
// dan bukan Super Admin.
// Pesan: message="failed to update user", error="Unauthorized"
type ErrUnauthorizedUpdateNonAdmin struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to update user"`
	Error   string `json:"error" example:"Unauthorized"`
	Path    string `json:"path,omitempty" example:"/api/user/sync/mahasiswa/10"`
}

// ErrDeleteUserFailed dikembalikan ketika penghapusan user gagal.
//
// Kemungkinan error:
//   - message="failed to delete user", error="user not found"  -> user tidak ada di database
//   - message="failed to delete user", error="role not found"  -> role_name tidak valid
type ErrDeleteUserFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to delete user"`
	Error   string `json:"error" example:"user not found"`
	Path    string `json:"path,omitempty" example:"/api/super/user/019748ae-beef-7abc-b123-abcdef012345"`
}

// ErrDeleteUserInternalServer dikembalikan ketika penghapusan user gagal akibat error server.
// Pesan: message="failed to delete user", error="Internal Error"
type ErrDeleteUserInternalServer struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to delete user"`
	Error   string `json:"error" example:"Internal Error"`
	Path    string `json:"path,omitempty" example:"/api/super/user/019748ae-beef-7abc-b123-abcdef012345"`
}

// ErrGetListUserFailed dikembalikan ketika pengambilan daftar user gagal.
//
// Kemungkinan error:
//   - message="failed to get user", error="role not found"  -> role_name tidak valid
type ErrGetListUserFailed struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed to get user"`
	Error   string `json:"error" example:"role not found"`
	Path    string `json:"path,omitempty" example:"/api/user/role/mahasiswa"`
}
