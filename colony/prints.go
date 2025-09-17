package colony

import (
	"fmt"
	"sort"
	"strings"
)

func PrintRoom(room Room) {
	var ids []string
	for _, neighbour := range room.Neighbours {
		ids = append(ids, neighbour.Name)
	}
	neighboursStr := strings.Join(ids, ", ")
	fmt.Println("Room", room.Name, "-> Neighbours:", neighboursStr)
}

func PrintPath(path []*Room) {
	for i, room := range path {
		if i == len(path)-1 {
			fmt.Print(room.Name + "\n")
		} else {
			fmt.Print(room.Name + "->")
		}
	}
}

func PrintColony(roomlist []*Room) {
	heigh, width := calculateSize(roomlist)
	for line := 0; line <= heigh; line++ {
		var strline string
		for coloumn := 0; coloumn <= width; coloumn++ {
			needPlacement := false
			for _, room := range roomlist {
				if room.Coordinates.Y == line && room.Coordinates.X == coloumn {
					strline += string(room.Name[0])
					needPlacement = true
				}
			}
			if !needPlacement {
				strline += "."
			}
		}
		fmt.Println(strline)
	}
}

func calculateSize(roomlist []*Room) (height, width int) {
	maxheigh := 0
	maxwidth := 0
	for _, room := range roomlist {
		if room.Coordinates.Y > maxheigh {
			maxheigh = room.Coordinates.Y
		}
		if room.Coordinates.X > maxwidth {
			maxwidth = room.Coordinates.X
		}
	}
	return maxheigh, maxwidth
}

func PrintResolve(nbAnt int, paths [][]*Room) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	_, antsPerPath := calculateTime(nbAnt, paths)

	var antIDs []int
	var antPaths []int
	var antPositions []int
	antsSent := 0
	antsFinished := 0
	pathCursor := make([]int, len(paths)) // combien de fourmis déjà envoyées sur chaque chemin

	for tour := 1; antsFinished < nbAnt; tour++ {

		// Avancer les fourmis déjà envoyées
		for i := 0; i < len(antIDs); i++ {
			if antPositions[i] < len(paths[antPaths[i]])-1 {
				antPositions[i]++
				room := paths[antPaths[i]][antPositions[i]].Name
				fmt.Printf("L%d-%s ", antIDs[i], room)
				if antPositions[i] == len(paths[antPaths[i]])-1 {
					antsFinished++
				}
			}
		}

		// Envoyer les fourmis selon le plan
		for i := range paths {
			if pathCursor[i] < antsPerPath[i] {
				antsSent++
				pathCursor[i]++
				antIDs = append(antIDs, antsSent)
				antPaths = append(antPaths, i)
				antPositions = append(antPositions, 0)
			}
		}
		if tour > 1 {
			fmt.Println("")
		}
	}
}
