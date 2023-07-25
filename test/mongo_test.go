package test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im/helper"
	"im/models"
	"testing"
	"time"
)

func TesMongotFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://47.103.139.14:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	ub := new(models.UserBasic)
	err = db.Collection("user_basic").FindOne(context.Background(), bson.D{}).Decode(ub)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ub ==> ", ub)
	fmt.Println(helper.GetMd5("123465"))
}

func TestMongoFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://47.103.139.14:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")

	cursor, err := db.Collection("user_room").Find(context.Background(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	userRooms := make([]*models.UserRoom, 0)
	for cursor.Next(context.Background()) {
		ub := new(models.UserRoom)
		err := cursor.Decode(ub)
		if err != nil {
			t.Fatal(err)
		}
		userRooms = append(userRooms, ub)
	}
	for _, v := range userRooms {
		fmt.Println("UserRoom==>", v)
	}
}
