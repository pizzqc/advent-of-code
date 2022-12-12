package main

import (
	"advent-code-2022/common"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	RopeSize int = 10
)

type CartesianPosition struct {
	X, Y int
}

type RopeKnot struct {
	Position CartesianPosition
	Number   int
	Parent   *RopeKnot
	Next     *RopeKnot
}

type Matrix struct {
	StartPos CartesianPosition
	Head     *RopeKnot
	Cells    [][]Cell
	Debug    bool
}

type Cell struct {
	Position     CartesianPosition
	visitedCount int
}

type CellOptions struct {
	X, Y int
}

func NewRope(size int) *RopeKnot {
	head := &RopeKnot{
		Position: CartesianPosition{X: 0, Y: 0},
		Number:   0,
		Parent:   nil,
		Next:     nil,
	}

	currentKnot := head
	for i := 1; i < size; i++ {
		currentKnot.Next = &RopeKnot{
			Position: CartesianPosition{X: 0, Y: 0},
			Number:   i,
			Parent:   currentKnot,
			Next:     nil,
		}
		currentKnot = currentKnot.Next
	}

	return head

}

func NewCell(options CellOptions) *Cell {
	return &Cell{
		Position:     CartesianPosition{X: options.X, Y: options.Y},
		visitedCount: 0,
	}
}

func NewRow(x, length int) []Cell {
	newRow := []Cell{}
	for y := 0; y < length; y++ {
		newRow = append(newRow, *NewCell(CellOptions{
			X: x,
			Y: y,
		}))
	}

	return newRow
}

func (c *Cell) String() string {
	return fmt.Sprint(c.visitedCount)
}

func (r *RopeKnot) isConnected() bool {
	if r.Next.Position.Y == r.Position.Y-1 ||
		r.Next.Position.Y == r.Position.Y ||
		r.Next.Position.Y == r.Position.Y+1 {
		if r.Next.Position.X == r.Position.X-1 ||
			r.Next.Position.X == r.Position.X ||
			r.Next.Position.X == r.Position.X+1 {
			return true
		}
	}
	return false
}

func (r *RopeKnot) shift(direction common.Direction) {
	currentKnot := r

	switch direction {
	case common.UP:
		currentKnot.Position.X++
		for currentKnot.Next != nil {
			currentKnot = currentKnot.Next
			currentKnot.Position.X++
		}
	case common.LEFT:
		currentKnot.Position.Y++
		for currentKnot.Next != nil {
			currentKnot = currentKnot.Next
			currentKnot.Position.Y++
		}
	}
}

func (m *Matrix) markVisited(pos CartesianPosition) {
	m.Cells[pos.X][pos.Y].visitedCount++
}

func (m *Matrix) moveHead(direction common.Direction) {
	switch direction {
	case common.RIGHT:
		m.Head.Position.Y++
	case common.UP:
		m.Head.Position.X--
	case common.LEFT:
		m.Head.Position.Y--
	case common.DOWN:
		m.Head.Position.X++
	}
}

func (r *RopeKnot) moveKnot() {
	// Using diff to figure out how to get closer to Parent
	diffX := r.Position.X - r.Parent.Position.X
	diffY := r.Position.Y - r.Parent.Position.Y

	if diffX == 0 { // SAME ROW
		if diffY < 0 {
			r.Position.Y++
		} else {
			r.Position.Y--
		}
	} else if diffX > 0 { // 1+ ROW BELOW
		if diffY < 0 {
			r.Position.X--
			r.Position.Y++
		} else if diffY == 0 {
			r.Position.X--
		} else if diffY > 0 {
			r.Position.X--
			r.Position.Y--
		}
	} else if diffX < 0 { // 1+ ROW BELOW
		if diffY < 0 {
			r.Position.X++
			r.Position.Y++
		} else if diffY == 0 {
			r.Position.X++
		} else if diffY > 0 {
			r.Position.X++
			r.Position.Y--
		}
	}
}

func (m *Matrix) Right(count int) {
	fmt.Printf("\n== R %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the right (i.e.: ..H)
		if m.Head.Position.Y == len(m.Cells[m.Head.Position.X])-1 {
			m.increaseMatrix(common.RIGHT)
		}

		m.Move(common.RIGHT)
	}
}

func (m *Matrix) Up(count int) {
	fmt.Printf("\n== U %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the top
		if m.Head.Position.X == 0 {
			m.increaseMatrix(common.UP)
		}

		m.Move(common.UP)
	}
}

func (m *Matrix) Down(count int) {
	fmt.Printf("\n== D %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the bottom
		if m.Head.Position.X == len(m.Cells)-1 {
			m.increaseMatrix(common.DOWN)
		}

		m.Move(common.DOWN)
	}
}

