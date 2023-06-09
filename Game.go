package main

import (
	"math/rand"
	"time"
)

func StartGame(roomCode string) {
	playerCounts := len(RoomsCache[FindRoomIndex(roomCode)].PlayersID)
	if playerCounts == 2 {

	} else {

	}
}

func EndGame(roomCode string) {

}

func CreateDeck(roomCode string) {
	foundRoomIndex := FindRoomIndex(roomCode)
	if foundRoomIndex != -1 {
		for _, card := range CardsCache {
			for i := 0; i < card.Counts; i++ {
				if card.InDeck {
					RoomsCache[foundRoomIndex].Deck.Cards = append(RoomsCache[foundRoomIndex].Deck.Cards, card)
				} else if !card.InDeck {
					RoomsCache[foundRoomIndex].Deck.CardsOutOfGame = append(RoomsCache[foundRoomIndex].Deck.CardsOutOfGame, card)
				}
			}
		}
	}
}

func ShuffleDecks(roomCode string) {
	var deck Deck
	foundRoomIndex := FindRoomIndex(roomCode)
	if foundRoomIndex != -1 {
		deck = RoomsCache[foundRoomIndex].Deck
	}

	rand.Seed(time.Now().UnixNano())

	// 从最后一个元素开始，逐个与随机位置的元素交换
	for i := len(deck.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	}

	RoomsCache[foundRoomIndex].Deck = deck
}

func DrawCard() {

}

func PutCard() {

}

func DiscardCard() {

}
