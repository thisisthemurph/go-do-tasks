package repository

import (
	"fmt"
	"godo/configuration"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"

	"github.com/jinzhu/gorm"
)

func CreateAndPopulateDatabase(logger ilog.StdLogger) {
	db := connect(logger)

	dropAllTables(db)
	migrate(db)
	populateTestData(db)
}

func GetDatabase(logger ilog.StdLogger) *gorm.DB {
	return connect(logger)
}

func connect(logger ilog.StdLogger) *gorm.DB {
	c := configuration.LoadConfig(logger)
	connectionString := makeConnectionString(c)

	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		panic(err)
	}

	return db
}

func makeConnectionString(config configuration.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseHost,
		config.DatabaseUsername,
		config.DatabasePassword,
		config.DatabaseName,
	)
}

func dropAllTables(db *gorm.DB) {
	db.DropTableIfExists(&entities.Account{})
	db.DropTableIfExists(&entities.User{})
	db.DropTableIfExists(&entities.Project{})
	db.DropTableIfExists(&entities.Story{})
	db.DropTableIfExists(&entities.Task{})
	db.DropTableIfExists(&entities.Tag{})
	db.DropTableIfExists("task_tags")
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&entities.Account{})
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Project{})
	db.AutoMigrate(&entities.Story{})
	db.AutoMigrate(&entities.Task{})
	db.AutoMigrate(&entities.Tag{})
}

func populateTestData(db *gorm.DB) {
	db.Create(&account)

	user.AccountId = account.ID
	db.Create(&user)

	project.Creator = user
	db.Create(&project)

	story1.Creator = user
	story1.Project = project
	db.Create(&story1)

	task1.Creator = user
	task1.Story = story1
	db.Create(&task1)

	tag1.Project = project
	db.Create(&tag1)
}
