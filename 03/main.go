package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// This problem could be solved VERY easily by setting any
//  non-numeric character as a delimiter and then converting
//  all remaining values to ints and summing them. Solving it
//  the hard way as a grid problem because it's more fun :)

// Go doesn't have sets? O.o

type Set struct {
	list map[rune]struct{}
}

func (s *Set) Has(v rune) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set) Add(v rune) {
	s.list[v] = struct{}{}
}

func (s *Set) AddMulti(list ...rune) {
	for _, v := range list {
		s.Add(v)
	}
}

func (s *Set) Remove(v rune) {
	delete(s.list, v)
}

func (s *Set) Clear() {
	s.list = make(map[rune]struct{})
}

func (s *Set) Size() int {
	return len(s.list)
}

func NewSet() *Set {
	s := &Set{}
	s.list = map[rune]struct{}{}
	return s
}

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

// ...
// ...

func parseStringArrToGrid(input []string) [][]rune {
	var grid [][]rune = [][]rune{}

    // This will read each horizontal line sequentially,
    //  which will make accessing the slice [y][x].
	for _, line := range input {
		grid = append(grid, []rune(line))
	}

	return grid
}

// A part number is any number adjacent to a symbol
func isSymbol(r rune) bool {
	return !notSymbols.Has(r)
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

type Position struct {
	Line int
	Char int
}

func (c Position) String() string {
	return fmt.Sprintf("%v:%v", c.Line, c.Char)
}

func (c Position) WithinBounds(min Position, max Position) bool {
	return c.Line >= min.Line && c.Line <= max.Line &&
		c.Char >= min.Char && c.Char <= max.Char
}

func findSymbols(grid [][]rune) []Position {
    coords := []Position{}

    for lineNum, line := range grid {
        for y, c := range line {
            if !notSymbols.Has(c) {
                coords = append(coords, Position{
                    Line: lineNum,
                    Char: y,
                })
            }
        }
    }

	return coords
}

func findSpecifiedSymbols(grid [][]rune, symbols Set) []Position {
    coords := []Position{}

    for lineNum, line := range grid {
        for y, c := range line {
            if symbols.Has(c) {
                coords = append(coords, Position{
                    Line: lineNum,
                    Char: y,
                })
            }
        }
    }

	return coords
}

func ScanAdjacentNumbers(grid [][]rune, pos Position) []Position {
    adjacentNumbers := []Position{}

    lowerBound := Position{ Line: 0, Char: 0 }
    upperBound := Position{ Line: len(grid), Char: len(grid[0]) }

    for i := -1; i < 2; i++ {
        for j := -1; j < 2; j++ {
            if i == 0 && j == 0 { continue }
            currentPosition := Position{
                Line: pos.Line+i,
                Char: pos.Char+j,
            }

            if currentPosition.WithinBounds(lowerBound, upperBound) {
                if isNumber(grid[currentPosition.Line][currentPosition.Char]) {
                    adjacentNumbers = append(adjacentNumbers, currentPosition)
                }
            }
        }
    }

    return adjacentNumbers
}

func ExpandNumber(grid [][]rune, pos Position, checkedPositions map[string]struct{}) int {
    leftBound := pos
    rightBound := pos

    // Find the edges of the number pos lands on
    for leftBound.Char-1 >= 0 && isNumber(grid[leftBound.Line][leftBound.Char-1]) {
        leftBound = Position{
            Line: leftBound.Line,
            Char: leftBound.Char - 1,
        }
    }
    for rightBound.Char+1 < len(grid[0]) && isNumber(grid[rightBound.Line][rightBound.Char+1]) {
        rightBound = Position{
            Line: rightBound.Line,
            Char: rightBound.Char + 1,
        }
    }

    currentPos := leftBound
    numString := ""
    // Collect the number, and add all positions the number inhabits to checkedPositions
    for currentPos.Char <= rightBound.Char {
        numString += string(grid[currentPos.Line][currentPos.Char])
        checkedPositions[currentPos.String()] = struct{}{}
        currentPos.Char += 1
    }

    number, err := strconv.Atoi(numString)
    if err != nil {
        panic(err)
    }
    return number
}

func PartOne(grid [][]rune) int {
    // Find all symbols, and store them as coordinates in a list.
    positions := findSymbols(grid)
    numbers := []int{}

    // Iterate over the symbols and find all numbers touching each.
    // Make sure we do NOT collect any given number more than once,
    //  it's possible for two symbols to touch one number.
    checkedPositions := map[string]struct{}{}
    for _, pos := range positions {
        adjacent := ScanAdjacentNumbers(grid, pos)
        
        for _, newPos := range adjacent {
            if _, ok := checkedPositions[newPos.String()]; !ok {
                numbers = append(numbers, ExpandNumber(grid, newPos, checkedPositions))
           } 
        }
    }
    

    // Sum all found numbers and return the result.
    sum := 0
    for _, number := range numbers {
        sum += number
    }

    return sum
}

func PartTwo(grid [][]rune) int {
    numbers := []int{}
    // Find all *'s and store them as coordinates in a list.
    gearSet := NewSet()
    gearSet.Add('*')

    gearPositions := findSpecifiedSymbols(grid, *gearSet)

    // Iterate over all gears to find all numbers touching each.
    // For each gear with two numbers touching it, multiply those
    //  numbers together and add them to the list to be summed.
    for _, pos := range gearPositions {
        checkedPositions := map[string]struct{}{}
        adjacent := ScanAdjacentNumbers(grid, pos)
        currentNumbers := []int{}
        for _, newPos := range adjacent {
            if _, ok := checkedPositions[newPos.String()]; !ok {
                currentNumbers = append(currentNumbers, ExpandNumber(grid, newPos, checkedPositions))
            }
        }
        if len(currentNumbers) == 2 {
            numbers = append(numbers, currentNumbers[0] * currentNumbers[1])
        }
    }

    sum := 0
    for _, number := range numbers {
        sum += number
    }

    return sum
}

var notSymbols *Set = NewSet()

func main() {
	notSymbols.AddMulti('1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '.')
	input := readToLineArray("input")

    // For the purpose of this challenge x/y orientation
    //  DOES matter, the numbers are read as normal,
    //  left to right, and summed as appropriate.
    grid := parseStringArrToGrid(input)
    partOneAnswer := PartOne(grid)

    // PART ONE!
    fmt.Println("Part One::")
    fmt.Println(partOneAnswer) // 525,119

    partTwoAnswer := PartTwo(grid)
    fmt.Println("Part Two::")
    fmt.Println(partTwoAnswer) // 76,504,829
}
