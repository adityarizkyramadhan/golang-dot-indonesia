package dto

type UserRegister struct {
	Username *string `json:"username" binding:"required,min=8,max=50"`
	Password *string `json:"password" binding:"required,min=8,max=50"`
	Name     *string `json:"name" binding:"required,min=1,max=50"`
}

type UserLogin struct {
	Username *string `json:"username" binding:"required,min=8,max=50"`
	Password *string `json:"password" binding:"required,min=8,max=50"`
}
