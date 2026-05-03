package swagger

type (
	UserAdminCreateRequest struct {
		Name     string `json:"name" form:"name" binding:"required,min=2,max=255" example:"rezi // required, min 2 characters, max 255 characters"`
		Email    string `json:"email" form:"email" binding:"required,email" example:"rezi@example.com // required, must be a valid email address"`
		Password string `json:"password" form:"password" binding:"required,min=8" example:"inipasswordrezi // required, min 8 characters"`
		RoleName string `json:"role_name" form:"role_kode" binding:"required" example:"raja-nyawit // required, must be a valid role name"`
		DetailId *uint  `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0" example:"1"`
	}

	UserNonAdminCreateRequest struct {
		Name     string `json:"name" form:"name" binding:"required,min=2,max=255" example:"Rezi // required, min 2 max 255 characters"`
		Email    string `json:"email" form:"email" binding:"required,email" example:"rezi@example.com // required, must be a valid email address"`
		Password string `json:"password" form:"password" binding:"required,min=8" example:"inipasswordrezi // required, min 8 characters"`
		RoleName string `json:"role_name" form:"role_kode" binding:"required,is_non_admin" example:"raja-nyawit // required, must be a valid role name"`
		DetailId *uint  `json:"detail_id" form:"detail_id" binding:"required,gt=0" example:"1"`
	}

	UserAdminUpdateRequest struct {
		Name     string `json:"name" form:"name" binding:"omitempty,min=2,max=255" example:"Rezi // optional, min 2 max 255 characters"`
		Email    string `json:"email" form:"email" binding:"omitempty,email" example:"rezi@example.com // optional, must be a valid email address"`
		Password string `json:"password" form:"password" binding:"omitempty,min=8" example:"inipasswordrezi // optional, min 8 characters"`
		RoleName string `json:"role_name" form:"role_name" binding:"omitempty" example:"raja-nyawit // optional"`
		DetailId *uint  `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0" example:"1"`
	}

	UserNonAdminUpdateRequest struct {
		Name     string `json:"name" form:"name" binding:"omitempty,min=2,max=255" example:"rezi"`
		Email    string `json:"email" form:"email" binding:"omitempty,email" example:"rezi@example.com // optional, must be a valid email address"`
		Password string `json:"password" form:"password" binding:"omitempty,min=8" example:"inipasswordrezi // optional, min 8 characters"`
		RoleName string `json:"role_name" form:"role_name" binding:"omitempty,is_non_admin" example:"raja-nyawit // optional"`
		DetailId *uint  `json:"detail_id" form:"detail_id" binding:"omitempty,gt=0" example:"1"`
	}
)
