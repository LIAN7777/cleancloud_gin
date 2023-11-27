package dto

type RegisterForm struct {
	Telephone    string `json:"telephone"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RegisterCode string `json:"registerCode"`
}
