package dto

type UserListRequest struct {
	// Cari nama/email/username
	Search   string `form:"search"` 
	// Filter admin/parent   
	Role     string `form:"role"`    
	// Filter status
	IsActive *bool  `form:"is_active"` 
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}
