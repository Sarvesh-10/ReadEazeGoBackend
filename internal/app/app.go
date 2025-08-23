package app

import (
	"database/sql"

	"github.com/Sarvesh-10/ReadEazeBackend/config"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain/repository"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/services"
	util "github.com/Sarvesh-10/ReadEazeBackend/utility"
	redis "github.com/redis/go-redis/v9"
)

type App struct {
	DB          *sql.DB
	Logger      util.Logger
	UserRepo    domain.UserRepository
	UserService *UserService
	UserHandler *UserHandler
	BookRepo    domain.BookRepositoryImpl
	BookService *BookService
	BookHandler *BookHandler
	ChatService *ChatService
	ChatHandler *ChatHandler
}

func NewApp() *App {
	logger := util.NewLogger()
	logger.Info("Creating a new App")
	dsn := config.AppConfig.DBURL + "?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Error connecting to database: %s", err.Error())
		return nil
	}

	opt, err := redis.ParseURL(config.AppConfig.REDIS_URL)
	if err != nil {
		panic("Invalid Redis URL: " + err.Error())
	}
	cache := redis.NewClient(opt)
	cacheRepo := repository.NewRedisBookCache(cache, logger)
	refreshRepo := domain.NewRefreshTokenRepository(db, logger)
	userRepo := domain.NewUserRepository(db, logger)
	userService := NewUserService(userRepo, logger, refreshRepo)
	userHandler := NewUserHandler(userService, logger)

	bookRepo := domain.NewBookRepository(db)
	bookService := NewBookService(bookRepo, logger, cacheRepo)
	userBookProfileService := services.NewUserBookProfileService(bookRepo, logger, cacheRepo)
	bookHandler := NewBookHandler(bookService, userBookProfileService, logger)

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
