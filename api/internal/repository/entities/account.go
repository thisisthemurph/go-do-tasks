package entities

//
type Account struct {
	Base

	Name  string `json:"name" validate:"required" gorm:"not null"`
	Email string `json:"email" validate:"required" gorm:"not null"`

	TimestampBase
}

type AccountKey struct{}

// AccountResponse an account representing the organisation, company, or group
// swagger:response accountResponse
type AccountResponse struct {
	Body Account
}
