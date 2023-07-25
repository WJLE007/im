package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType     int    `bson:"room_type"` // 房间 类型 【1-独聊房间 2-群聊房间】
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

func GetUserRoomByUserIDRoomID(userIdentity, roomIdentity string) (*UserRoom, error) {
	ub := new(UserRoom)
	err := Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"user_identity", userIdentity}, {"room_identity", roomIdentity}}).
		Decode(ub)
	return ub, err
}

func GetUserRoomByRoomIdenrity(roomIdentity string) ([]*UserRoom, error) {

	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{{"room_identity", roomIdentity}})
	if err != nil {
		return nil, err
	}
	userRooms := make([]*UserRoom, 0)
	for cursor.Next(context.Background()) {
		ub := new(UserRoom)
		err := cursor.Decode(ub)
		if err != nil {
			return nil, err
		}
		userRooms = append(userRooms, ub)
	}
	return userRooms, err
}
