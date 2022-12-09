package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type elf struct {
	index    int
	calories int
}

func main() {
	intputFile, _ := os.Open("input.txt")
	defer intputFile.Close()

	scanner := bufio.NewScanner(intputFile)

	listOfElf := []elf{}
	elfTotal := 0
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			listOfElf = append(listOfElf, elf{
				index:    idx,
				calories: elfTotal,
			})
			idx++
			elfTotal = 0
		} else {
			calories, _ := strconv.Atoi(line)
			elfTotal = elfTotal + calories
		}
	}

	sort.Slice(listOfElf, func(i, j int) bool {
		return listOfElf[i].calories > listOfElf[j].calories
	})
	grandTotal := listOfElf[0].calories + listOfElf[1].calories + listOfElf[2].calories
	fmt.Println(grandTotal)
}
