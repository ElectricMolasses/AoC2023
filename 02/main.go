package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PART ONE

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

type Game struct {
	GameId int
	Rounds []Round
}

type Round struct {
	Red   int
	Green int
	Blue  int
}

func (r Round) Validate(allowedTotals Round) bool {
	if r.Red > allowedTotals.Red {
		return false
	}
	if r.Blue > allowedTotals.Blue {
		return false
	}
	if r.Green > allowedTotals.Green {
		return false
	}

	return true
}

func (s Round) Max(o Round) Round {
	maximums := Round{}

	if s.Red >= o.Red {
		maximums.Red = s.Red
	} else {
		maximums.Red = o.Red
	}
	if s.Green >= o.Green {
		maximums.Green = s.Green
	} else {
		maximums.Green = o.Green
	}
	if s.Blue >= o.Blue {
		maximums.Blue = s.Blue
	} else {
		maximums.Blue = o.Blue
	}

	return maximums
}

func (s Round) Power() int {
	return s.Red * s.Green * s.Blue
}

func parseRound(input string) Round {
	var round Round = Round{}

	pieces := strings.Split(input, ",")
	for _, piece := range pieces {
		// Retrieve the number and parse to int
		trimmed := strings.Trim(piece, " ")
		num, err := strconv.Atoi(strings.Split(trimmed, " ")[0])
		checkErr(err)

		// Assign the number to the relevant color field
		if strings.Contains(piece, "red") {
			round.Red = num
			continue
		}
		if strings.Contains(piece, "blue") {
			round.Blue = num
			continue
		}
		if strings.Contains(piece, "green") {
			round.Green = num
			continue
		}
	}

	return round
}

func parseGame(input string) Game {
	var game Game = Game{}

	// Break up the string by Game ID:, roundData, roundData, ...
	// Split by delims : and ;
	slices := strings.FieldsFunc(input, func(r rune) bool {
		return r == ':' || r == ';'
	})

	// Assign the gameId
	gameId, err := strconv.Atoi(strings.Split(slices[0], " ")[1])
	checkErr(err)
	game.GameId = gameId

	// Process rounds and append them to the game objects rounds.
	for _, round := range slices[1:] {
		game.Rounds = append(game.Rounds, parseRound(round))
	}

	return game
}

func parseInput(input []string) []Game {
	var games []Game = []Game{}

	for _, line := range input {
		games = append(games, parseGame(line))
	}

	return games
}

func validateGame(game Game, allowedTotals Round) bool {
	for _, round := range game.Rounds {
		if !round.Validate(allowedTotals) {
			return false
		}
	}

	return true
}

func findOnlyValidGames(games []Game, allowedTotals Round) []Game {
	validGames := []Game{}

	for _, game := range games {
		if validateGame(game, allowedTotals) {
			validGames = append(validGames, game)
		}
	}

	return validGames
}

func sumGameIds(games []Game) int {
	sum := 0

	for _, game := range games {
		sum += game.GameId
	}

	return sum
}

// PART TWO
func findGameMinimum(game Game) Round {
	minimums := Round{
		Red:   0,
		Green: 0,
		Blue:  0,
	}

	for _, round := range game.Rounds {
		minimums = minimums.Max(round)
	}

	return minimums
}

func findGamesMinimums(games []Game) []Round {
	minimums := []Round{}

	for _, game := range games {
		minimums = append(minimums, findGameMinimum(game))
	}

	return minimums
}

func sumPowersOfRounds(rounds []Round) int {
	sum := 0

	for _, round := range rounds {
		sum += round.Power()
	}

	return sum
}

func main() {
	// Read each game and validate whether or not the games are possible when
	//  checked against a round representing the number of each colour in
	//  each game.
	input := readToLineArray("input")
	allowedTotals := Round{
		Red:   12,
		Green: 13,
		Blue:  14,
	}

	parsedGames := parseInput(input)

	validGames := findOnlyValidGames(parsedGames, allowedTotals)

	fmt.Println("Part 1 Answer::")
	fmt.Println(sumGameIds(validGames)) // => 2348

	// For part two find the fewest number of cubes possible for each game to be possible.
	gameMinimums := findGamesMinimums(parsedGames)
	fmt.Println("Minimums")
	for i, game := range gameMinimums {
		fmt.Printf("%v:%v", i, game)
	}

	fmt.Println("Part 2 Answer::")
	fmt.Println(sumPowersOfRounds(gameMinimums)) // => 76008
}
