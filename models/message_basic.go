package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MessageBasic struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (MessageBasic) CollectionName() string {
	return "message_basic"
}

func InsertOneMessageBasic(mb *MessageBasic) error {
	_, err := Mongo.Collection(MessageBasic{}.CollectionName()).InsertOne(context.Background(), mb)
	return err
}
func GetMessageListByRoomID(roomidentity string, limt, skip *int64) ([]*MessageBasic, error) {
	cursor, err := Mongo.Collection(MessageBasic{}.CollectionName()).Find(context.Background(), bson.M{"room_identity": roomidentity}, &options.FindOptions{
		Limit: limt,
		Skip:  skip,
		Sort:  bson.D{{"created_at", -1}},
	})
	if err != nil {
		log.Println("[DB ERROR]")
		return nil, err
	}
	messageBasics := make([]*MessageBasic, 0)
	for cursor.Next(context.Background()) {
		mb := new(MessageBasic)
		err := cursor.Decode(mb)
		if err != nil {
			return nil, err
		}
		messageBasics = append(messageBasics, mb)
	}

	return messageBasics, err
}
