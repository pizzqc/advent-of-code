package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ElfDir struct {
	Parent        *ElfDir
	Name          string
	Size          int
	Dirs          map[string]*ElfDir
	Files         map[string]*ElfFile
	inListingMode bool
}

type ElfFile struct {
	Parent *ElfDir
	Name   string
	Size   int
}

func (d *ElfDir) containsDirectory(directory string) bool {
	if _, ok := d.Dirs[directory]; ok {
		return true
	}
	return false
}

func (d *ElfDir) containsFile(filename string) bool {
	if _, ok := d.Files[filename]; ok {
		return true
	}
	return false
}

func (d *ElfDir) initDirectory(directory string) {
	d.Dirs[directory] = &ElfDir{
		Parent: d,
		Name:   directory,
		Size:   0,
		Dirs:   make(map[string]*ElfDir),
		Files:  make(map[string]*ElfFile),
	}
}

func (d *ElfDir) initFile(filesize int, filename string) {
	d.Files[filename] = &ElfFile{
		Parent: d,
		Name:   filename,
		Size:   filesize,
	}
}

func (d *ElfDir) cd(directory string) *ElfDir {
	// We want to go all the way back to the first directory (root).  Root doesn't have a parent
	if directory == "/" {
		for d.Parent != nil {
			d = d.Parent
		}
	} else if directory == ".." {
		d = d.Parent
	} else { // Last cases are going into a specific directory
		// Check if the directory already exist in current folder
		if d.containsDirectory(directory) {
			// If exist change dir to new directory
			d = d.Dirs[directory]
		} else {
			// Changing to a directory never initialized before
			d.initDirectory(directory)
			d = d.Dirs[directory]
		}
	}

	return d // return new current directory
}

func (d *ElfDir) increaseSize(filesize int) {
	// Increase the current directory total size
	d.Size += filesize

	for d.Parent != nil {
		d = d.Parent
		d.Size += filesize
	}
}

func (d *ElfDir) treeView(level int) {
	fmt.Printf("%s- %s (dir, size=%v)\n", strings.Repeat(" ", level*4), d.Name, d.Size)
	for _, subFolder := range d.Dirs {
		subFolder.treeView(level + 1)
	}

	for _, file := range d.Files {
		fmt.Printf("%s- %s (file, size=%v)\n", strings.Repeat(" ", (level+1)*4), file.Name, file.Size)
	}
}

func (d *ElfDir) calculateSum(atMost int) int {
	totalSum := 0
	if d.Size <= atMost {
		totalSum = d.Size
	}

	for _, directory := range d.Dirs {
		totalSum += directory.calculateSum(atMost)
	}

	return totalSum
}

func (d *ElfDir) findSmallestWithAtLeast(neededSpace int) *ElfDir {
	var goodCandidate *ElfDir
	for _, directory := range d.Dirs {
		candidate := directory.findSmallestWithAtLeast(neededSpace)
		if candidate != nil && candidate.Size > neededSpace {
			if goodCandidate == nil {
				goodCandidate = candidate
			} else if candidate.Size < goodCandidate.Size {
				goodCandidate = candidate
			}
		}
	}

	if d.Size > neededSpace {
		if goodCandidate == nil {
			return d
		} else if d.Size < goodCandidate.Size {
			return d
		}
	}
	return goodCandidate
}

func (d *ElfDir) String() string {
	return fmt.Sprintf("Name: %s, Size: %v", d.Name, d.Size)
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	// Initialize the root directory
	currentDirectory := &ElfDir{
		Parent: nil,
		Name:   "/",
		Size:   0,
		Dirs:   make(map[string]*ElfDir),
		Files:  make(map[string]*ElfFile),
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineToken := strings.Split(line, " ")

		// Processing directory listing
		if currentDirectory.inListingMode {
			switch lineToken[0] {
			case "dir":
				if !currentDirectory.containsDirectory(lineToken[1]) {
					currentDirectory.initDirectory(lineToken[1])
				}
				continue
			case "$": // this line is a new command
				currentDirectory.inListingMode = false
			default: // if not a DIR it is a file
				if !currentDirectory.containsFile(lineToken[1]) {
					filesize, _ := strconv.Atoi(lineToken[0])
					currentDirectory.initFile(filesize, lineToken[1])
					currentDirectory.increaseSize(filesize)
				}
				continue
			}
		}

		if lineToken[0] == "$" {
			// Process a command
			switch lineToken[1] {
			case "cd":
				currentDirectory = currentDirectory.cd(lineToken[2])
			case "ls":
				currentDirectory.inListingMode = true
			}
		}
	}

	// Set the current folder to root
	currentDirectory = currentDirectory.cd("/")

	// Debug printout of the treeView structure of dir/files
	currentDirectory.treeView(0)

	// Part #1:
	//   Find all of the directories with a total size of at most 100000, then calculate the sum of their total sizes.
	part1Answer := currentDirectory.calculateSum(100000)
	fmt.Printf("The sum is: %v\n", part1Answer)

	// Part #2:
	totalDisk := 70000000
	needSpace := 30000000

	currentFreeSpace := totalDisk - currentDirectory.Size

	bestCandidate := currentDirectory.findSmallestWithAtLeast(needSpace - currentFreeSpace)
	fmt.Printf("The best candidate for deletion is: %v\n", bestCandidate)
}
