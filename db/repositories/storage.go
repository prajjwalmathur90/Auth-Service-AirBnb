package db

type Storage struct { // facilitates dependency injection for repositories
	UserRepository UserRepository
}

func NewStorage() *Storage {
	return &Storage{
		UserRepository: &UserReposityImpl{},
	}
}