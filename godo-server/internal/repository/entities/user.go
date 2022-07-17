package entities

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name          string `json:"name" validate:"required,min=1,max=25"`
	Username      string `json:"username" gorm:"index" validate:"required"`
	Discriminator uint32 `json:"discriminator" validate:"required"`
	Email         string `json:"email" gorm:"unique,index" validate:"required"`
	Password      string `json:"-" validate:"required"`
	AccountId     string `json:"-" validate:"required"`
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

func (u *User) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	err := enc.Encode(u)
	if err != nil {
		log.Println("Issue encoding User JSON", err)
	}

	return err
}

func (u *User) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(u)
	if err != nil {
		log.Println("Issue decoding User JSON:", err)
	}

	return err
}

func (u User) FromHttpRequest(r *http.Request) {
	log.Println("Decoding", u)
	u = r.Context().Value(UserKey{}).(User)
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}