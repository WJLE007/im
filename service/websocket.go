package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im/define"
	"im/helper"
	"im/models"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

var wc = make(map[string]*websocket.Conn)

func WebsocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) //让http协议升级为websocket协议
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
		return
	}
	defer conn.Close()
	uc := c.MustGet("user_claims").(*helper.UserClaim)
	wc[uc.Identity] = conn
	for {
		ms := new(define.MessageStruct)
		err := conn.ReadJSON(ms) //把查询到的数据直接读成接送格式

		if err != nil {
			log.Println("Read Error", err)
			return
		}

		//用户是否是属于消息体传递的房间
		_, err = models.GetUserRoomByUserIDRoomID(uc.Identity, ms.RoomIdentity)
		if err != nil {
			log.Printf("UserIdentity:%v or RoomIdentity:%v NOT EXITS", uc.Identity, ms.RoomIdentity)
		}

		//TODO:保存消息
		mb := &models.MessageBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: ms.RoomIdentity,
			Data:         ms.Message,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		err = models.InsertOneMessageBasic(mb)
		if err != nil {
			log.Println("[DB ERROR]", err)
			return
		}
		//查询特定房间内的在线用户
		userRooms, err := models.GetUserRoomByRoomIdenrity(ms.RoomIdentity)
		if err != nil {
			log.Println("[DB ERROR]")
			return
		}
		//给房间内用户发送消息
		for _, room := range userRooms {
			if w, ok := wc[room.UserIdentity]; ok {
				err := w.WriteMessage(websocket.TextMessage, []byte(ms.Message))
				if err != nil {
					log.Println("Write error", err)
					return
				}
			}
		}
		//将消息发动给所有人
		//for _, w := range wc {
		//	err := w.WriteMessage(websocket.TextMessage, []byte(ms.Message))
		//	if err != nil {
		//		log.Println("Write error", err)
		//		return
		//	}
		//}
	}

}
