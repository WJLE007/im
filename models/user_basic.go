package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UserBasic struct {
	// "account": "账号",
	//    "password": "密码",
	//    "nickname": "昵称",
	//    "sex": 1, //0是未知1是男2是女
	//    "email": "邮箱",
	//    "avatar": "头像",
	//    "creat_at": 1, //创建时间
	//    "updated_at": 1 // 更新时间
	Identity string `bson:"identity"`
	Account  string `bson:"account"`
	//Password string  `bson:"password"`
	Nickname  string `bson:"nickname"`
	Sex       int    `bson:"sex"`
	Email     string `bson:"email"`
	Avatar    string `bson:"avatar"`
	CreatAt   int64  `bson:"creat_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

func (UserBasic) CollectionName() string {
	return "user_basic"
}

func GetUserBasicByAP(account, password string) (*UserBasic, error) {
	ub := new(UserBasic)
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"account", account}, {"password", password}}).
		Decode(ub)
	return ub, err
}

func GetUserBasicById(identity string) (*UserBasic, error) {
	ub := new(UserBasic)
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"identity", identity}}).
		Decode(ub)
	return ub, err
}
func GetUserBasicCountByEmail(email string) (int64, error) {
	return Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{{"email", email}})

}
