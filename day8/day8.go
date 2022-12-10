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

func isVisible(rowIdx, colIdx, treeSize int, forest [][]int, direction Direction) bool {
	treeVisible := false
	if forest[rowIdx][colIdx] < treeSize {
		treeVisible = true
	} else if forest[rowIdx][colIdx] >= treeSize {
		return false
	}

	// Check if position is valid to keep traversing further
	if rowIdx >= 0 || colIdx >= 0 || rowIdx < len(forest)-1 || colIdx < len(forest[rowIdx])-1 {
		switch direction {
		case UP:
			if rowIdx > 0 {
				// Recurse UP all the way or until we know it is hidden
				return isVisible(rowIdx-1, colIdx, treeSize, forest, UP)
			}
		case DOWN:
			if rowIdx < len(forest)-1 {
				// Recurse DOWN all the way or until we know it is hidden
				return isVisible(rowIdx+1, colIdx, treeSize, forest, DOWN)
			}
		case LEFT:
			if colIdx > 0 {
				// Recurse LEFT all the way or until we know it is hidden
				return isVisible(rowIdx, colIdx-1, treeSize, forest, LEFT)
			}
		case RIGHT:
			if colIdx < len(forest[rowIdx])-1 {
				// Recurse RIGHT all the way or until we know it is hidden
				return isVisible(rowIdx, colIdx+1, treeSize, forest, RIGHT)
			}
		}
	}

	return treeVisible
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

	// Traverse the 2D arrays and build a mask of all the visible ones
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
				visibleTrees[i][j] = 1
			} else if isVisible(i-1, j, col, forest, UP) {
				visibleTrees[i][j] = 1
			} else if isVisible(i+1, j, col, forest, DOWN) {
				visibleTrees[i][j] = 1
			} else if isVisible(i, j+1, col, forest, RIGHT) {
				visibleTrees[i][j] = 1
			} else if isVisible(i, j-1, col, forest, LEFT) {
				visibleTrees[i][j] = 1
			}
		}
	}

	// Part 1: Total visible trees
	totalVisible := 0
	for _, row := range visibleTrees {
		for _, col := range row {
			if col == 1 {
				totalVisible++
			}
		}
	}
	fmt.Printf("Visible trees: %v", totalVisible)
}
