package swagger

type (
	// USER
	UserAdminCreateRequest struct {
		Name     string  `json:"name" form:"name" binding:"required,min=2,max=255" example:"rezi // required, min 2 characters, max 255 characters"`
		Email    string  `json:"email" form:"email" binding:"required,email" example:"rezi@example.com // required, must be a valid email address"`
		Password string  `json:"password" form:"password" binding:"required,min=8" example:"inipasswordrezi // required, min 8 characters"`
		RoleName string  `json:"role_name" form:"role_kode" binding:"required" example:"raja-nyawit // required, must be a valid role name"`
		DetailId *string `json:"detail_id" form:"detail_id" binding:"omitempty" example:"01965a1d-7777-7777-7777-777777777777"`
	}

	UserNonAdminCreateRequest struct {
		Name     string  `json:"name" form:"name" binding:"required,min=2,max=255" example:"Rezi // required, min 2 max 255 characters"`
		Email    string  `json:"email" form:"email" binding:"required,email" example:"rezi@example.com // required, must be a valid email address"`
		Password string  `json:"password" form:"password" binding:"required,min=8" example:"inipasswordrezi // required, min 8 characters"`
		RoleName string  `json:"role_name" form:"role_kode" binding:"required,is_non_admin" example:"raja-nyawit // required, must be a valid role name"`
		DetailId *string `json:"detail_id" form:"detail_id" binding:"required" example:"01965a1d-7777-7777-7777-777777777777"`
	}

	UserAdminUpdateRequest struct {
		Name     string  `json:"name" form:"name" binding:"omitempty,min=2,max=255" example:"Rezi // optional, min 2 max 255 characters"`
		Email    string  `json:"email" form:"email" binding:"omitempty,email" example:"rezi@example.com // optional, must be a valid email address"`
		Password string  `json:"password" form:"password" binding:"omitempty,min=8" example:"inipasswordrezi // optional, min 8 characters"`
		RoleName string  `json:"role_name" form:"role_name" binding:"omitempty" example:"raja-nyawit // optional"`
		DetailId *string `json:"detail_id" form:"detail_id" binding:"omitempty" example:"01965a1d-7777-7777-7777-777777777777"`
	}

	UserNonAdminUpdateRequest struct {
		Name     string  `json:"name" form:"name" binding:"omitempty,min=2,max=255" example:"rezi"`
		Email    string  `json:"email" form:"email" binding:"omitempty,email" example:"rezi@example.com // optional, must be a valid email address"`
		Password string  `json:"password" form:"password" binding:"omitempty,min=8" example:"inipasswordrezi // optional, min 8 characters"`
		RoleName string  `json:"role_name" form:"role_name" binding:"omitempty,is_non_admin" example:"raja-nyawit // optional"`
		DetailId *string `json:"detail_id" form:"detail_id" binding:"omitempty" example:"01965a1d-7777-7777-7777-777777777777"`
	}

	// AKADEMIK
	AkademikCreateRequest struct {
		ID           uint   `json:"id" binding:"required,gte=0" example:"20241"`
		TipeSemester string `json:"tipe_semester" binding:"required,enumTipeSemester" example:"genap"`
		TahunAwal    string `json:"tahun_awal" binding:"required" example:"2024-01-01"`
		TahunAkhir   string `json:"tahun_akhir" binding:"required" example:"2025-01-01"`
	}

	AkademikUpdateRequest struct {
		ID           uint   `json:"id" binding:"omitempty,gte=0" example:"20241"`
		TipeSemester string `json:"tipe_semester" binding:"omitempty,enumTipeSemester" example:"genap"`
		TahunAwal    string `json:"tahun_awal" binding:"omitempty" example:"2024-01-01"`
		TahunAkhir   string `json:"tahun_akhir" binding:"omitempty" example:"2025-01-01"`
		Status       string `json:"status" binding:"omitempty,enumStatus" example:"aktif"`
	}
)
