package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// I love sets.
type Set struct {
	list map[int]struct{}
}

func (s *Set) Has(v int) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set) Add(v int) {
	s.list[v] = struct{}{}
}

func (s *Set) AddMulti(list ...int) {
	for _, v := range list {
		s.Add(v)
	}
}

func (s *Set) Remove(v int) {
	delete(s.list, v)
}

func (s *Set) Clear() {
	s.list = make(map[int]struct{})
}

func (s *Set) Size() int {
	return len(s.list)
}

func NewSet() *Set {
	s := &Set{}
	s.list = map[int]struct{}{}
	return s
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readToLineArray(path string) []string {
	file, err := os.Open(path)
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var outArray []string = []string{}
	for scanner.Scan() {
		outArray = append(outArray, scanner.Text())
	}

	checkErr(scanner.Err())

	return outArray
}

// PART ONE
type Card struct {
    CardNumber int
    WinningNumbers []int
    WinningNumbersSet Set
    Numbers []int
}

func (card Card)String() string {
    return fmt.Sprintf("%v\n%v\n%v", card.CardNumber, card.Numbers, card.WinningNumbers)
}

func parseCardDelims(r rune) bool {
    return r == ':' || r == '|'
}

func parseCard(line string) Card {
    parts := strings.FieldsFunc(line, parseCardDelims)
    winningNumbers := strings.Split(parts[1], " ")
    numbers := strings.Split(parts[2], " ")

    card := Card{ WinningNumbersSet: *NewSet() }
    for _, number := range winningNumbers {
        num, err := strconv.Atoi(number)
        if err == nil {
            card.WinningNumbers = append(card.WinningNumbers, num)
            card.WinningNumbersSet.Add(num)
        }
    }
    for _, number := range numbers {
        num, err := strconv.Atoi(number)
        if err == nil {
            card.Numbers = append(card.Numbers, num)
        }
    }

    return card
}

func parseCards(lines []string) []Card {
    cards := []Card{}
    
    for i, line := range lines {
        cards = append(cards, parseCard(line))
        cards[i].CardNumber = i+1
    }

    return cards
}

func PartOne(input []string) int {
    cards := parseCards(input)
    points := 0

    for _, card := range cards {
        cardPoints := 0
        for _, cardNumber := range card.Numbers {
            if card.WinningNumbersSet.Has(cardNumber) {
                cardPoints *= 2
                if cardPoints == 0 {
                    cardPoints = 1
                }
            }
        }
        points += cardPoints
    }

    return points 
}

func PartTwo(input []string) int {
    cards := parseCards(input)

    // At a glance this is a memoization problem.
    // Since we already accelerated the results of the
    //  first step with sets, we may not need it?
    // NOPE, Go threw out signal: killed. Guess we 
    //  can't be that lazy. :v

    // Iterate over the cards from back to front, since the 
    //  challenge states that cards will never add cards 
    //  beyond the last. As we move backwards we can add 
    //  the cumulative sum to cards, since every card after 
    //  will already have been processed.
    cardSums := []int{}
    for i := 0; i < len(cards); i++ {
        cardSums = append(cardSums, 0)
    }

    for i := len(cards)-1; i >= 0; i-- {
        card := cards[i]
        cardPoints := 0

        for _, cardNumber := range card.Numbers {
            if card.WinningNumbersSet.Has(cardNumber) {
                cardPoints += 1
            }
        }

        // Add sums of all cards this card would create
        for j := 0; j < cardPoints; j++ {
            cardSums[i] += cardSums[i+j+1]
        }
        // Don't forget to add this card itself.
        cardSums[i] += 1
    }

    // Dis is too slow.
    //// We need to use a for loop to iterate over a 
    ////  growing slice, range will only evaluate the 
    ////  length on the initial iteration, this will 
    ////  re-evaluate every iteration.
    //for i := 0; i < len(cardStack); i++ {
    //    card := cardStack[i]
    //    cardPoints := 0
    //    for _, cardNumber := range card.Numbers {
    //        if card.WinningNumbersSet.Has(cardNumber) {
    //            cardPoints += 1
    //        }
    //    }
    //    // Now we win a copy of the next <cardPoints> cards.
    //    // This is based on the card number. We can use the
    //    //  original stack to pull the correct ones.
    //    for j := 0; j <= cardPoints; j++ {
    //        cardStack = append(cardStack, cards[card.CardNumber-1+j])
    //    }
    //}
    fmt.Println(cardSums)

    sum := 0
    for _, cardSum := range cardSums {
        sum += cardSum
    }

    return sum
}

func main() {
	input := readToLineArray("input")
    
    firstAnswer := PartOne(input)
    fmt.Println("Part One::")
    fmt.Println(firstAnswer) // 20667

    secondAnswer := PartTwo(input)
    fmt.Println("Part Two::")
    fmt.Println(secondAnswer) // 5833065
}
