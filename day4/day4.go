package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func isFullyOverlap(pairs []string) bool {
	p1Min, _ := strconv.Atoi(strings.Split(pairs[0], "-")[0])
	p1Max, _ := strconv.Atoi(strings.Split(pairs[0], "-")[1])
	p2Min, _ := strconv.Atoi(strings.Split(pairs[1], "-")[0])
	p2Max, _ := strconv.Atoi(strings.Split(pairs[1], "-")[1])

	// Check if P1 fits in P2
	if (p1Min <= p2Min) && (p1Max >= p2Max) {
		return true
	}

	// Check if P2 fits in P1
	if (p2Min <= p1Min) && (p2Max >= p1Max) {
		return true
	}

	return false
}

func main() {

	intputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer intputFile.Close()

	scanner := bufio.NewScanner(intputFile)

	needRework := 0
	for scanner.Scan() {
		line := scanner.Text()
		pairs := strings.Split(line, ",")
		if isFullyOverlap(pairs) {
			needRework++
		}
	}

	log.Printf("Total rework needed: %v", needRework)
}
