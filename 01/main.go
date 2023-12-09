package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Input answers in case I decide to perform a much needed refactor
//  for practice later:
//   Part 1 Answer::
//   55816
//   Part 2 Answer::
//   54980

// Simple start. The challenge is to read each line, find
// the first and last number in each line, and concatenate
// them to create a two digit number. The puzzle answer is
// the sum of all these two digit numbers.

// Functions written for part one

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

func extractNumbersFromLine(line string) []int {
	var numbers []int = []int{}

	for _, char := range line {
		if num, err := strconv.Atoi(string(char)); err == nil {
			numbers = append(numbers, num)
		}
	}

	return numbers
}

func extractNumbersFromAllLines(lines []string) [][]int {
	var number_lines [][]int = [][]int{}

	for _, line := range lines {
		number_lines = append(number_lines, extractNumbersFromLine(line))
	}

	return number_lines
}

// Note this will create a two digit number if there is only one
//
//	number provided in the array. This seems to be the desired
//	result for the challenge
func concatFirstAndLast(array []int) int {
	if len(array) == 0 {
		panic("Impossible concenation, array is empty")
	}

	output, err := strconv.Atoi(strconv.Itoa(array[0]) + strconv.Itoa(array[len(array)-1]))
	checkErr(err)

	return output
}

func concatFirstAndLastAllLines(array [][]int) []int {
	var output []int = []int{}

	for _, line := range array {
		output = append(output, concatFirstAndLast(line))
	}

	return output
}

func sumArray(array []int) int {
	var sum int = 0

	for _, v := range array {
		sum += v
	}

	return sum
}

var wordDict = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

// Functions written for part two
func upgradeExtractNumbersFromLine(line string) []int {
	var output []int = []int{}

	for iLine, char := range line {
		if num, err := strconv.Atoi(string(char)); err == nil {
			output = append(output, num)
			continue
		}

		for word, num := range wordDict {
			if len(word) > len(line[iLine:]) {
				continue
			}
			if word == line[iLine:iLine+len(word)] {
				output = append(output, num)
			}
		}
	}

	return output
}

func upgradeExtractNumbersFromAllLines(lines []string) [][]int {
	var output [][]int = [][]int{}

	for _, line := range lines {
		output = append(output, upgradeExtractNumbersFromLine(line))
	}

	return output
}

func main() {
	// Read input file
	input := readToLineArray("input")

	// PART ONE

	// Parse numbers from each line and store them
	//  in an array. Then store this array in an
	//  array for all lines.
	extractedNumbers := extractNumbersFromAllLines(input)

	// String concatonate the first and last integers
	//  from each line
	concatonatedLines := concatFirstAndLastAllLines(extractedNumbers)

	// Summed lines for the answer
	sum := sumArray(concatonatedLines)

	fmt.Println("Part 1 Answer::")
	fmt.Println(sum)

	// PART TWO

	// All we need to change for part two is to update
	//  the extractNumbersFromAllLines function
	extractedNumbers = upgradeExtractNumbersFromAllLines(input)
	concatonatedLines = concatFirstAndLastAllLines(extractedNumbers)

	sum = sumArray(concatonatedLines)

	fmt.Println("Part 2 Answer::")
	fmt.Println(sum)
}
