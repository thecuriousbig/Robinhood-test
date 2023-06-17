package dto

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	ProfileImage string `json:"profileImage"`
}

type RegisterRequest struct {
	Username string `json:"username" valid:"required,length(3|20)"`
	Password string `json:"password" valid:"required,length(6|20)"`
	Email    string `json:"email" valid:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateUserRequest struct {
	ProfileImage string `json:"profileImage"`
}
