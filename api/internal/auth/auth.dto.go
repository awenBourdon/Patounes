package auth

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
