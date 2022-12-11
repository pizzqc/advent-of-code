package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Direction int

const (
	UP Direction = iota + 1
	DOWN
	LEFT
	RIGHT
)

func isVisible(rowIdx, colIdx, treeSize int, forest [][]int, direction Direction) (bool, int) {
	treeVisible := false
	nbTreeVisible := 0
	if forest[rowIdx][colIdx] < treeSize {
		treeVisible = true
		nbTreeVisible++
	} else if forest[rowIdx][colIdx] >= treeSize {
		nbTreeVisible++
		return false, nbTreeVisible
	}

	// Check if position is valid to keep traversing further
	if rowIdx >= 0 || colIdx >= 0 || rowIdx < len(forest)-1 || colIdx < len(forest[rowIdx])-1 {
		switch direction {
		case UP:
			if rowIdx > 0 {
				// Recurse UP all the way or until we know it is hidden
				visible, count := isVisible(rowIdx-1, colIdx, treeSize, forest, UP)
				nbTreeVisible += count
				return visible, nbTreeVisible
			}
		case DOWN:
			if rowIdx < len(forest)-1 {
				// Recurse DOWN all the way or until we know it is hidden
				visible, count := isVisible(rowIdx+1, colIdx, treeSize, forest, DOWN)
				nbTreeVisible += count
				return visible, nbTreeVisible
			}
		case LEFT:
			if colIdx > 0 {
				// Recurse LEFT all the way or until we know it is hidden
				visible, count := isVisible(rowIdx, colIdx-1, treeSize, forest, LEFT)
				nbTreeVisible += count
				return visible, nbTreeVisible
			}
		case RIGHT:
			if colIdx < len(forest[rowIdx])-1 {
				// Recurse RIGHT all the way or until we know it is hidden
				visible, count := isVisible(rowIdx, colIdx+1, treeSize, forest, RIGHT)
				nbTreeVisible += count
				return visible, nbTreeVisible
			}
		}
	}

	return treeVisible, nbTreeVisible
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	// Build 2D array of the entire forest
	forest := make([][]int, 0)
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		forest = append(forest, make([]int, len(line)))
		for col := 0; col < len(line); col++ {
			treeSize, _ := strconv.Atoi(string(line[col]))
			forest[row][col] = treeSize
		}
		row++
	}

	// Initialize the 2D visbleTrees arrays tracker
	visibleTrees := make([][]int, len(forest))
	for i := 0; i < len(visibleTrees); i++ {
		visibleTrees[i] = make([]int, len(forest[i]))
		for j := 0; j < len(visibleTrees[i]); j++ {
			visibleTrees[i][j] = 0 // 0 = not-visible ; 1 = visible
		}
	}

	// Traverse the forest
	for i, row := range forest {
		for j, col := range row {
			if i == 0 || j == 0 || i == len(forest)-1 || j == len(forest[i])-1 {
				// Anything on the edge is visible
				visibleTrees[i][j] = 0
			} else {
				_, UpCount := isVisible(i-1, j, col, forest, UP)
				_, DownCount := isVisible(i+1, j, col, forest, DOWN)
				_, RightCount := isVisible(i, j+1, col, forest, RIGHT)
				_, LeftCount := isVisible(i, j-1, col, forest, LEFT)
				visibleTrees[i][j] = (UpCount * DownCount * RightCount * LeftCount)
			}
		}
	}

	// Part 2: Scenic Score
	maxScenicScore := 0
	for _, row := range visibleTrees {
		for _, col := range row {
			if col > maxScenicScore {
				maxScenicScore = col
			}
		}
	}
	fmt.Printf("Max Scenic Score: %v\n", maxScenicScore)

}
