package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// Structs
// Card (Value -- Between ace and king, Suite -- Spade, Heart, Diamond, Club)
// Ace -- 0
type Card struct {
	Value int
	Suite string
	Rank  string
}

type Player struct {
	Hand     []*Card
	IsDealer bool
}

func (player *Player) SortHand() {
	sort.Slice(player.Hand, func(i, j int) bool {
		return player.Hand[i].Value < player.Hand[j].Value
	})
}

// func (player *Player) ShowCards(range TYPE)
func (player *Player) ShowHand() {
	fmt.Println()
	if player.IsDealer {
		fmt.Println("Dealer's Hand (Card 1 is hidden)")
		for i := 1; i < len(player.Hand); i++ {
			fmt.Printf("Card %v -- %s\n", i+1, player.Hand[i].GetCardDetails())
		}

	} else {
		fmt.Println("Player's hand")
		for i := 0; i < len(player.Hand); i++ {
			fmt.Printf("Card %v -- %s\n", i+1, player.Hand[i].GetCardDetails())
		}
	}
	fmt.Println()
}

func (player *Player) DealTo(deck *Deck) *Card {
	dealtCard := deck.DealCard()
	player.Hand = append(player.Hand, dealtCard)
	return dealtCard
}

func (player *Player) GetDealerTotal() int {
	if !player.IsDealer {
		fmt.Errorf("error: player is not a dealer")
		return 100
	}

	numOfAces := 0
	cardTotalWithoutAces := 0
	aceTotal := 0

	for _, card := range player.Hand {
		if card.Value == 0 {
			numOfAces++
		}
		cardTotalWithoutAces += card.Value
	}

	if cardTotalWithoutAces < 11 && numOfAces > 0 {
		// Hand without aces = 10
		// 1 ace
		// Ace with value of 1: 21 - 11 = 10
		// Ace with value of 11: 21 - 21 = 0
		// Pick closest hand to 0 that isn't over 0
		aceDiffWithValue1 := 21 - (cardTotalWithoutAces + numOfAces)
		aceDiffWithValue11 := 21 - (cardTotalWithoutAces + numOfAces + 10)

		if aceDiffWithValue11 < aceDiffWithValue1 && aceDiffWithValue11 > -1 {
			aceTotal += 11
		} else {
			aceTotal += 1
		}
	}

	if numOfAces > 1 {
		return cardTotalWithoutAces + aceTotal + numOfAces - 1
	} else {
		return cardTotalWithoutAces + aceTotal
	}
}

func (player *Player) DealAndGetTotalAndCard(deck *Deck) (int, *Card) {

	dealtCard := player.DealTo(deck)

	if player.IsDealer {
		return player.GetDealerTotal(), dealtCard
	} else {
		return player.GetPlayerTotal(), dealtCard
	}
}

func (player *Player) GetPlayerTotal() int {

	numOfAces := 0
	cardTotalWithoutAces := 0
	firstAceValue := 1

	for _, card := range player.Hand {
		if card.Rank == "Ace" {
			numOfAces++
		}
		cardTotalWithoutAces += card.Value
	}

	if cardTotalWithoutAces < 11 && numOfAces > 0 {
		isValid := false
		for !isValid {
			fmt.Println("Card total without first ace:", strconv.Itoa(cardTotalWithoutAces+numOfAces-1))

			var input string

			fmt.Print("Value for first ace (1 or 11): ")
			fmt.Scanln(&input)

			value, err := strconv.Atoi(input)
			if err != nil || (value != 1 && value != 11) {
				fmt.Println("Invalid Input!")
				continue
			}

			firstAceValue += value
			isValid = true
		}
	}

	if numOfAces >= 1 {
		if numOfAces > 1 {
			return cardTotalWithoutAces + firstAceValue + numOfAces - 1
		} else {
			return cardTotalWithoutAces + firstAceValue
		}
	} else {
		return cardTotalWithoutAces
	}
}

func (card *Card) PrintCardDetails() {
	fmt.Println(card.GetCardDetails())
}

func (card *Card) GetCardDetails() string {
	return fmt.Sprintf("%s of %ss", card.Rank, card.Suite)
}

type Deck struct {
	Cards []*Card
}

func (deck *Deck) PrintCards() {
	for _, card := range deck.Cards {
		card.PrintCardDetails()
	}
}

func (deck *Deck) GenerateCards() {
	suites := [4]string{"spade", "heart", "diamond", "club"}
	for _, suite := range suites {
		deck.generateCardsBySuite(suite)
	}
}

func (deck *Deck) ShuffleCards() {
	rand.Shuffle(len(deck.Cards), func(i, j int) { deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i] })
}

