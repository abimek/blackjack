package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var reset = false

func main() {
	deck := create_deck(read_deck("card_deck.json"), 52)
	dealer := Player{}
	players := setup()
	print(len(players))
	start(&deck, players, &dealer)
}

func setup() []*Player {
	players := make([]*Player, 0, 10)
	fmt.Println("Welcome to a custom take on blackjack with infinite cards (Choose weather its best for you to reset the deck or not for the number you want)")
	fmt.Println("When resetting the deck, the number of cards in the deck will now range between 30 and 52 being randomly selected, but it will be the first n values of the deck")
	fmt.Println("Game Instructions: To hit use (H), to stand press (S), to reset the deck press (R)")
	reader := bufio.NewReader(os.Stdin)
	var v int
	for {
		fmt.Print("Select How Many Players: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error reading in string")
		}
		line = strings.TrimSuffix(line, "\r\n")
		if v, err = strconv.Atoi(line); err != nil {
			fmt.Println(err.Error())
			fmt.Println(v)
			fmt.Println("Not an integer, try again!")
		} else {
			break
		}
	}

	for i := 0; i < v; i++ {
		players = append(players, &Player{playerNumber: i + 1})
	}
	return players
}

func start(deck *Deck, players []*Player, dealer *Player) {
	dealInitial(players, dealer, deck)
	gameLoop(deck, players)
	dealerFinal(deck, dealer)
	broadcastWinners(deck, players, dealer)
}

func gameLoop(deck *Deck, players []*Player) {
	stopped := 0
	for {
		for _, p := range players {
			if stopped == len(players) {
				return
			}
			if p.standing || p.busted() {
				continue
			}
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Printf("Player %d turn with a deck value of %d\n", p.playerNumber, p.runningTotal())
				fmt.Printf("Type (s) to stand, (h) to hit, and (r) to reset the deck\n")
				move, err := reader.ReadString('\n')
				move = strings.TrimSuffix(move, "\r\n")
				if err != nil {
					log.Fatal("error reading in string")
				}
				time.Sleep(time.Millisecond * 500)
				if move == "h" || move == "s" || move == "r" {
					if move == "h" {
						p.giveCard(deck.randomCardAndRemove(), false, false)
						fmt.Printf("Deck length is %d\n", len(deck.cards))
						if p.busted() {
							fmt.Printf("Player %d has busted with a deck(len %d) value of %d\n", p.playerNumber, len(deck.cards), p.runningTotal())
							stopped += 1
						}
					} else if move == "s" {
						p.stand()
						stopped += 1
						fmt.Printf("Player %d has stood\n", p.playerNumber)
					} else {
						if reset == false {
							*deck = create_deck(read_deck("card_deck.json"), 30+rand.Intn(12))
							fmt.Printf("Player %d has used the only reset in the game and a deck length of %d was selected!\n", p.playerNumber, len(deck.cards))
							reset = true
						} else {
							fmt.Println("Reset has already been used")
							continue
						}
					}
					break
				} else {
					fmt.Println("invalid character")
				}
			}
		}
	}
}

func dealerFinal(deck *Deck, dealer *Player) {
	fmt.Printf("...Hitting Dealer with inital value of %d....\n", dealer.runningTotal())
	fmt.Println("")
	for dealer.runningTotal() < 17 {
		dealer.giveCard(deck.randomCardAndRemove(), true, true)
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Printf("Done hitting dealer with final value of %d\n", dealer.runningTotal())
}

func broadcastWinners(deck *Deck, players []*Player, dealer *Player) {
	winners := 0
	for _, p := range players {
		if !p.busted() && p.runningTotal() > dealer.runningTotal() {
			fmt.Printf("Player %d has beat the dealer with a deck value of %d\n", p.playerNumber, p.runningTotal())
			time.Sleep(time.Millisecond * 500)
			winners++
		}
	}
	if winners == 0 {
		fmt.Println("Noone has beat the dealer\n")
	}
}

func dealInitial(players []*Player, dealer *Player, deck *Deck) {
	for i := 0; i < len(players); i++ {
		players[i].giveCard(deck.randomCardAndRemove(), false, false)
		time.Sleep(time.Millisecond * 500)
	}
	dealer.giveCard(deck.randomCardAndRemove(), true, true)
	time.Sleep(time.Millisecond * 500)
	for i := 0; i < len(players); i++ {
		players[i].giveCard(deck.randomCardAndRemove(), false, false)
		time.Sleep(time.Millisecond * 500)
	}
	dealer.giveCard(deck.randomCardAndRemove(), true, false)
}
