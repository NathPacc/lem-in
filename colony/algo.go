package colony

import (
	"slices"
	"sort"
)

func FindAllPaths(entry *Room, exit *Room, currentPath []*Room) [][]*Room {
	copyPath := append([]*Room{}, currentPath...)
	copyPath = append(copyPath, entry)

	var allPaths [][]*Room

	if entry == exit {
		allPaths = append(allPaths, copyPath)
		return allPaths
	}

	for _, neighbour := range entry.Neighbours {
		if slices.Contains(copyPath, neighbour) {
			continue
		}
		subPaths := FindAllPaths(neighbour, exit, copyPath)
		allPaths = append(allPaths, subPaths...)
	}

	return allPaths
}

func isRedundant(pathA, pathB []*Room) bool {
	// On ne considère pas l'entrée et la sortie
	set := make(map[string]bool)
	for _, room := range pathB[1 : len(pathB)-1] {
		set[room.Name] = true
	}

	for _, room := range pathA[1 : len(pathA)-1] {
		if !set[room.Name] {
			return false
		}
	}

	// Si pathA est plus long, alors il est redondant
	return len(pathA) > len(pathB)
}

func OptimizePaths(paths [][]*Room) [][]*Room {
	var results [][]*Room
	for i, pathA := range paths {
		redundant := false
		for j, pathB := range paths {
			if i != j && isRedundant(pathA, pathB) {
				redundant = true
				break
			}
		}
		if !redundant {
			results = append(results, pathA)
		}
	}
	return results
}

func IndepPaths(paths [][]*Room) [][][]*Room {
	var allSets [][][]*Room

	var explore func(current [][]*Room, start int)
	explore = func(current [][]*Room, start int) {
		if len(current) > 0 {
			// Vérifie que ce groupe n'existe pas déjà
			duplicate := false
			for _, existing := range allSets {
				if sameSet(existing, current) {
					duplicate = true
					break
				}
			}
			if !duplicate {
				allSets = append(allSets, append([][]*Room{}, current...))
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

	explore([][]*Room{}, 0)
	return allSets
}

func areIndep(pathA, pathB []*Room) bool {
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

func sameSet(a, b [][]*Room) bool {
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

func samePath(a, b []*Room) bool {
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

func calculateTime(nbAnt int, paths [][]*Room) (int, []int) {
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

func Resolve(nbAnt int, colony []*Room) [][]*Room {
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