func (deck *Deck) DealCard() *Card {
	// Get the last index of cards from the deck
	lastIndex := len(deck.Cards) - 1
	// Temporary value of the last element (which is a card)
	dealtCard := deck.Cards[lastIndex]
	// Remove the card that we just pulled from slice
	deck.Cards = deck.Cards[:len(deck.Cards)-1]
	// Return the last card
	//fmt.Printf("Dealt Card -- Rank: %s, Suite: %s\n", dealtCard.Rank, dealtCard.Suite)
	return dealtCard
}

func (deck *Deck) GetCardsLeft() int {
	return len(deck.Cards)
}

func (deck *Deck) generateCardsBySuite(suite string) {
	rankToValue := map[string]int{
		"Ace":   0,
		"Two":   2,
		"Three": 3,
		"Four":  4,
		"Five":  5,
		"Six":   6,
		"Seven": 7,
		"Eight": 8,
		"Nine":  9,
		"Ten":   10,
		"Jack":  10,
		"Queen": 10,
		"King":  10,
	}

	for rank, value := range rankToValue {
		card := Card{Value: value, Suite: suite, Rank: rank}
		deck.Cards = append(deck.Cards, &card)
	}
}
func PlayBlackJack(deck *Deck) {
	playerTotal := 0
	dealerTotal := 0

	// Deal first hand
	player := Player{IsDealer: false}
	player.DealTo(deck)
	playerTotal, _ = player.DealAndGetTotalAndCard(deck)
	player.ShowHand()
	fmt.Println("Player total after initial deals:", strconv.Itoa(playerTotal))

	fmt.Println()

	dealer := Player{IsDealer: true}
	dealer.DealTo(deck)
	dealerTotal, _ = dealer.DealAndGetTotalAndCard(deck)
	dealer.ShowHand()

	isPlaying := true
	for playerTotal <= 21 && isPlaying {
		fmt.Println()
		isPlaying = playerHasHit(&player)
		if isPlaying {
			newTotal, dealtCard := player.DealAndGetTotalAndCard(deck)
			playerTotal = newTotal
			fmt.Printf("Player receives: %s\n", dealtCard.GetCardDetails())
			fmt.Println("Player current Total:", strconv.Itoa(playerTotal))
		}
		if playerTotal > 21 {
			fmt.Println()
			fmt.Println("Player busts! Dealer wins.")
			break
		}
	}
	fmt.Println("Player Final Total:", strconv.Itoa(playerTotal))
	if playerTotal > 21 {
		fmt.Println()
		fmt.Println("Player busts! Dealer wins.")
		return
	}

	fmt.Println()
	fmt.Println("Dealer hidden card:", dealer.Hand[0].GetCardDetails())
	fmt.Println("Dealer current total:", strconv.Itoa(dealerTotal))

	for dealerTotal < 17 {
		newTotal, dealtCard := dealer.DealAndGetTotalAndCard(deck)
		dealerTotal = newTotal
		fmt.Printf("Dealer receives: %s\n", dealtCard.GetCardDetails())
		fmt.Println("Dealer new total:", strconv.Itoa(dealerTotal))
	}
	fmt.Println("Dealer final total:", strconv.Itoa(dealerTotal))

	fmt.Println()
	if dealerTotal > 21 {
		fmt.Println("Dealer Busts! Player wins.")
		return
	}

	fmt.Println()
	fmt.Printf("Player Total: %v, Dealer Total: %v\n", playerTotal, dealerTotal)
	if playerTotal > dealerTotal {
		fmt.Println("Player Wins!")
	} else if playerTotal < dealerTotal {
		fmt.Println("Dealer Wins!")
	} else {
		fmt.Println("Draw!")
	}
}

// func playerHasHit() bool
func playerHasHit(player *Player) bool {
	isValidInput := false
	var input string
	for !isValidInput {

		fmt.Print("Stay, hit, or view hand? [stay/hit/view-hand]: ")
		_, err := fmt.Scanln(&input)
		if err != nil || (input != "stay" && input != "hit" && input != "view-hand") {
			fmt.Errorf(err.Error())
			fmt.Println("Invalid Input.")
			continue
		}

		if input == "view-hand" {
			player.ShowHand()
			continue
		}

		isValidInput = true
	}

	if input == "hit" {
		return true
	} else {
		return false
	}
}

// func
func initDeck() *Deck {
	deck := Deck{}
	deck.GenerateCards()
	deck.ShuffleCards()

	return &deck
}
func main() {
	deck := initDeck()
	PlayBlackJack(deck)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
