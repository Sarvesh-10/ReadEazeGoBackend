package domain

import (
	"database/sql"
	"fmt"

	"github.com/Sarvesh-10/ReadEazeBackend/utility"
)

type UserRepository interface {
	CreateUser(user *User) (int64, error)
	GetUserByEmail(username string) (*User, error)
}

type UserRepositoryImpl struct {
	db     *sql.DB
	logger *utility.Logger
}

func NewUserRepository(db *sql.DB, logger *utility.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{db, logger}
}

func (ur *UserRepositoryImpl) CreateUser(user *User) (int64, error) {
	var id int64
	err := ur.db.QueryRow(
		`INSERT INTO users (name, password, email) 
         VALUES ($1, $2, $3) RETURNING id`,
		user.Name, user.Password, user.Email,
	).Scan(&id)

	if err != nil {
		ur.logger.Error("Error inserting user: %s", err.Error())
		return -1, err
	}

	return id, nil
}

// func for getting user by email
func (ur *UserRepositoryImpl) GetUserByEmail(email string) (*User, error) {
	ur.logger.Info("Getting user by email: %s", email)
	row := ur.db.QueryRow("SELECT id,name,password, email FROM users WHERE email = $1", email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found, return nil to indicate no user was found
			ur.logger.Info("No user found with email: %s", email)
			return nil, nil
		}
		ur.logger.Error("Error getting user: %s", err.Error())
		return nil, err
	}
	fmt.Printf("USER IS %+v", user)

	return user, nil
}
