package router

import (
	"net/http"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/app"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterBookRoutes(router *mux.Router, bookHandler *app.BookHandler) {

	router.Handle("/upload", middleware.JWTMiddleWare(http.HandlerFunc(bookHandler.UploadBook))).Methods("POST")
	router.Handle("/books/{id}", middleware.JWTMiddleWare(http.HandlerFunc(bookHandler.GetBook))).Methods("GET")
	router.Handle("/books", middleware.JWTMiddleWare(http.HandlerFunc(bookHandler.GetBooksMeta))).Methods("GET")
}
