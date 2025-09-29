// Package colony contains algorithms for pathfinding and optimization in the lem-in project.
package colony

import (
	"lem-in/modules"
	"slices"
	"sort"
	"strings"
)

func FindAllPaths(room *modules.Room, exit *modules.Room, currentPath []*modules.Room) [][]*modules.Room {
	// On ajoute la salle actuelle
	currentPath = append(currentPath, room)

	// Stock tous les chemins valides reliant cette salle à la sortie
	var allPaths [][]*modules.Room

	// Si la salle actuelle est la sortie, on ajoute le chemin qu'on a emprunter pour l'atteindre
	if room == exit {
		allPaths = append(allPaths, currentPath)
		return allPaths
	}

	// Pour chaque voisin de notre salle
	for _, neighbour := range room.Neighbours {
		// Si la salle a déjà été visitée, on l'ignore (sinon on tourne en rond)
		if slices.Contains(currentPath, neighbour) {
			continue
		}
		// On récupère récursivement tous les chemins qui relient ce voisin à la sortie
		subPaths := FindAllPaths(neighbour, exit, currentPath)
		// Si aucun chemin ne relie ce voisin à la sortie, on ignore cette direction (gestion des impasses)
		if subPaths == nil {
			continue
		}
		// Ajoute tous les chemins qui ont atteins la sortie depuis cette salle
		// !!! subPaths peut vous induire en erreur !!!
		// !!! Lorsque l'on atteint le exit, la fonction renvoie le chemin depuis le tout début !!!
		// !!! En effet, on append currentPath à chaque fois qu'on arrive à exit, mais currentPath est le chemin réalisé jusqu'à maintenant !!!
		// !!! En résumé, subPath renvoie tous les chemins qui atteignent la fin en ayant le même début !!!
		allPaths = append(allPaths, subPaths...)
	}
	// On restaure currentPath à l'étape précédente
	currentPath = currentPath[:len(currentPath)-1]
	return allPaths
}

func isRedundant(pathA, pathB []*modules.Room) bool {
	if len(pathA) == 2 || len(pathB) == 2 {
		return false
	}
	if len(pathB) >= len(pathA) {
		return false
	}
	set := make(map[string]bool, len(pathA))
	for _, room := range pathA[1 : len(pathA)-1] {
		set[room.Name] = true
	}
	for _, room := range pathB[1 : len(pathB)-1] {
		if !set[room.Name] {
			return false
		}
	}
	return true
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

// Récupère toutes les combinaisons de chemin indépendants entre eux (ne partagent aucune salle en dehors de l'entrée/sortie)
func IndepPaths(paths [][]*modules.Room) [][][]*modules.Room {
	var allSets [][][]*modules.Room

	// On génère une clé pour chaque groupe de chemin pour l'identifier efficacement quelque soit l'ordre au sein de ce dernier
	seen := make(map[string]bool)

	var explore func(current [][]*modules.Room, start int)
	// Explore est une fonction récursive qui va construire les combinaisons et les ajouter à allSets lorsqu'elles sont terminées
	// Elle est imbriquée dans IndepPaths pour simplifier la gestion de ses arguments.
	explore = func(current [][]*modules.Room, start int) {
		canExtend := false

		for i := start; i < len(paths); i++ {
			compatible := true
			for _, p := range current {
				// Test l'indépendance entre chaque chemin et l'ensemble des chemins de la combinaison actuelle
				if !areIndep(p, paths[i]) {
					compatible = false
					break
				}
			}
			// Si le chemin est bien indépendant avec l'ensemble des chemins de la combinaison actuelle, on l'ajoute à la combinaison
			if compatible {
				// A chaque fois qu'on a ajouté un chemin, on part du principe que la combinaison peut être étendue
				canExtend = true
				explore(append(current, paths[i]), i+1)
			}
		}

		// Le dernier tour du for n'a pas trouvé de chemin compatible donc canExtend est false
		// La combinaison ne peut plus être étendue, on regarde si elle n'existe pas déjà dans notre liste
		// Pour cela, on utilise les clés de chaque chemins
		if !canExtend && len(current) > 0 {
			key := pathSetKey(current)
			if !seen[key] {
				seen[key] = true
				allSets = append(allSets, append([][]*modules.Room{}, current...))
			}
		}
	}
	// On appelle explore avec une combinaison vide pour commencer.
	explore([][]*modules.Room{}, 0)
	return allSets
}

// Vérifie que deux chemins ne partagent pas de salle autre que start et end
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

func pathSetKey(paths [][]*modules.Room) string {
	var keys []string
	for _, p := range paths {
		var names []string
		for _, r := range p {
			names = append(names, r.Name)
		}
		keys = append(keys, strings.Join(names, "-"))
	}
	sort.Strings(keys) // ordre canonique
	return strings.Join(keys, "|")
}

// Calcule le temps de résolution d'une combinaison de chemins
func calculateTime(nbAnt int, paths [][]*modules.Room) (int, []int) {
	// Trie les chemins par longueur
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Enregistre le nombre de fourmis à envoyer dans chaque chemin
	antsPerPath := make([]int, len(paths))

	// Tant qu'il reste des fourmis
	for nbAnt > 0 {
		// On regarde quelle chemin est le plus rapide pour cette fourmi là
		bestIndex := 0
		bestTime := len(paths[0]) + antsPerPath[0]
		// Pour chaque chemin, on compare le temps de parcours
		for i := 1; i < len(paths); i++ {
			// Le temps de parcours se calcule à partir de la longueur du chemin et le nombre de fourmi déjà envoyé dedans
			t := len(paths[i]) + antsPerPath[i]
			// Si t est meilleur que le bestTime jusque maintenant, on enregistre son index et son temps
			if t < bestTime {
				bestTime = t
				bestIndex = i
			}
		}
		// On envoie la fourmi dans ce chemin
		antsPerPath[bestIndex]++
		nbAnt--
	}

	// On compare le temps utilisé par chaque chemin pour savoir quand la dernière fourmi arrivera
	maxTime := 0
	for i := range paths {
		t := len(paths[i]) + antsPerPath[i] - 1
		if t > maxTime {
			maxTime = t
		}
	}

	return maxTime, antsPerPath
}

// Calcule le temps de résolution de toutes les combinaisons d'une colonie et renvoie la plus rapide
func Resolve(nbAnt int, colony []*modules.Room) [][]*modules.Room {
	paths := OptimizePaths(FindAllPaths(colony[0], colony[len(colony)-1], nil))
	for _, p := range paths {
		PrintPath(p)
	}
	indepPaths := IndepPaths(paths)
	bestset := indepPaths[0]
	// On initialise le meilleur temps
	bestTime, _ := calculateTime(nbAnt, indepPaths[0])
	for _, set := range indepPaths {
		// Si une combinaison de chemin est plus rapide à traverser, on la sauvegarde
		time, _ := calculateTime(nbAnt, set)
		if time < bestTime {
			bestset = set
			bestTime = time
		}
	}
	return bestset
}
