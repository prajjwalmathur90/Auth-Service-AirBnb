package db

import (
	"AuthInGo/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	GetByID() (*models.User, error)
	Create() (error)
}

type UserReposityImpl struct {
	db *sql.DB
}

func NewUserRepository(_db *sql.DB) UserRepository {
	return &UserReposityImpl{
		db: _db,
	}
}

func (u *UserReposityImpl) Create() (error) {
	query := "INSERT INTO users (username, email, password) VALUES (? , ?, ?)"

	result, err := u.db.Exec(query, "testUser", "test@example.com", "test123")

	if err != nil {
		fmt.Println("Error creating user : ", err)
		return err
	}

	rowsAffected, rowErr := result.RowsAffected()

	if rowErr != nil {
		fmt.Println("Error getting rows affected : ", rowErr)
		return rowErr
	}

	if(rowsAffected == 0){
		fmt.Println("No rows affected!")
		return nil
	}

	fmt.Println("User created successfully! Rows Affected : ", rowsAffected)

	return nil
}

func (u *UserReposityImpl) GetByID() (*models.User, error) {
	
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?"

	row := u.db.QueryRow(query, 1)

	user := &models.User{}

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User Not Found")
			return nil, err
		} else {
			fmt.Println("Error scanning user : ", err)
			return nil, err
		}
	}

	fmt.Println("User fetched successfully!", user)
	
	return user, nil
}
