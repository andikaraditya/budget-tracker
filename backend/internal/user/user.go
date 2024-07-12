package user

type UserSimple struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	UserSimple
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}
