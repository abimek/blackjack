package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type BlackJack struct {
	dealer  Player
	players []*Player
	deck    Deck
}

func main() {
	deck := create_deck(read_deck("card_deck.json"))
	players := make([]*Player, 0, 10)
	dealer := Player{}
	players = setup(players)
	print(len(players))
	start(&deck, players, &dealer)
}

func setup(players []*Player) []*Player {
	fmt.Println("Game Instructions: To hit use (H), to pass press (S: Per User)")
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
out:
	for {
		for _, p := range players {
			if stopped == len(players) {
				break out
			}
			if p.standing || p.busted() {
				continue
			}
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Printf("Player %d turn with a deck value of %d\n", p.playerNumber, p.runningTotal())
				fmt.Printf("Type (s) to stand and (h) to hit\n")
				move, err := reader.ReadString('\n')
				move = strings.TrimSuffix(move, "\r\n")
				if err != nil {
					log.Fatal("error reading in string")
				}
				time.Sleep(time.Millisecond * 500)
				if move == "h" || move == "s" {
					if move == "h" {
						p.giveCard(deck.randomCardAndRemove(), false, false)
						if p.busted() {
							fmt.Printf("Player %d has busted with a deck value of %d\n", p.playerNumber, p.runningTotal())
							stopped += 1
						}
					} else {
						p.stand()
						stopped += 1
						fmt.Printf("Player %d has stood\n", p.playerNumber)
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
	winners := make([]*Player, 0, len(players))
	for _, p := range players {
		if !p.busted() && p.runningTotal() > dealer.runningTotal() {
			winners = append(winners, p)
		}
	}
	for _, p := range winners {
		fmt.Printf("Player %d has beat the dealer with a deck value of %d\n", p.playerNumber, p.runningTotal())
		time.Sleep(time.Millisecond * 500)
	}
	if len(winners) == 0 {
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
