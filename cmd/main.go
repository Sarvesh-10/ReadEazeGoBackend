package main

import (
	"net/http"

	"github.com/Sarvesh-10/ReadEazeBackend/config"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/app"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/services"
	"github.com/Sarvesh-10/ReadEazeBackend/router"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// func init() {
// 	err := godotenv.Load() // Only needed for local dev
// 	if err != nil {
// 		panic("Error loading .env file")
// 	}
// }

var (
	DB_USER     = config.AppConfig.DBUser
	DB_PASSWORD = config.AppConfig.DBPassword
	DB_HOST     = config.AppConfig.DBHost
	DB_PORT     = config.AppConfig.DBPort
	DB_NAME     = config.AppConfig.DBName
)

func main() {
	logger := utility.NewLogger()

	app := app.NewApp()
	defer app.DB.Close()
	r := router.SetupRoutes(app.UserHandler, app.BookHandler, app.ChatHandler, logger)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // React frontend address
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	logger.Info("Starting server on port 8080")
	go services.ListenStatusQueue(app.Cache)

	error := http.ListenAndServe(":8080", handler)
	if error != nil {
		logger.Error("Error starting server: %s", error.Error())
	}

}
