package dto

type LoginRequestDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDto struct {
	Token string `json:"token"`
}

type RegistrationRequestDto struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	AccountId string `json:"account_id" validate:"required"`
}
