package models

type UpdateAuth struct {
	Username    string `json:"username" binding:"required,min=5,max=20" minLength:"5" maxLength:"20"`
	Password    string `json:"password" binding:"required,min=8,max=20" minLength:"8" maxLength:"20" format:"password"`
	OldPassword string `json:"old_password" binding:"required,min=5,max=20" minLength:"5" maxLength:"20" format:"password"`
}

type Auth struct {
	Username string `json:"username" binding:"required,min=5,max=20" minLength:"5" maxLength:"20"`
	Password string `json:"password" binding:"required,min=5,max=20" minLength:"5" maxLength:"20" format:"password"`
}
