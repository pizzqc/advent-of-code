package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Command int64

const (
	NOOP Command = iota
	ADDX
)

type Pos struct {
	start, end int
}

type CPU struct {
	X            int
	Cycle        int
	Stack        []Instruction
	SignalChecks map[int]int
	CRT          []string
	SpritePos    Pos
}

type Instruction struct {
	cmd   Command
	value int
}

func NewCPU() CPU {
	return CPU{
		X:     1,
		Cycle: -1,
		Stack: []Instruction{},
		SignalChecks: map[int]int{
			20:  0,
			60:  0,
			100: 0,
			140: 0,
			180: 0,
			220: 0,
		},
		CRT:       make([]string, 6),
		SpritePos: Pos{start: 0, end: 2},
	}
}

func (c *CPU) addInstruction(line string) {
	instructionToken := strings.Split(line, " ")
	if instructionToken[0] == "noop" {
		c.Stack = append(c.Stack, Instruction{cmd: NOOP, value: 0})
	} else {
		v, _ := strconv.Atoi(instructionToken[1])
		c.Stack = append(c.Stack, Instruction{cmd: ADDX, value: v})
	}
}

func (c *CPU) computeSignalStr() {
	_, found := c.SignalChecks[c.Cycle]
	if found {
		ss := c.Cycle * c.X
		c.SignalChecks[c.Cycle] = ss
	}
}

func (c *CPU) drawCRT() {
	crtVerticalIdx := int(math.Trunc(float64(c.Cycle) / float64(40)))
	nextDrawPos := len(c.CRT[crtVerticalIdx])
	if nextDrawPos <= c.SpritePos.end && nextDrawPos >= c.SpritePos.start {
		c.CRT[crtVerticalIdx] = c.CRT[crtVerticalIdx] + "#"
	} else {
		c.CRT[crtVerticalIdx] = c.CRT[crtVerticalIdx] + "."
	}
}

func (c *CPU) increaseCycle() {
	c.Cycle++
	c.drawCRT()
}

func (c *CPU) printCRT() {
	for _, crtLine := range c.CRT {
		fmt.Printf("%s\n", crtLine)
	}
}

func (c *CPU) processStack() {
	for _, inst := range c.Stack {
		fmt.Printf("Processing instruction: %v\n", inst)
		switch inst.cmd {
		case NOOP:
			c.increaseCycle()
			c.computeSignalStr()
		case ADDX:
			c.increaseCycle()
			c.computeSignalStr()
			c.increaseCycle()
			c.X += inst.value
			c.computeSignalStr()
			c.SpritePos.start = c.X - 1
			c.SpritePos.end = c.X + 1
		}
		c.printCRT()
	}
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	cpu := NewCPU()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		cpu.addInstruction(line)
	}

	cpu.processStack()

	//  signal strength (the cycle number multiplied by the value of the X register)
	//    during the 20th cycle and every 40 cycles after that (that is, during the 20th, 60th, 100th, 140th, 180th, and 220th cycles).
	signalSum := 0
	for _, val := range cpu.SignalChecks {
		signalSum += val
	}
	fmt.Printf("The sum of these signal strengths is: %v", signalSum)
}
