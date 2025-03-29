package app

import (
	"database/sql"

	"github.com/Sarvesh-10/ReadEazeBackend/config"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	util "github.com/Sarvesh-10/ReadEazeBackend/utility"
)

type App struct {
	DB          *sql.DB
	Logger      util.Logger
	UserRepo    domain.UserRepository
	UserService *UserService
	UserHandler *UserHandler
	BookRepo    domain.BookRepository
	BookService *BookService
	BookHandler *BookHandler
	ChatService *ChatService
	ChatHandler *ChatHandler
}

func NewApp() *App {
	logger := util.NewLogger()
	logger.Info("Creating a new App")
	dsn := config.AppConfig.DBURL

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Error connecting to database: %s", err.Error())
		return nil
	}

	userRepo := domain.NewUserRepository(db, logger)
	userService := NewUserService(userRepo, logger)
	userHandler := NewUserHandler(userService, logger)

	bookRepo := domain.NewBookRepository(db)
	bookService := NewBookService(bookRepo, logger)
	bookHandler := NewBookHandler(bookService, logger)

	chatService := NewChatService(config.AppConfig.LlamaAPIKey, logger)
	chatHandler := NewChatHandler(chatService, logger)
	return &App{
		DB:          db,
		Logger:      *logger,
		UserRepo:    userRepo,
		UserService: userService,
		UserHandler: userHandler,
		BookRepo:    *bookRepo,
		BookService: bookService,
		BookHandler: bookHandler,
		ChatService: chatService,
		ChatHandler: chatHandler,
	}
}
