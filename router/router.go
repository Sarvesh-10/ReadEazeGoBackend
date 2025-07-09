package router

import (
	"net/http"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/app"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/gorilla/mux"
)

// SetupRoutes initializes all routes
func SetupRoutes(userHandler *app.UserHandler, bookHandler *app.BookHandler, chatHandler *app.ChatHandler, logger *utility.Logger) *mux.Router {
	r := mux.NewRouter()
	logger.Info("Setting up routes")

	// User routes
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if err := userHandler.SignUpHandler(w, r); err != nil {
			logger.Error("Error signing up user: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := userHandler.LoginHandler(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("POST")

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		userHandler.LogoutHandler(w, r)
	}).Methods("POST")

	r.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if err := userHandler.RefreshTokenHandler(w, r); err != nil {
			logger.Error("Error refreshing token: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("POST")

	RegisterBookRoutes(r, bookHandler)

	r.HandleFunc("/chat", chatHandler.HandleChat).Methods("POST")

	logger.Info("✅ Routes setup complete!")
	return r
}
