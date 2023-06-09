package main

import "time"

type Status string

const (
	Unselectable Status = "Unselectable"
	Undead       Status = "Undead"
	Non          Status = "Non"
)

type Player struct {
	ID              string    `json:"openid"`
	Name            string    `json:"name"`
	HandCard        []Card    `json:"handCard"`
	FrontCards      []Card    `json:"frontCards"`
	Insane          int       `json:"insane"`
	SaneToken       int       `json:"saneToken"`
	InsaneToken     int       `json:"insaneToken"`
	LastLogin       time.Time `json:"lastLogin"`
	CurrentRoomCode string    `json:"currentRoomCode"`
	Status          Status    `json:"status"`
}

type Room struct {
	Code      string   `json:"code"`
	Log       []Log    `json:"log"`
	PlayersID []string `json:"playersOpenid"`
	Deck      Deck     `json:"deck"`
}

type Card struct {
	Name         string `json:"name"`
	Insane       bool   `json:"insane"`
	SaneEffect   string `json:"saneEffect"`
	InsaneEffect string `json:"insaneEffect"`
	Counts       int    `json:"counts"`
	Number       int    `json:"number"`
	InDeck       bool   `json:"inDeck"`
}

type Deck struct {
	Cards          []Card `json:"cards"`
	CardsOutOfGame []Card `json:"cardOutOfGame"`
}

type Log struct {
	Info       string `json:"info"`
	Time       string `json:"time"`
	UserOpenid string `json:"userOpenid"`
}