func (m *Matrix) Left(count int) {
	fmt.Printf("\n== L %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the Left
		if m.Head.Position.Y == 0 {
			m.increaseMatrix(common.LEFT)
		}

		m.Move(common.LEFT)

	}
}

func (m *Matrix) Move(direction common.Direction) {
	// Move: HEAD
	m.moveHead(direction)

	// Move: Tail Knots
	currentKnot := m.Head
	for currentKnot.Next != nil {
		if !currentKnot.isConnected() {
			currentKnot.Next.moveKnot()

			// Check if we moved the Tail of the Rope
			if currentKnot.Next.Next == nil {
				m.markVisited(currentKnot.Next.Position)
			}
		}
		currentKnot = currentKnot.Next
	}
}

func (m *Matrix) increaseMatrix(direction common.Direction) {
	switch direction {
	case common.RIGHT:
		for i, row := range m.Cells {
			nc := NewCell(CellOptions{
				X: i,
				Y: len(row)})
			m.Cells[i] = append(row, *nc)
		}
	case common.UP:
		nr := NewRow(len(m.Cells), len(m.Cells[m.Head.Position.X]))
		m.Cells = append([][]Cell{nr}, m.Cells...)
		m.StartPos.X++
		m.Head.shift(common.UP)

	case common.LEFT:
		for i, row := range m.Cells {
			nc := NewCell(CellOptions{
				X: i,
				Y: 0})
			m.Cells[i] = append([]Cell{*nc}, row...)
		}
		m.StartPos.Y++
		m.Head.shift(common.LEFT)
	case common.DOWN:
		nr := NewRow(len(m.Cells), len(m.Cells[m.Head.Position.X]))
		m.Cells = append(m.Cells, nr)
	}
}

func (m *Matrix) Print() {
	if m.Debug {
		fmt.Println("")
		for i, row := range m.Cells {
			var line strings.Builder
			overlap := make(map[string][]RopeKnot)
			for j, col := range row {
				knotFound := m.knotsFound(i, j)
				if len(knotFound) > 0 {
					knotIdx := knotFound[0].Number
					knotName := strconv.Itoa(knotIdx)
					if knotIdx == 0 {
						knotName = "H"
					}

					line.WriteString(knotName)
					if len(knotFound) > 1 {
						overlap[knotName] = knotFound[1:]
					}
				} else if m.isStart(i, j) {
					line.WriteString("s")
				} else if col.visitedCount > 0 {
					line.WriteString("#")
				} else {
					line.WriteString(".")
				}
			}
			fmt.Print(line.String())
			// Print covers
			for key, over := range overlap {
				coversList := []string{}
				for _, knot := range over {
					coversList = append(coversList, strconv.Itoa(knot.Number))
				}
				if m.isStart(over[0].Position.X, over[0].Position.Y) {
					coversList = append(coversList, "s")
				}

				fmt.Printf("  (%s covers %v)", key, strings.Join(coversList, ","))
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func (m *Matrix) knotsFound(i, j int) []RopeKnot {
	knotsFound := []RopeKnot{}
	currentKnot := m.Head
	if currentKnot.Position.X == i && currentKnot.Position.Y == j {
		knotsFound = append(knotsFound, *currentKnot)
	}

	for currentKnot.Next != nil {
		currentKnot = currentKnot.Next
		if currentKnot.Position.X == i && currentKnot.Position.Y == j {
			knotsFound = append(knotsFound, *currentKnot)
		}
	}

	return knotsFound
}

func (m *Matrix) isStart(i, j int) bool {
	return m.StartPos.X == i && m.StartPos.Y == j
}

func (r *RopeKnot) size() int {
	size := 1
	for r.Next != nil {
		size++
		r = r.Next
	}

	return size
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	m := &Matrix{
		StartPos: CartesianPosition{X: 0, Y: 0},
		Head:     NewRope(10),
		Cells: [][]Cell{
			{
				Cell{
					Position:     CartesianPosition{X: 0, Y: 0},
					visitedCount: 1,
				},
			},
		},
		Debug: false,
	}
	fmt.Printf("Rope Size: %v\n", m.Head.size())
	m.Print()
	for scanner.Scan() {
		line := scanner.Text()
		direction := strings.Split(line, " ")[0]
		count, _ := strconv.Atoi(strings.Split(line, " ")[1])

		switch direction {
		case "R":
			m.Right(count)
		case "U":
			m.Up(count)
		case "L":
			m.Left(count)
		case "D":
			m.Down(count)
		}
		m.Print()
	}

	// Part1 + Part2 : How many positions does the tail of the rope visit at least once?
	tailVisitedCount := 0
	for _, row := range m.Cells {
		for _, cell := range row {
			if cell.visitedCount >= 1 {
				tailVisitedCount++
			}
		}
	}
	fmt.Printf("Tail Visited at least once: %v", tailVisitedCount)
}
