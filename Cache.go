package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

var CardsCache []Card
var RoomsCache []Room
var PlayersCache []Player

func init() {
	LoadCards()
	for _, card := range CardsCache {
		fmt.Println(card.Name)
	}
}

func LoadCards() {
	filePath := "cards.csv" // CSV文件路径

	// 打开CSV文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("无法打开CSV文件:", err)
		return
	}
	defer file.Close()

	// 创建CSV阅读器
	reader := csv.NewReader(file)

	// 读取CSV记录
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("读取CSV记录时出错:", err)
		return
	}

	// 将CSV记录转化为对象
	var cards []Card
	for _, record := range records {
		// 假设CSV文件的列顺序为Name,Age,Email
		name := record[0]
		number, _ := strconv.Atoi(record[1])
		saneEffect := record[2]
		counts, _ := strconv.Atoi(record[3])
		insaneEffect := record[4]
		inDeck := record[5]

		// 将记录转化为Person对象
		card := Card{
			Name:         name,
			Number:       number,
			SaneEffect:   saneEffect,
			Counts:       counts,
			InsaneEffect: insaneEffect,
			Insane:       insaneEffect != "",
			InDeck:       inDeck == "1",
		}

		cards = append(cards, card)
	}

	CardsCache = cards
}
