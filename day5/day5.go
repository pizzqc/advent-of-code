package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type GiantCargo struct {
	Cargo []Stack
}
type Stack struct {
	Crates []string
}

func (s *GiantCargo) init(nbOfStack int) {
	s.Cargo = []Stack{}
	for i := 0; i < nbOfStack; i++ {
		s.Cargo = append(s.Cargo, Stack{
			Crates: []string{},
		})
	}
}

func (s *GiantCargo) move(amount int, source int, destination int) {
	// Extract slice matching amount to move
	crateToMove := s.Cargo[source-1].Crates[len(s.Cargo[source-1].Crates)-amount : len(s.Cargo[source-1].Crates)]
	// Remove all crates amount from source
	s.Cargo[source-1].Crates = s.Cargo[source-1].Crates[:len(s.Cargo[source-1].Crates)-amount]

	// Add slice to: destination
	s.Cargo[destination-1].Crates = append(s.Cargo[destination-1].Crates, crateToMove...)
}

func (s *GiantCargo) printTopStack() {
	topStack := []string{}
	for _, stack := range s.Cargo {
		topStack = append(topStack, stack.Crates[len(stack.Crates)-1])
	}
	log.Printf("Top Stack is: %s", strings.Join(topStack, ""))
}

func parseInitialState(scanner *bufio.Scanner) GiantCargo {
	rawCargo := []string{}
	nbOfStack := 0
	biggestStack := 0

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		if string(line[0]) == "[" {
			rawCargo = append(rawCargo, line)
			biggestStack++
		} else if string(line[1]) == "1" {
			lastStack := string(line[len(line)-2])
			nbOfStack, _ = strconv.Atoi(lastStack)
		}
	}

	newCargo := GiantCargo{}
	newCargo.init(nbOfStack)

	for i := 1; i <= nbOfStack; i++ {
		stackPos := 1
		if i > 1 {
			// pos: 1 , 5 , 9 ...
			// 1 = 1 => 0
			// 2 = 5 => 3
			// 3 = 9 => 6
			// 4 = 13 => 9
			stackPos = ((i - 1) * 4) + 1
		}
		for j := biggestStack - 1; j >= 0; j-- {
			cargoType := string(rawCargo[j][stackPos])
			if cargoType != " " {
				newCargo.Cargo[i-1].Crates = append(newCargo.Cargo[i-1].Crates, cargoType)
			}
		}
	}

	return newCargo
}

func main() {

	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	giantCargo := parseInitialState(scanner)

	// i.e.: move 3 from 2 to 1
	re := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		amountToMove, _ := strconv.Atoi(match[1])
		fromStack, _ := strconv.Atoi(match[2])
		toStack, _ := strconv.Atoi(match[3])
		giantCargo.move(amountToMove, fromStack, toStack)
	}

	log.Printf("Update Cargo: %v\n", giantCargo)
	giantCargo.printTopStack()
}
