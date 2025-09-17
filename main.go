// main.go is the entry point for the lem-in program.
package main

import (
	"fmt"
	"lem-in/colony"
	"lem-in/datas"
	"os"
)

// main parses arguments, loads data, checks for errors, builds the colony, and prints the solution.
func main() {
	// Check for correct usage
	if len(os.Args) != 2 {
		fmt.Println("Error : Usage is './lem-in filename'")
		return
	}
	filename := os.Args[1]
	// Read and parse the input file
	rawdatas := datas.GetDatas(filename)
	filedatas := datas.SaveDatas(rawdatas)
	// Run error checks
	datas.CheckErrors(&filedatas)
	if len(filedatas.Errors) != 0 {
		for _, err := range filedatas.Errors {
			fmt.Println(fmt.Errorf("error : %w", err))
		}
		return
	}
	// Build the rooms and colony structure
	rooms := colony.CreatRooms(filedatas)
	colony.CreatColony(filedatas, rooms)
	// Find and print the best solution
	bestset := colony.Resolve(filedatas.NbAnts, rooms)
	colony.PrintResolve(filedatas.NbAnts, bestset)
}
