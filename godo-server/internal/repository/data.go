package repository

import (
	"godo/internal/repository/entities"
)

func userFactory(user *entities.User, accountId string) {
	user.AccountId = accountId
}

var (
	account = entities.Account{
		Name:  "TestingPalace",
		Email: "mike@email.com.",
	}

	user = entities.User{
		Name:          "Mike",
		Username:      "mike",
		Discriminator: 1,
		Email:         "mike@email.com",
		Password:      "$2a$14$tJqsjfLC29X9edPbbOoNnuiahW1fWPPz9911TDeg3KdqszShhtpya",
	}

	project = entities.Project{
		Name:        "TheMainSolution_Dev",
		Description: "This is the description of this project/solution",
	}
)
