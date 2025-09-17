package main

import (
	"fmt"
	"lem-in/colony"
	"lem-in/datas"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error : Usage is './lem-in filename'")
		return
	}
	filename := os.Args[1]
	rawdatas := datas.GetDatas(filename)
	filedatas := datas.SaveDatas(rawdatas)
	datas.CheckErrors(&filedatas)
	if len(filedatas.Errors) != 0 {
		for _, err := range filedatas.Errors {
			fmt.Println(fmt.Errorf("error : %w", err))
		}
		return
	}
	rooms := colony.CreatRooms(filedatas)
	colony.CreatColony(filedatas, rooms)
	bestset := colony.Resolve(filedatas.NbAnts, rooms)
	colony.PrintResolve(filedatas.NbAnts, bestset)
}
