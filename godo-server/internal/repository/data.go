package repository

import "godo/internal/repository/entities"

var (
	user1 = entities.Person{Name: "Ash"}
	user2 = entities.Person{Name: "Misty"}

	qaTag = entities.Tag{Name: "QA"}
	onHoldTag = entities.Tag{Name: "On hold"}
	easyWinTag = entities.Tag{Name: "Easy win"}

	allTags = []entities.Tag{qaTag, onHoldTag, easyWinTag}

	tasks1 = []entities.Task{
		{
			Name:        "Database structure design",
			Description: "Design the structure of the database",
			Type:        0,
			Status:      0,
			Creator: user2,
		},
		{
			Name:        "Database implementation",
			Description: "Implement the actual database in code",
			Type:        0,
			Status:      0,
			Creator: user2,
		},
	}

	tasks2 = []entities.Task{
		{
			Name:        "UI/UX Design",
			Description: "Design the user interface and UX considerations",
			Type:        0,
			Status:      0,
			Creator: user1,
		},
	}

	stories = []entities.Story{
		{Name: "Building user authentication", Tasks: tasks1, Creator: user2},
		{Name: "User interface and UX", Tasks: tasks2, Creator: user1},
	}

	project = entities.Project{
		Name: "TheMainSolution_Dev",
		Description: "This is the description of this project/solution",
		Stories: stories,
		Tags: allTags,
		Creator: user1,
	}
)