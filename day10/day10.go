package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Command int64

const (
	NOOP Command = iota
	ADDX
)

type CPU struct {
	X            int
	Cycle        int
	Stack        []Instruction
	SignalChecks map[int]int
}

type Instruction struct {
	cmd   Command
	value int
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

func (c *CPU) processStack() {
	for _, inst := range c.Stack {
		fmt.Printf("Processing instruction: %v\n", inst)
		switch inst.cmd {
		case NOOP:
			c.Cycle++
			c.computeSignalStr()
		case ADDX:
			c.Cycle++
			c.computeSignalStr()
			c.Cycle++
			c.X += inst.value
			c.computeSignalStr()
		}
	}
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	cpu := CPU{
		X:     1,
		Cycle: 0,
		Stack: []Instruction{},
		SignalChecks: map[int]int{
			20:  0,
			60:  0,
			100: 0,
			140: 0,
			180: 0,
			220: 0,
		},
	}

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
