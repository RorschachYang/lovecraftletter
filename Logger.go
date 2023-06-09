package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var logClients = make(map[*websocket.Conn]string)

func handleLog(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("升级连接为 WebSocket 失败：", err)
		return
	}

	// 生成唯一的客户端 ID
	clientID := generateClientID()

	// 添加客户端到 clients 映射
	logClients[conn] = clientID

	// 当函数返回时，断开 WebSocket 连接并删除客户端
	defer func() {
		conn.Close()
		delete(clients, conn)
	}()

	for {
		// 读取客户端发送的消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取客户端消息失败：", err)
			break
		}
		if string(msg) == "players" {
			ShowPlayers()
		}
		if string(msg) == "rooms" {
			ShowRooms()
		}
	}

}

func SendLog(text string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// 将消息发送给所有连接的客户端
	for client := range logClients {
		err := client.WriteMessage(websocket.TextMessage, []byte(currentTime+"  "+text))
		if err != nil {
			log.Println("发送消息给客户端失败：", err)
			break
		}
	}
}

func ShowPlayers() {
	for _, player := range PlayersCache {
		SendLog("名称:" + player.Name + ",ID:" + player.ID + ",房间:" + player.CurrentRoomCode)
	}
}

func ShowRooms() {
	for _, room := range RoomsCache {
		SendLog("房间:" + room.Code)
		for _, playerID := range room.PlayersID {
			SendLog("玩家:" + playerID)
		}
	}
}
