package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	START_MARKER string = "S"
	END_MARKER   string = "E"
)

type ElevationMap struct {
	Elevation [][]int
	Start     Position
	End       Position
}

type Position struct {
	X int
	Y int
}

type Traveler struct {
	Elevation       ElevationMap
	DistanceTracker [][]int
	Start           Position
	Queue           *list.List
}

func NewTraveler(emap ElevationMap, start Position) *Traveler {
	d := make([][]int, len(emap.Elevation))
	for i, row := range emap.Elevation {
		d[i] = make([]int, len(row))
		for j, _ := range row {
			d[i][j] = -1
		}
	}

	t := &Traveler{
		Elevation:       emap,
		Start:           start,
		Queue:           list.New(),
		DistanceTracker: d,
	}

	t.Queue.PushBack(start)

	return t
}

func NewPos(x, y int) Position {
	return Position{X: x, Y: y}
}

func (p Position) GetNeighbors() []Position {
	up := Position{p.X - 1, p.Y}
	right := Position{p.X, p.Y + 1}
	down := Position{p.X + 1, p.Y}
	left := Position{p.X, p.Y - 1}

	return []Position{up, right, down, left}
}

func (t *Traveler) travel() error {

	// If nothing to travel to we are done
	if t.Queue.Len() == 0 {
		return errors.New("nothing left to traverse")
	}

	position := t.Queue.Front()
	t.Queue.Remove(position)
	pos := position.Value.(Position)

	// If start set to 0 the tracker
	if pos == t.Start {
		t.DistanceTracker[pos.X][pos.Y] = 0
	}

	neighbors := t.Elevation.FindNeighbors(pos.X, pos.Y)

	for _, n := range neighbors {
		if t.DistanceTracker[n.X][n.Y] == -1 {
			t.DistanceTracker[n.X][n.Y] = t.DistanceTracker[pos.X][pos.Y] + 1
			t.Queue.PushBack(n)
		}
	}

	return nil
}

func (e *ElevationMap) FindNeighbors(row, col int) []Position {
	validOptions := []Position{}

	currentElevation := int(e.Elevation[row][col])

	// Check UP is valid
	if row > 0 {
		if int(e.Elevation[row-1][col]) <= currentElevation+1 {
			validOptions = append(validOptions, Position{X: row - 1, Y: col})
		}
	}

	// Check RIGHT is valid
	if len(e.Elevation[row])-1 > col {
		if int(e.Elevation[row][col+1]) <= currentElevation+1 {
			validOptions = append(validOptions, Position{X: row, Y: col + 1})
		}
	}

	// Check DOWN is valid
	if len(e.Elevation)-1 > row {
		if int(e.Elevation[row+1][col]) <= currentElevation+1 {
			validOptions = append(validOptions, Position{X: row + 1, Y: col})
		}
	}

	// Check LEFT is valid
	if col > 0 {
		if int(e.Elevation[row][col-1]) <= currentElevation+1 {
			validOptions = append(validOptions, Position{X: row, Y: col - 1})
		}
	}

	return validOptions
}

func (t *Traveler) DistanceFrom(pos Position) int {
	return t.DistanceTracker[pos.X][pos.Y]
}

func GetShortestPath(emap ElevationMap, start Position) int {
	t := NewTraveler(emap, start)
	for t.Queue.Len() > 0 {
		t.travel()
	}

	return t.DistanceFrom(emap.End)
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	emap := ElevationMap{
		Elevation: [][]int{},
	}

	scanner := bufio.NewScanner(inputFile)
	lineNb := 0
	for scanner.Scan() {
		line := scanner.Text()
		emap.Elevation = append(emap.Elevation, make([]int, len(line)))
		for i, ch := range line {
			emap.Elevation[lineNb][i] = int(ch)
			if string(ch) == START_MARKER {
				emap.Elevation[lineNb][i] = int(rune('a'))
				emap.Start = NewPos(lineNb, i)
			} else if string(ch) == END_MARKER {
				emap.Elevation[lineNb][i] = int(rune('z'))
				emap.End = NewPos(lineNb, i)
			}
		}
		lineNb++
	}

	// Results
	allLowestPointSteps := make(map[Position]int)

	for i, row := range emap.Elevation {
		for j, col := range row {
			if col == int(rune('a')) {
				pos := Position{X: i, Y: j}
				count := GetShortestPath(emap, pos)
				allLowestPointSteps[pos] = count
			}
		}
	}

	var lowestPos Position
	var lowestStep int
	for key, count := range allLowestPointSteps {
		if count != -1 {
			if lowestStep == 0 {
				lowestStep = count
				lowestPos = key
			} else if count < lowestStep {
				lowestStep = count
				lowestPos = key
			}
		}
	}
	fmt.Printf("What is the fewest steps required to move from your current position to the location that should get the best signal?: %v from position %v\n", lowestStep, lowestPos)
}
