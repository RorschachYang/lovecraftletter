package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients = make(map[*websocket.Conn]string) // 存储所有已连接的客户端及其对应的 ID
)

func handleCreateRoom(c *gin.Context) {
	var requestBody map[string]interface{}
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerName, ok := requestBody["playerName"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'name' field"})
		return
	}

	playerID := Login(playerName)

	roomCode := CreateRoom()

	JoinRoom(roomCode, playerID)

	SendLog("玩家:" + playerName + "创建了房间:" + roomCode)

	// 构建响应
	response := map[string]interface{}{
		"roomCode": roomCode,
		"playerID": playerID,
	}

	c.JSON(http.StatusOK, response)
}

func handleJoinRoom(c *gin.Context) {
	var requestBody map[string]interface{}
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomCode, ok := requestBody["roomCode"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'roomCode' field"})
		return
	}

	playerName, ok := requestBody["playerName"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'name' field"})
		return
	}

	playerID := Login(playerName)

	JoinRoom(roomCode, playerID)

	SendLog("玩家:" + playerName + "加入了房间:" + roomCode)

	// 构建响应
	response := map[string]interface{}{
		"playerID": playerID,
	}
	c.JSON(http.StatusOK, response)
}

func handleConnectInRoom(c *gin.Context) {
	var requestBody map[string]interface{}
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerID, ok := requestBody["playerID"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'playerID' field"})
		return
	}

	roomCode, ok := requestBody["roomCode"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'roomCode' field"})
		return
	}

	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("升级连接为 WebSocket 失败：", err)
		return
	}

	// 生成唯一的客户端 ID
	clientID := playerID

	// 发送客户端 ID 给连接的客户端
	err = conn.WriteMessage(websocket.TextMessage, []byte("你的id是"+clientID))
	if err != nil {
		log.Println("发送客户端 ID 失败：", err)
		conn.Close()
		return
	}

	// 添加客户端到 clients 映射
	clients[conn] = clientID

	var roomClients []*websocket.Conn
	for _, room := range RoomsCache {
		if room.Code == roomCode {
			for _, playerID := range room.PlayersID {
				for client, clientPlayerID := range clients {
					if clientPlayerID == playerID {
						roomClients = append(roomClients, client)
					}
				}
			}
		}
	}

	// 当函数返回时，断开 WebSocket 连接并删除客户端
	defer func() {
		conn.Close()
		delete(clients, conn)
		DeletePlayer(playerID)
	}()

	for {
		// 读取客户端发送的消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取客户端消息失败：", err)
			break
		}

		for _, client := range roomClients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("发送消息给客户端失败：", err)
				break
			}
		}
	}
}

func handleWebSocket(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("升级连接为 WebSocket 失败：", err)
		return
	}

	// 生成唯一的客户端 ID
	clientID := generateClientID()

	// 发送客户端 ID 给连接的客户端
	err = conn.WriteMessage(websocket.TextMessage, []byte("你的id是"+clientID))
	if err != nil {
		log.Println("发送客户端 ID 失败：", err)
		conn.Close()
		return
	}

	// 添加客户端到 clients 映射
	clients[conn] = clientID

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

		// 解析消息中的 ID 和内容
		parts := strings.SplitN(string(msg), ":", 2)
		if len(parts) != 2 {
			log.Println("无效的消息格式：", string(msg))
			continue
		}

		fromClient := clients[conn]
		toClientID := parts[0]
		message := parts[1]

		if toClientID == "all" {
			// 将消息发送给所有连接的客户端
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(fromClient+"说:"+message))
				if err != nil {
					log.Println("发送消息给客户端失败：", err)
					break
				}
			}
		} else {
			// 将消息发送给特定 ID 的客户端
			found := false
			for client, id := range clients {
				if id == toClientID {
					err := client.WriteMessage(websocket.TextMessage, []byte(fromClient+"对你说"+":"+message))
					if err != nil {
						log.Println("发送消息给客户端失败:", err)
					}
					found = true
					break
				}
			}

			if !found {
				log.Println("未找到指定的客户端 ID:", toClientID)
			}
		}
	}
}

// 生成唯一的客户端 ID
func generateClientID() string {
	id := 1
	for _, clientId := range clients {
		cid, _ := strconv.Atoi(clientId)
		if id == cid {
			id++
		}
	}
	return strconv.Itoa(id)
}
