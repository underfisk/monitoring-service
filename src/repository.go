package core

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	mongo *mongo.Client
}

func NewRepository (mongoClient *mongo.Client) *Repository {
	return &Repository {
		mongoClient,
	}
}

func (r *Repository) CreateLog (log Log) ( interface {}, error){
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := r.mongo.Database("monitor_service").Collection("app_logs")
	res, err := collection.InsertOne(ctx, bson.M{
		"name": log.name,
		"applicationId": log.applicationId,
		"payload": log.payload,
		"_type": log._type,
	})

	if err != nil {
		return -1, err
	}

	id := res.InsertedID

	return id, nil
}

