package models

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Mongo = InitMongo()
var Redis = InitRedis()

func InitMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://47.103.139.14:27017"))
	if err != nil {
		log.Println("Connect to Mongo Error", err)
	}
	return client.Database("im")
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "47.103.139.14:6379",
		Password: "Ll5211314",
		DB:       0,
	})
}
