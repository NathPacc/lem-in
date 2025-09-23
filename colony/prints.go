// Package colony provides printing utilities for rooms, paths, and colony visualization.
package colony

import (
	"fmt"
	"lem-in/modules"
	"sort"
	"strings"
)

// Print le nom du salle et la liste de ses voisins.
func PrintRoom(room modules.Room) {
	var ids []string
	for _, neighbour := range room.Neighbours {
		ids = append(ids, neighbour.Name)
	}
	neighboursStr := strings.Join(ids, ", ")
	fmt.Println("Room", room.Name, "-> Neighbours:", neighboursStr)
}

// Print un chemin en affichant les salles en suivant l'ordre parcouru.
func PrintPath(path []*modules.Room) {
	for i, room := range path {
		if i == len(path)-1 {
			fmt.Print(room.Name + "\n")
		} else {
			fmt.Print(room.Name + "->")
		}
	}
}

// Print les salles de la colonie en prenant en compte les coordonnées (mais sans les liens)
func PrintColony(roomlist []*modules.Room) {
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

// Calcule la ligne la plus basse et la ligne la plus à droite de la colonie.
func calculateSize(roomlist []*modules.Room) (height, width int) {
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

// Print la résolution de l'algorithme.
func PrintResolve(nbAnt int, paths [][]*modules.Room) {
	// On trie les différents chemins utilisés par longueur
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// On utilise calculateTime pour savoir combien de fourmi va être envoyé dans chaque chemin.
	_, antsPerPath := calculateTime(nbAnt, paths)

	// id de la fourmi
	var antIDs []int
	// id du chemin emprunter par la fourmi
	var antPaths []int
	var antPositions []int
	antsSent := 0
	antsFinished := 0

	// Garde en mémoire combien de fourmi a été envoyée dans chaque chemin.
	pathCursor := make([]int, len(paths))

	// Boucle qui tourne tant que toutes les fourmis n'ont pas atteint la fin.
	for turn := 1; antsFinished < nbAnt; turn++ {

		// On déplace les fourmis déjà présentes et on print leur position.
		for i := 0; i < len(antIDs); i++ {
			if antPositions[i] < len(paths[antPaths[i]])-1 {
				// La fourmi avance d'un rang dans le chemin
				antPositions[i]++
				// On récupère le nom de la salle ou elle se trouve maintenant.
				room := paths[antPaths[i]][antPositions[i]].Name
				fmt.Printf("L%d-%s ", antIDs[i], room)
				// On vérifie si elle a atteint la fin ce tour-ci
				if antPositions[i] == len(paths[antPaths[i]])-1 {
					antsFinished++
				}
			}
		}

		// On prépare les nouvelles fourmis (mais sans les envoyés, elle sont positionnés sur le start! )
		for path := range paths {
			// Pour chaque chemin, on compare le nombre de fourmis à envoyer au nombre de fourmis déjà) envoyées.
			if pathCursor[path] < antsPerPath[path] {
				// Une fourmi de plus est placée et sa place est reservée dans le chemin qu'elle va emprunter.
				antsSent++
				pathCursor[path]++
				antIDs = append(antIDs, antsSent)
				antPaths = append(antPaths, path)
				antPositions = append(antPositions, 0)
			}
		}
		// Saute la ligne au changement de tour
		// Le premier tour ne faisant qu'envoyer les premières fourmis sans faire de print, on l'ignore pour éviter une double ligne vide
		if turn > 1 {
			fmt.Println("")
		}
	}
}
