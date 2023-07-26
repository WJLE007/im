package define

type MessageStruct struct {
	RoomIdentity string `json:"room_identity"`
	Message      string `json:"message"`
}

var RegisterPer = "TOKEN_"
var ExpireTime = 300
