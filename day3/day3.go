package main

import (
	"bufio"
	"log"
	"os"
	"unicode"
)

func main() {

	intputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer intputFile.Close()

	scanner := bufio.NewScanner(intputFile)

	commonItems := []int{}

	groupSackCount := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(groupSackCount) < 3 {
			groupSackCount = append(groupSackCount, line)
			if len(groupSackCount) == 3 {
				commonFound := false
				// log.Printf("%s   /   %s", h1, h2)
				for _, ch := range groupSackCount[0] {
					for _, ch2 := range groupSackCount[1] {
						if ch == ch2 {
							for _, ch3 := range groupSackCount[2] {
								if ch == ch3 {
									// rune a = 97, b = 98 , etc
									// runa A = 65, etc
									priority := 0
									if unicode.IsUpper(ch) {
										priority = int(ch) - 64 + 26
									} else {
										priority = int(ch) - 96
									}

									log.Printf("%s = %v => prio = %v", string(ch), ch, priority)
									commonItems = append(commonItems, priority)
									commonFound = true
									groupSackCount = []string{}
									break
								}
							}
							if commonFound {
								break
							}
						}
					}
					if commonFound {
						break
					}
				}
			}
		}
	}

	totalPrio := 0
	for _, prio := range commonItems {
		totalPrio += prio
	}
	log.Printf("Total bag: %v, total sum: %v", len(commonItems), totalPrio)
}
