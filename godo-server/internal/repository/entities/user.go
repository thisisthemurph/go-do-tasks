package entities

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	ID            uint       `gorm:"primary_key"`
	Name          string     `json:"name" validate:"required,min=1,max=25"`
	Username      string     `json:"username" gorm:"index" validate:"required"`
	Discriminator uint32     `json:"discriminator" validate:"required"`
	Email         string     `json:"email" gorm:"unique,index" validate:"required"`
	Password      string     `json:"-" validate:"required"`
	AccountId     string     `json:"-" validate:"required"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"-" sql:"index"`
}

type UserKey struct{}

func (u *User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s}", u.ID, u.Name)
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("Issue hashing password")
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) VerifyPassword(passwordCheck string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwordCheck))
	if err != nil {
		log.Println("Issue verifying the user's password")
		return err
	}

	return nil
}
