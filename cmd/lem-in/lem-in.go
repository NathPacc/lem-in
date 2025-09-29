package main

import (
	"fmt"
	"lem-in/colony"
	"lem-in/datas"
	"os"
	"strings"
	"time"
)

// main parses arguments, loads data, checks for errors, builds the colony, and prints the solution.
func main() {
	// Check for correct usage
	start := time.Now()
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Println("Error : Usage is './lem-in filename' or './lem-in filename | ./visualizer")
		return
	}
	filename := os.Args[1]
	// Read and parse the input file
	rawdatas := datas.GetDatas(filename)
	instructions := strings.Join(rawdatas, "\n")
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
	durationColony := time.Since(start)
	startAlgo := time.Now()
	bestset := colony.Resolve(filedatas.NbAnts, rooms)
	fmt.Println(instructions + "\n")
	colony.PrintResolve(filedatas.NbAnts, bestset)
	durationAll := time.Since(start)
	durationAlgo := time.Since(startAlgo)
	fmt.Println("--------------------")
	fmt.Printf("Colony constructed in %s\n", durationColony)
	fmt.Printf("Algo resolution done in %s\n", durationAlgo)
	fmt.Printf("Total resolution done in %s\n", durationAll)
}
