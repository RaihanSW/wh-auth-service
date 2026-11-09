package entity

type User struct {
	UID       string `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	RoleID    int    `json:"role_id"`
	Password  string `json:"-"`
}
