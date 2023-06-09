package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func CreateRoom() string {
	code := generateRoomCode()
	newRoom := Room{
		Code: code,
	}
	RoomsCache = append(RoomsCache, newRoom)

	return code
}

func JoinRoom(code string, id string) {
	foundRoomIndex := FindRoomIndex(code)
	foundPlayerIndex := FindPlayerIndex(id)
	if foundRoomIndex != -1 {
		if len(RoomsCache[foundRoomIndex].PlayersID) <= 6 {
			RoomsCache[foundRoomIndex].PlayersID = append(RoomsCache[foundRoomIndex].PlayersID, id)
			PlayersCache[foundPlayerIndex].CurrentRoomCode = code
		} else {
			fmt.Println("房间" + code + "人数已达上限")
		}
	} else {
		fmt.Println("没有找到房间" + code)
	}
}

func ExitFromRoom(id string) {
	foundPlayerIndex := FindPlayerIndex(id)
	if foundPlayerIndex != -1 {
		foundRoomIndex := FindRoomIndex(PlayersCache[foundPlayerIndex].CurrentRoomCode)
		if foundRoomIndex != -1 {
			PlayersCache[foundPlayerIndex].CurrentRoomCode = ""
			RoomsCache[foundRoomIndex].PlayersID = removeStringFromArray(RoomsCache[foundRoomIndex].PlayersID, id)
			if len(RoomsCache[foundRoomIndex].PlayersID) == 0 {
				DeleteRoom(RoomsCache[foundRoomIndex].Code)
			}
		}
	}
}

func FindRoomIndex(code string) int {
	for i, room := range RoomsCache {
		if room.Code == code {
			return i
		}
	}
	return -1
}

func DeleteRoom(code string) {
	indexToDelete := -1
	for i, room := range RoomsCache {
		if room.Code == code {
			indexToDelete = i
			break
		}
	}

	if indexToDelete != -1 {
		RoomsCache = append(RoomsCache[:indexToDelete], RoomsCache[indexToDelete+1:]...)
	}
}

func generateRoomCode() string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成6位随机数字
	randomNumber := rand.Intn(900000) + 100000

	randomNumberStr := strconv.Itoa(randomNumber)

	for _, room := range RoomsCache {
		if room.Code == randomNumberStr {
			randomNumberStr = generateRoomCode()
		}
	}

	return randomNumberStr
}

func removeStringFromArray(arr []string, target string) []string {
	result := []string{}
	for _, str := range arr {
		if str != target {
			result = append(result, str)
		}
	}
	return result
}
