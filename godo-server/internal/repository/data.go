package repository

import "godo/internal/repository/entities"

var (
	stories = []entities.Story{
		{Name: "Building user authentication"},
		{Name: "User interface and UX"},
	}

	tasks = []entities.Task{
		{
			Name:        "Database structure design",
			Description: "Design the structure of the database",
			Type:        0,
			Status:      0,
			StoryID:     1,
		},
		{
			Name:        "Database implementation",
			Description: "Implement the actual database in code",
			Type:        0,
			Status:      0,
			StoryID:     1,
		},
		{
			Name:        "UI/UX Design",
			Description: "Design the user interface and UX considerations",
			Type:        0,
			Status:      0,
			StoryID:     2,
		},
	}
)