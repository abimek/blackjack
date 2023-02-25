package main

import "fmt"

type Player struct {
	playerNumber int
	standing     bool
	cards        []Card
}

func (p *Player) giveCard(c Card, dealer bool, showDealer bool) {
	p.cards = append(p.cards, c)
	if dealer {
		if showDealer {
			fmt.Printf("dealer has been dealt a %s of %s \n", c.symbol, CARD_HOUSE_REV_MAP[c.house])
			return
		}
		fmt.Println("Dealer face down")
		return
	}
	fmt.Printf("player %d has been dealt a %s of %s \n", p.playerNumber, c.symbol, CARD_HOUSE_REV_MAP[c.house])
}

func (p *Player) stand() {
	p.standing = true
	fmt.Sprintf("Plyer %d is standing with a deck value of %d", p.playerNumber, p.runningTotal())
}

func (p *Player) runningTotal() int {
	total := 0
	for _, card := range p.cards {
		total += card.value
	}
	return total
}

func (p *Player) busted() bool {
	return p.runningTotal() > 21
}
