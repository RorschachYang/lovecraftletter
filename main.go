package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 配置静态文件路径
	router.Static("/static", "./static")

	// 配置 WebSocket 路由
	router.GET("/ws", handleWebSocket)
	router.GET("/log", handleLog)
	router.POST("/createRoom", handleCreateRoom)
	router.POST("/joinRoom", handleJoinRoom)
	router.GET("/connectInRoom", handleConnectInRoom)

	// 启动服务器
	log.Println("服务器已启动，监听端口 80...")
	err := router.Run(":80")
	if err != nil {
		log.Fatal("服务器启动失败：", err)
	}
}
