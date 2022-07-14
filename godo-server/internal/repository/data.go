package repository

import (
	"godo/internal/repository/entities"
)

var (
	tasks1 = []entities.Task{
		{
			Name:        "Database structure design",
			Description: "Design the structure of the database",
			Type:        0,
			Status:      0,
			// Creator:     user2,
		},
		{
			Name:        "Database implementation",
			Description: "Implement the actual database in code",
			Type:        0,
			Status:      0,
			// Creator:     user2,
		},
	}

	tasks2 = []entities.Task{
		{
			Name:        "UI/UX Design",
			Description: "Design the user interface and UX considerations",
			Type:        0,
			Status:      0,
			// Creator:     user1,
		},
	}

	project = entities.Project{
		Name:        "TheMainSolution_Dev",
		Description: "This is the description of this project/solution",
	}
)
