package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type House int32

const (
	HEARTS House = iota
	DIAMONDS
	CLUBS
	SPADES
)

var (
	CARD_VALUE_MAP     = map[string]int{"Q": 10, "J": 10, "K": 10, "A": 11}
	CARD_HOUSE_MAP     = map[string]House{"hearts": HEARTS, "diamonds": DIAMONDS, "clubs": CLUBS, "spades": SPADES}
	CARD_HOUSE_REV_MAP = map[House]string{HEARTS: "hearts", DIAMONDS: "diamonds", CLUBS: "clubs", SPADES: "spades"}
)

func create_deck(data_deck card_deck, numCards int) Deck {
	deck := Deck{}
	for _, card := range data_deck.Cards {
		if len(deck.cards) > numCards {
			break
		}
		c := Card{}
		c.symbol = card.Value
		if v, err := strconv.Atoi(card.Value); err != nil {
			c.value = CARD_VALUE_MAP[card.Value]
		} else {
			c.value = v
		}
		c.house = CARD_HOUSE_MAP[card.Suit]
		deck.cards = append(deck.cards, c)
	}
	return deck
}

func read_deck(file string) card_deck {

	byteValue, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("Error when trying to open file: ", err)
	}
	var data_deck card_deck
	err = json.Unmarshal(byteValue, &data_deck)
	if err != nil {
		log.Fatal("error while trying to unmarshal file: ", err)
	}
	return data_deck
}

type Deck struct {
	cards []Card
}

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

func (d *Deck) randomCardAndRemove() Card {
	index := rand.Intn(len(d.cards))
	card := d.cards[index]
	remove(d.cards, index)
	return card
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(s []Card, i int) []Card {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Card struct {
	value  int
	symbol string
	house  House
}

type card_data struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
}

type card_deck struct {
	Cards []card_data `json:"cards"`
}
