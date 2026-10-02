package db

import "fmt"

type UserRepository interface {
	Create() error
}

type UserReposityImpl struct {
	// db *sql.DB
}

func NewUserRepository() UserRepository {
	return &UserReposityImpl{
		// db: _db,
	}
}

func (u *UserReposityImpl) Create() error {
	fmt.Println("Registering user in repo")
	return nil
}
