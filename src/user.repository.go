package monitorservice

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	mongo *mongo.Client
}

type User struct {
	Id uuid.UUID
	Name string
	Email string
	CreatedAt string
	JobPosition string /** what he does */
	Company string
	Avatar string
}

func NewUserRepository (mongo *mongo.Client) *UserRepository {
	return &UserRepository{
		mongo: mongo,
	}
}


func (ur *UserRepository) Create () User {
	return User{}
}

func (ur *UserRepository) GetById (id uuid.UUID) (User, error) {
	return User{}, nil
}

func (ur *UserRepository) GetByCredentials (username string, password string) (User,error) {
	return User{}, nil
}

func (ur *UserRepository) Update (id uuid.UUID, field string, value interface {}) {

}