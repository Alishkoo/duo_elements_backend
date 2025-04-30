package data

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

type UserResponse struct {
	Username string `json:"username"`
}
