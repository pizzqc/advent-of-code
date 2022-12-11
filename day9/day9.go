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

type Matrix struct {
	StartPos []int
	HeadPos  []int
	TailPos  []int
	Cells    [][]Cell
	Debug    bool
}

type Cell struct {
	x, y         int
	visitedCount int
}

type CellOptions struct {
	X, Y int
}

func NewCell(options CellOptions) *Cell {
	return &Cell{
		x:            options.X,
		y:            options.Y,
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

func (m *Matrix) isHeadTailTouch() bool {
	if m.TailPos[1] == m.HeadPos[1]-1 ||
		m.TailPos[1] == m.HeadPos[1] ||
		m.TailPos[1] == m.HeadPos[1]+1 {		
		if m.TailPos[0] == m.HeadPos[0]-1 ||
			m.TailPos[0] == m.HeadPos[0] ||
			m.TailPos[0] == m.HeadPos[0]+1 {
			return true
		}
	}

	return false
}

func (m *Matrix) markVisited() {
	m.Cells[m.TailPos[0]][m.TailPos[1]].visitedCount++
}

func (m *Matrix) moveHead(direction common.Direction) {
	switch direction {
	case common.RIGHT:
		m.HeadPos[1]++
	case common.UP:
		m.HeadPos[0]--
	case common.LEFT:
		m.HeadPos[1]--
	case common.DOWN:
		m.HeadPos[0]++
	}
}

func (m *Matrix) moveTail(direction common.Direction) {
	switch direction {
	case common.RIGHT:
		m.TailPos[1]++
		if m.TailPos[0] < m.HeadPos[0] {
			// Go Diag DownRight
			m.TailPos[0]++
		} else if m.TailPos[0] > m.HeadPos[0] {
			// Go Diag UpRight
			m.TailPos[0]--
		}

	case common.UP:
		m.TailPos[0]--
		if m.TailPos[1] < m.HeadPos[1] {
			// Go Diag UpRight
			m.TailPos[1]++
		} else if m.TailPos[1] > m.HeadPos[1] {
			// Go Diag UpLeft
			m.TailPos[1]--
		}

	case common.LEFT:
		m.TailPos[1]--
		if m.TailPos[0] < m.HeadPos[0] {
			// Go Diag DownLeft
			m.TailPos[0]++
		} else if m.TailPos[0] > m.HeadPos[0] {
			// Go Diag UpLeft
			m.TailPos[0]--
		}
	case common.DOWN:
		m.TailPos[0]++
		if m.TailPos[1] < m.HeadPos[1] {
			// Go Diag UpRight
			m.TailPos[1]++
		} else if m.TailPos[1] > m.HeadPos[1] {
			// Go Diag UpLeft
			m.TailPos[1]--
		}
	}
	m.markVisited()
}

func (m *Matrix) Right(count int) {
	fmt.Printf("\n== R %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the right (i.e.: ..H)
		if m.HeadPos[1] == len(m.Cells[m.HeadPos[0]])-1 {
			m.increaseMatrix(common.RIGHT)
			m.Print()
		}

		m.Move(common.RIGHT)
	}
}

func (m *Matrix) Up(count int) {
	fmt.Printf("\n== U %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the top
		if m.HeadPos[0] == 0 {
			m.increaseMatrix(common.UP)
			m.Print()
		}

		m.Move(common.UP)
	}
}

func (m *Matrix) Down(count int) {
	fmt.Printf("\n== D %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the bottom
		if m.HeadPos[0] == len(m.Cells)-1 {
			m.increaseMatrix(common.DOWN)
			m.Print()
		}

		m.Move(common.DOWN)
	}
}

func (m *Matrix) Left(count int) {
	fmt.Printf("\n== L %v ==\n", count)
	for i := 1; i <= count; i++ {
		// Check if we need to increase matrix size to the Left
		if m.HeadPos[1] == 0 {
			m.increaseMatrix(common.LEFT)
			m.Print()
		}

		m.Move(common.LEFT)

	}
}

func (m *Matrix) Move(direction common.Direction) {
	// Move: HEAD
	m.moveHead(direction)
	m.Print()

	// Move: TAIL
	if !m.isHeadTailTouch() {
		m.moveTail(direction)
		m.Print()
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
		nr := NewRow(len(m.Cells), len(m.Cells[m.HeadPos[0]]))
		m.Cells = append([][]Cell{nr}, m.Cells...)
		// Shift Start, Head and Tail X-axis by 1 index down
		m.StartPos[0]++
		m.HeadPos[0]++
		m.TailPos[0]++
	case common.LEFT:
		for i, row := range m.Cells {
			nc := NewCell(CellOptions{
				X: i,
				Y: 0})
			m.Cells[i] = append([]Cell{*nc}, row...)
		}
		m.StartPos[1]++
		m.HeadPos[1]++
		m.TailPos[1]++
	case common.DOWN:
		nr := NewRow(len(m.Cells), len(m.Cells[m.HeadPos[0]]))
		m.Cells = append(m.Cells, nr)
	}
}

func (m *Matrix) Print() {
	if m.Debug {
		fmt.Println("")
		for i, row := range m.Cells {
			var line strings.Builder
			for j, col := range row {
				if m.isHead(i, j) && m.isTail(i, j) && m.isStart(i, j) {
					line.WriteString("(HTs)")
				} else if m.isHead(i, j) && m.isTail(i, j) && !m.isStart(i, j) {
					line.WriteString("(HT)")
				} else if m.isHead(i, j) && !m.isTail(i, j) && m.isStart(i, j) {
					line.WriteString("(Hs)")
				} else if !m.isHead(i, j) && m.isTail(i, j) && m.isStart(i, j) {
					line.WriteString("(Ts)")
				} else if m.isHead(i, j) {
					line.WriteString("H")
				} else if m.isTail(i, j) {
					line.WriteString("T")
				} else if m.isStart(i, j) {
					line.WriteString("s")
				} else if col.visitedCount > 0 {
					line.WriteString("#")
				} else {
					line.WriteString(".")
				}
			}
			fmt.Println(line.String())
		}
		fmt.Println("")
	}
}

func (m *Matrix) isHead(i, j int) bool {
	return m.HeadPos[0] == i && m.HeadPos[1] == j
}

func (m *Matrix) isTail(i, j int) bool {
	return m.TailPos[0] == i && m.TailPos[1] == j
}

func (m *Matrix) isStart(i, j int) bool {
	return m.StartPos[0] == i && m.StartPos[1] == j
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	m := &Matrix{
		StartPos: []int{0, 0},
		HeadPos:  []int{0, 0},
		TailPos:  []int{0, 0},
		Cells: [][]Cell{
			{
				Cell{
					x:            0,
					y:            0,
					visitedCount: 1,
				},
			},
		},
		Debug: false,
	}

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
	}

	// Part1 : How many positions does the tail of the rope visit at least once?
	part1Result := 0
	for _, row := range m.Cells {
		for _, cell := range row {
			if cell.visitedCount >= 1 {
				part1Result++
			}
		}
	}
	m.Debug = true
	m.Print()
	fmt.Printf("Part1 Answer: %v", part1Result)
}
