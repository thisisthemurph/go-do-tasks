package auth

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

type TokenRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenRequestKey struct{}

func (t *TokenRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	err := e.Decode(t)
	if err != nil {
		log.Println("Auth TokenRequest: there was an error decoding the JSON:", err)
	}

	return err
}

func (t TokenRequest) FromHttpRequest(r *http.Request) {
	t = r.Context().Value(TokenRequestKey{}).(TokenRequest)
}

func (t *TokenRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
