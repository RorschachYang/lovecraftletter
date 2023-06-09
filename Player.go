package main

import (
	"math/rand"
	"strconv"
	"time"
)

func Login(playerName string) string {
	index := FindPlayerIndexByName(playerName)
	if index == -1 {
		newPlayer := Player{
			Name: playerName,
			ID:   generatePlayerID(),
		}
		PlayersCache = append(PlayersCache, newPlayer)
		return newPlayer.ID
	} else {
		return ""
	}
}

func FindPlayerIndex(id string) int {
	for i, player := range PlayersCache {
		if player.ID == id {
			return i
		}
	}
	return -1
}

func FindPlayerIndexByName(name string) int {
	for i, player := range PlayersCache {
		if player.Name == name {
			return i
		}
	}
	return -1
}

func GetPlayerByID(id string) *Player {
	for _, player := range PlayersCache {
		if player.ID == id {
			return &player
		}
	}
	return nil
}

func DeletePlayer(id string) {
	ExitFromRoom(id)

	indexToDelete := -1
	for i, player := range PlayersCache {
		if player.ID == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete != -1 {
		PlayersCache = append(PlayersCache[:indexToDelete], PlayersCache[indexToDelete+1:]...)
	}

}

func generatePlayerID() string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成6位随机数字
	randomNumber := rand.Intn(9000) + 1000

	randomNumberStr := strconv.Itoa(randomNumber)

	for _, player := range PlayersCache {
		if player.ID == randomNumberStr {
			randomNumberStr = generateRoomCode()
		}
	}

	return randomNumberStr
}
