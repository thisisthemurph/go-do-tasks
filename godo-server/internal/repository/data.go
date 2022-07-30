package repository

import (
	"godo/internal/repository/entities"
)

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

	story1 = entities.Story{
		Name:        "Story name",
		Description: "This is only a test story",
	}

	task1 = entities.Task{
		Name:        "Task 1",
		Description: "This is only a task",
	}

	tag1 = entities.Tag{
		Name: "Tag 1",
	}
)
