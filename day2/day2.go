package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func decode(s string) string {
	switch s {
	case "A", "X":
		return "rock"
	case "B", "Y":
		return "paper"
	case "C", "Z":
		return "scissors"
	}

	return ""
}

func decodeP2(p1Move, p2EncodeMove string) string {
	switch p1Move {
	case "rock":
		switch p2EncodeMove {
		case "X": // lose
			return "scissors"
		case "Y": // draw
			return p1Move
		case "Z": // win
			return "paper"
		}
	case "paper":
		switch p2EncodeMove {
		case "X": // lose
			return "rock"
		case "Y": // draw
			return p1Move
		case "Z": // win
			return "scissors"
		}
	case "scissors":
		switch p2EncodeMove {
		case "X": // lose
			return "paper"
		case "Y": // draw
			return p1Move
		case "Z": // win
			return "rock"
		}
	}

	return ""
}

func pointDistribution(p1Hand string, p2Hand string, pointGuide map[string]int) (int, int) {
	p1 := 0
	p2 := 0

	p1Move := decode(p1Hand)
	p2Move := decodeP2(p1Move, p2Hand)

	// Draw
	if p1Move == p2Move {
		p1 = 3 + pointGuide[p1Move]
		p2 = 3 + pointGuide[p2Move]
		return p1, p2
	}

	// P1 plays Rock
	if p1Move == "rock" {
		if p2Move == "scissors" {
			// P1 Wins
			p1 = 6 + pointGuide[p1Move]
			p2 = pointGuide[p2Move]
			return p1, p2
		} else {
			// P2 Wins
			p1 = pointGuide[p1Move]
			p2 = 6 + pointGuide[p2Move]
			return p1, p2
		}
	}

	// P1 plays Paper
	if p1Move == "paper" {
		if p2Move == "rock" {
			// P1 Wins
			p1 = 6 + pointGuide[p1Move]
			p2 = pointGuide[p2Move]
			return p1, p2
		} else { // P2 Wins
			p1 = pointGuide[p1Move]
			p2 = p2 + 6 + pointGuide[p2Move]
			return p1, p2
		}
	}

	// P1 plays Scissors
	if p1Move == "scissors" {
		if p2Move == "paper" {
			// P1 Wins
			p1 = 6 + pointGuide[p1Move]
			p2 = pointGuide[p2Move]
			return p1, p2
		} else { // P2 Wins
			p1 = pointGuide[p1Move]
			p2 = 6 + pointGuide[p2Move]
			return p1, p2
		}
	}

	return p1, p2
}

// GPT optimized version
// func pointDistribution(p1Hand string, p2Hand string, pointGuide map[string]int) (int, int) {
// 	score := make(map[string]int)
// 	score["p1"] = 0
// 	score["p2"] = 0

// 	p1Move := decode(p1Hand)
// 	p2Move := decodeP2(p1Move, p2Hand)

// 	// Draw
// 	if p1Move == p2Move {
// 		score["p1"] = 3 + pointGuide[p1Move]
// 		score["p2"] = 3 + pointGuide[p2Move]
// 		return score["p1"], score["p2"]
// 	}

// 	// P1 wins
// 	if (p1Move == "rock" && p2Move == "scissors") ||
// 		(p1Move == "paper" && p2Move == "rock") ||
// 		(p1Move == "scissors" && p2Move == "paper") {
// 		score["p1"] = 6 + pointGuide[p1Move]
// 		score["p2"] = pointGuide[p2Move]
// 		return score["p1"], score["p2"]
// 	}

// 	// P2 wins
// 	if (p1Move == "rock" && p2Move == "paper") ||
// 		(p1Move == "paper" && p2Move == "scissors") ||
// 		(p1Move == "scissors" && p2Move == "rock") {
// 		score["p1"] = pointGuide[p1Move]
// 		score["p2"] = 6 + pointGuide[p2Move]
// 		return score["p1"], score["p2"]
// 	}
	
// 	return score["p1"], score["p2"]
// }

func main() {

	// Create the guide
	pointGuide := map[string]int{
		"rock":     1, // Rock
		"paper":    2, // Paper
		"scissors": 3, // Scissor
	}

	intputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer intputFile.Close()

	scanner := bufio.NewScanner(intputFile)

	totalP1 := 0
	totalP2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		turn := strings.Split(line, " ")
		p1Pts, p2Pts := pointDistribution(turn[0], turn[1], pointGuide)
		totalP1 += p1Pts
		totalP2 += p2Pts
	}

	log.Printf("P1 = %v , P2 = %v", totalP1, totalP2)
}
