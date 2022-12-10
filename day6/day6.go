package main

import (
	"bufio"
	"log"
	"os"
)

type Marker struct {
	marker []string
}

func (m *Marker) contains(newChar string) (int, bool) {
	for idx, markerChar := range m.marker {
		if markerChar == newChar {
			return idx, true
		}
	}
	return -1, false
}

func (m *Marker) add(newChar string) bool {
	m.marker = append(m.marker, newChar)
	if len(m.marker) == 4 {
		return true
	} else {
		return false
	}
}

func (m *Marker) shiftMarker(idx int) {
	m.marker = m.marker[idx+1:]
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		sop := Marker{
			marker: []string{},
		}

		for idx, ch := range line {
			currentChar := string(ch)
			markerIdx, found := sop.contains(currentChar)
			if found {
				sop.shiftMarker(markerIdx)
			}

			markerResolved := sop.add(currentChar)
			if markerResolved {
				log.Printf("First marker after character: %v", idx+1)
				break
			}

		}
	}
}
