// Package colony contains algorithms for pathfinding and optimization in the lem-in project.
package colony

import (
	"lem-in/modules"
	"slices"
	"sort"
)

func FindAllPaths(room *modules.Room, exit *modules.Room, currentPath []*modules.Room) [][]*modules.Room {
	// Copie le chemin emprunter jusqu'à la salle actuelle pour gérer correctement les modifications de slices.
	copyPath := append([]*modules.Room{}, currentPath...)
	// On ajoute la salle actuelle
	copyPath = append(copyPath, room)

	// Stock tous les chemins valides reliant cette salle à la sortie
	var allPaths [][]*modules.Room

	// Si la salle actuelle est la sortie, on ajoute le chemin qu'on a emprunter pour l'atteindre
	if room == exit {
		allPaths = append(allPaths, copyPath)
		return allPaths
	}

	// Pour chaque voisin de notre salle
	for _, neighbour := range room.Neighbours {
		// Si la salle a déjà été visitée, on l'ignore (sinon on tourne en rond)
		if slices.Contains(copyPath, neighbour) {
			continue
		}
		// On récupère récursivement tous les chemins qui relient ce voisin à la sortie
		subPaths := FindAllPaths(neighbour, exit, copyPath)
		// Si aucun chemin ne relie ce voisin à la sortie, on ignore cette direction (gestion des impasses)
		if subPaths == nil {
			continue
		}
		// Ajoute tous les chemins qui ont atteins la sortie depuis cette salle
		// !!! subPaths peut vous induire en erreur !!!
		// !!! Lorsque l'on atteint le exit, la fonction renvoie le chemin depuis le tout début !!!
		// !!! En effet, on append copyPath à chaque fois qu'on arrive à exit, mais copyPath est le chemin réalisé jusqu'à maintenant !!!
		// !!! En résumé, subPath renvoie tous les chemins qui atteignent la fin en ayant le même début !!!
		allPaths = append(allPaths, subPaths...)
	}
	return allPaths
}

// Vérifie si un chemin A est inclus dans un chemin B, le rendant inutile
func isRedundant(pathA, pathB []*modules.Room) bool {
	set := make(map[string]bool)
	for _, room := range pathB[1 : len(pathB)-1] {
		set[room.Name] = true
	}

	for _, room := range pathA[1 : len(pathA)-1] {
		if !set[room.Name] {
			return false
		}
	}

	return len(pathA) > len(pathB)
}

// Élimine les chemins "redondants", c'est à dire tous les chemins incluant un chemin valide (ex : A -> B -> C est redondant avec A -> C)
func OptimizePaths(paths [][]*modules.Room) [][]*modules.Room {
	var results [][]*modules.Room
	for i, pathA := range paths {
		redundant := false
		for j, pathB := range paths {
			// Si un chemin est redondant, on ne le compare plus avec le reste et on ne le recopie pas
			if i != j && isRedundant(pathA, pathB) {
				redundant = true
				break
			}
		}
		// Si un chemin n'est pas redondant avec tous les autres, on le recopie.
		if !redundant {
			results = append(results, pathA)
		}
	}
	return results
}

// IndepPaths finds all sets of independent (non-overlapping) paths.
func IndepPaths(paths [][]*modules.Room) [][][]*modules.Room {
	var allSets [][][]*modules.Room

	var explore func(current [][]*modules.Room, start int)
	explore = func(current [][]*modules.Room, start int) {
		if len(current) > 0 {
			// Check if this group already exists
			duplicate := false
			for _, existing := range allSets {
				if sameSet(existing, current) {
					duplicate = true
					break
				}
			}
			if !duplicate {
				allSets = append(allSets, append([][]*modules.Room{}, current...))
			}
		}

		for i := start; i < len(paths); i++ {
			compatible := true
			for _, p := range current {
				if !areIndep(p, paths[i]) {
					compatible = false
					break
				}
			}
			if compatible {
				explore(append(current, paths[i]), i+1)
			}
		}
	}

	explore([][]*modules.Room{}, 0)
	return allSets
}

// areIndep returns true if two paths are independent (no shared rooms except entry/exit).
func areIndep(pathA, pathB []*modules.Room) bool {
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, room := range pathA[1 : len(pathA)-1] {
		set1[room.Name] = true
	}
	for _, room := range pathB[1 : len(pathB)-1] {
		set2[room.Name] = true
	}

	for name := range set1 {
		if set2[name] {
			return false
		}
	}
	return true
}

// sameSet checks if two sets of paths are equivalent (same paths, any order).
func sameSet(a, b [][]*modules.Room) bool {
	if len(a) != len(b) {
		return false
	}

	used := make([]bool, len(b))

	for _, pathA := range a {
		found := false
		for i, pathB := range b {
			if used[i] {
				continue
			}
			if samePath(pathA, pathB) {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// samePath checks if two paths are identical (same rooms in order).
func samePath(a, b []*modules.Room) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}
	return true
}

// calculateTime distributes ants across paths to minimize the total time.
// Returns the max time and the number of ants per path.
func calculateTime(nbAnt int, paths [][]*modules.Room) (int, []int) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	antsPerPath := make([]int, len(paths))

	for nbAnt > 0 {
		bestIndex := 0
		bestTime := len(paths[0]) + antsPerPath[0]
		for i := 1; i < len(paths); i++ {
			t := len(paths[i]) + antsPerPath[i]
			if t < bestTime {
				bestTime = t
				bestIndex = i
			}
		}
		antsPerPath[bestIndex]++
		nbAnt--
	}

	maxTime := 0
	for i := range paths {
		t := len(paths[i]) + antsPerPath[i] - 1
		if t > maxTime {
			maxTime = t
		}
	}

	return maxTime, antsPerPath
}

// Resolve finds the best set of independent paths for the given number of ants.
// Returns the set of paths that minimizes the total time.
func Resolve(nbAnt int, colony []*modules.Room) [][]*modules.Room {
	paths := OptimizePaths(FindAllPaths(colony[0], colony[len(colony)-1], nil))
	indepPaths := IndepPaths(paths)
	bestset := indepPaths[0]
	bestTime, _ := calculateTime(nbAnt, indepPaths[0])
	for _, set := range indepPaths {
		time, _ := calculateTime(nbAnt, set)
		if time < bestTime {
			bestset = set
			bestTime = time
		}
	}
	return bestset
}
