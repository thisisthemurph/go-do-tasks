package repository

import (
	"fmt"
	"godo/internal/config"
	"godo/internal/repository/entities"

	"github.com/jinzhu/gorm"
)

func CreateAndPopulateDatabase() {
	db := connect()
	populateDatabase(db)
}

func GetDatabase() *gorm.DB {
	return connect()
}

func connect() *gorm.DB {
	configuration := config.LoadConfig()
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

func populateDatabase(db *gorm.DB) {
	// defer db.Close();

	// Drop the tables to refresh the data
	db.DropTableIfExists(&entities.Story{})
	db.DropTableIfExists(&entities.Task{})

	// Create the tables and add the data
	db.AutoMigrate(&entities.Story{})
	db.AutoMigrate(&entities.Task{})

	for index := range stories {
		db.Create(&stories[index])
	}

	for index := range tasks {
		db.Create(&tasks[index])
	}
}