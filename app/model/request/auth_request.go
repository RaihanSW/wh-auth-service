package request

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	RoleID    int    `json:"role_id"`
	Password  string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	RoleID    int    `json:"role_id" validate:"required"`
	Password  string `json:"password"` // optional — only updated when non-empty
}
