package main

import (
	"log"
	"net/http"
	"time"

	"godo/internal/api/handler"
	"godo/internal/config"
	"godo/internal/repository"
	"godo/internal/services"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

func main() {
	config := config.LoadConfig();
	repository.CreateAndPopulateDatabase()

	dao := repository.NewDAO()
	storyService := services.NewStoryService(dao)
	handlers := handler.MakeHandlers(
		&storyService,
	)

	router := mux.NewRouter()
	router.HandleFunc("/api/story", handlers.StoryHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/story/{id:[0-9]+}", handlers.StoryHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Addr: 			"0.0.0.0:" + config.ApiPort,
		WriteTimeout: 	time.Second * 15,
		IdleTimeout:	time.Second * 15,
		Handler: 		cors.Default().Handler(router),	
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
