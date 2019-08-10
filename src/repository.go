package monitorservice

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"github.com/google/uuid"
	"fmt"
)

type Repository struct {
	mongo *mongo.Client
}

func NewRepository (mongoClient *mongo.Client) *Repository {
	return &Repository {
		mongoClient,
	}
}

func (r *Repository) CreateLog (log Log) ( Log, error){
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := r.mongo.Database("monitor_service").Collection("app_logs")
	id, _ := uuid.NewUUID()
	_, err := collection.InsertOne(ctx, bson.M{
		"_id":			 id.String(),
		"name":          log.Name,
		"userId": 		log.UserId,
		"content":       log.Content,
		"_type":         log.Type,
		"occuredAt":	 time.Now(),
	})

	if err != nil {
		fmt.Println("Error: ", err)
		return Log{}, err
	}

	return Log{
		Id: id,
		Name: log.Name,
		Content: log.Content,
		Type: log.Type,
		OccurredAt: time.Now(),
		UserId: log.UserId,
	}, nil
}

func (r *Repository) FindLog (id uuid.UUID) ( interface {}, error ) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := r.mongo.Database("monitor_service").Collection("app_logs")

	var log Log
	err := collection.FindOne(ctx, bson.M{ "_id": id.String() }).Decode(&log)

	if err != nil {
		return nil, err
	}

	log.Id = id
	return log, nil
}