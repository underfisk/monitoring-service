package monitorservice

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserRepository struct {
	mongo *mongo.Client
}

type User struct {
	Id string	`bson:"_id"`
	Username string
	Password []byte
	Email string
	CreatedAt string
	JobPosition string /** what he does */
	Company string
	Avatar string
}

type RegisterDto struct {
	Username string	`json:username`
	Password string	`json:password`
	Email string	`json:email`
	JobPosition string	`json:jobPosition`
	Company string	`json:company`

}

func NewUserRepository (mongo *mongo.Client) *UserRepository {
	return &UserRepository{
		mongo: mongo,
	}
}


func (ur *UserRepository) Create (data *RegisterDto) (uuid.UUID, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := ur.mongo.Database("monitor_service").Collection("users")
	id, _ := uuid.NewUUID()
	pwd, hashErr := bcrypt.GenerateFromPassword( []byte(data.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return uuid.UUID{}, hashErr
	}


	_, err := collection.InsertOne(ctx, bson.M{
		"_id":			 	id.String(),
		"username":         data.Username,
		"password": 		pwd,
		"email":       	data.Email,
		"job_position":         	data.JobPosition,
		"avatar":	 	"",
		"company":		data.Company,
		"createdAt":	time.Now(),
	})

	if err != nil {
		log.Println("Error: ", err)
		return uuid.UUID{}, err
	}

	return id, nil
}

func (ur *UserRepository) GetById (id uuid.UUID) (User, error) {
	return User{}, nil
}

func (ur *UserRepository) GetByEmail(email string) (User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := ur.mongo.Database("monitor_service").Collection("users")

	var usr User
	err := collection.FindOne(ctx, bson.M{ "email": email }).Decode(&usr)

	if err != nil {
		return User{}, err
	}

	return User{
		Id: usr.Id,
		Username: usr.Username,
		Password: usr.Password,
		Company: usr.Company,
		JobPosition: usr.JobPosition,
		CreatedAt: usr.CreatedAt,
		Avatar: usr.Avatar,
		Email: usr.Email,
	}, nil
}

func (ur *UserRepository) Update (id uuid.UUID, field string, value interface {}) {

}