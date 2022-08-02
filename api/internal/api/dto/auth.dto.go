package dto

// LoginRequestDto model logging in a user
// swagger:model loginRequestDto
type LoginRequestDto struct {
	// the email address of the user
	// required: true
	Email string `json:"email" validate:"required"`

	// the user's password
	// required: true
	Password string `json:"password" validate:"required"`
}

// LoginResponseDto model for returning a JWT
// swagger:model JWTTokenResponse
type LoginResponseDto struct {
	// the token authenticating the user
	Token string `json:"token"`
}

// RegistrationRequestDto model for registering a new user
// swagger:model registrationRequestDto
type RegistrationRequestDto struct {
	// required: true
	Name string `json:"name" validate:"required"`

	// required: true
	Email string `json:"email" validate:"required"`

	// required: true
	Username string `json:"username" validate:"required"`

	// required: true
	Password string `json:"password" validate:"required"`

	// the account to which the user is to be associated
	// required: true
	AccountId string `json:"account_id" validate:"required"`
}
