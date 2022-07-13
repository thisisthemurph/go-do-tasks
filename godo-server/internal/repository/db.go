package repository

import (
	"fmt"
	"godo/internal/auth"
	"godo/internal/config"
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
	configuration := config.LoadConfig(logger)
	connectionString := makeConnectionString(configuration)

	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		panic(err)
	}

	return db
}

func makeConnectionString(config config.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseHost,
		config.DatabaseUsername,
		config.DatabasePassword,
		config.DatabaseName,
	)
}

func dropAllTables(db *gorm.DB) {
	// db.DropTableIfExists(&entities.Person{})
	db.DropTableIfExists(&auth.User{})
	db.DropTableIfExists(&entities.Project{})
	db.DropTableIfExists(&entities.Story{})
	db.DropTableIfExists(&entities.Task{})
	db.DropTableIfExists(&entities.Tag{})
	db.DropTableIfExists("task_tags")

}

func migrate(db *gorm.DB) {
	// db.AutoMigrate(&entities.Person{})
	db.AutoMigrate(&auth.User{})
	db.AutoMigrate(&entities.Project{})
	db.AutoMigrate(&entities.Story{})
	db.AutoMigrate(&entities.Task{})
	db.AutoMigrate(&entities.Tag{})
}

func populateTestData(db *gorm.DB) {
	db.Create(&project)
}
