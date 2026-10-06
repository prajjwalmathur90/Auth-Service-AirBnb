package db

import (
	"AuthInGo/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Create(username string, email string, hashedPassword string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetByUsername(username string) (*models.User, error)
	DeleteById(id int) (error)
}

type UserReposityImpl struct {
	db *sql.DB
}

func NewUserRepository(_db *sql.DB) UserRepository {
	return &UserReposityImpl{
		db: _db,
	}
}

func (u *UserReposityImpl) Create(username string, email string, hashedPassword string) (*models.User, error) {
	query := "INSERT INTO users (username, email, password) VALUES (? , ?, ?)"

	result, err := u.db.Exec(query, username, email, hashedPassword)

	if err != nil {
		fmt.Println("Error creating user : ", err)
		return  nil, err
	}

	rowsAffected, rowErr := result.RowsAffected()

	if rowErr != nil {
		fmt.Println("Error getting rows affected : ", rowErr)
		return nil, rowErr
	}

	if(rowsAffected == 0){
		fmt.Println("No rows affected!")
		return nil, nil
	}

	fmt.Println("User created successfully! Rows Affected : ", rowsAffected)

	user, err := u.GetByEmail(email)

	if err != nil {
		fmt.Println("Error fetching user by email : ", err)
		return nil, err
	}	

	return user, nil
}

func (u *UserReposityImpl) GetByID(id int) (*models.User, error) {
	
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?"

	row := u.db.QueryRow(query, id)

	user := &models.User{}

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User Not Found with this id")
			return nil, nil
		} else {
			fmt.Println("Error scanning user : ", err)
			return nil, err
		}
	}

	fmt.Println("User fetched successfully!", user)
	
	return user, nil
}

func (u *UserReposityImpl) GetByEmail(email string) (*models.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = ?"

	row := u.db.QueryRow(query, email)

	user := &models.User{}

	err := row.Scan(&user.Id, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User Not Found with this email")
			return nil, nil
		} else {
			fmt.Println("Error scanning user : ", err)
			return nil, err
		}
	}

	fmt.Println("User fetched successfully!", user)
	
	return user, nil
}

func (u *UserReposityImpl) GetAll() ([]*models.User, error) {
	query := "SELECT id, username, email, created_at, updated_at FROM users"

	rows, err := u.db.Query(query)

	if err != nil {
		fmt.Println("Error fetching users : ", err)
		return nil, err
	}

	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Println("Error scanning user : ", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows : ", err)
		return nil, err
	}
	
	return users, nil
}

func (u *UserReposityImpl) GetByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE username = ?"

	row := u.db.QueryRow(query, username)

	user := &models.User{}

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User Not Found with this username")
			return nil, nil
		} else {
			fmt.Println("Error scanning user : ", err)
			return nil, err
		}
	}

	fmt.Println("User fetched successfully!", user)
	
	return user, nil
}

func (u *UserReposityImpl) DeleteById(id int) (error) {
	query := "DELETE from users WHERE id = ?"

	result, err := u.db.Exec(query, id)

	if err != nil {
		fmt.Println("Error deleting user : ", err)
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

	fmt.Println("User deleted successfully! Rows Affected : ", rowsAffected)

	return nil
}