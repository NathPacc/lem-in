// Package datas provides data loading and parsing utilities for the lem-in project.
package datas

import (
	"bufio"
	"errors"
	"lem-in/modules"
	"log"
	"os"
	"strconv"
	"strings"
)

// Stock les instructions lignes par lignes
func GetDatas(filename string) []string {
	file, err := os.Open("files/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var results []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		results = append(results, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return results
}

// Répartie les instructions dans la struct data
func SaveDatas(filecontent []string) modules.Datas {
	var datas modules.Datas
	var err error
	// isRoom permet de savoir si l'on est sur une ligne qui définie une salle ou non
	isRoom := true
	// isStart, doubleStart et leurs cousins pour end servent à trouver le start/end et vérifier son existence unique.
	isStart := false
	doubleStart := false
	doubleEnd := false
	isEnd := false
	for i, line := range filecontent {
		// Si une ligne est vide on l'ignore
		if line == "" {
			continue
		}
		// La première ligne est forcément le nombre de fourmis. On vérifira plus tard que le nombre est logique.
		if i == 0 {
			datas.NbAnts, err = strconv.Atoi(line)
			if err != nil {
				datas.Errors = append(datas.Errors, errors.New("Bad format for number of ants"))
				return datas
			}
			continue
		}

		// Vérifie que start/end a bien été rencontré, n'est pas en double et que la ligne a un format valide pour une salle.
		if isStart && !doubleStart && checkRoomFormat(line) == "" {
			datas.Start = line
			doubleStart = true
			continue
		}
		if isEnd && !doubleEnd && checkRoomFormat(line) == "" {
			datas.End = line
			doubleEnd = true
			continue
		}

		// Localise les marqueurs start et end.
		if line == "##start" {
			if doubleStart {
				datas.Errors = append(datas.Errors, errors.New("More than one start"))
				continue
			}
			isStart = true
			continue
		}
		if line == "##end" {
			isEnd = true
			if doubleEnd {
				datas.Errors = append(datas.Errors, errors.New("More than one end"))
				continue
			}
			continue
		}

		// Si la ligne commence par un #, c'est un commentaire qu'on ignore
		if rune(line[0]) == '#' {
			continue
		}

		// Lorsque l'on croise un tiret, on entre dans la définition des liens.
		if strings.Contains(line, "-") {
			isRoom = false
		}

		// Si l'on est encore sur une ligne de room
		if isRoom {
			// Si la ligne n'est pas valide, on ajoute une erreur
			if checkRoomFormat(line) != "" {
				datas.Errors = append(datas.Errors, errors.New(checkRoomFormat(line)))
				continue
			}
			datas.Rooms = append(datas.Rooms, line)
			continue
			// Si ce n'est pas une ligne de salle, alors on la stock comme un lien.
		} else if !isRoom {
			if strings.Contains(line, "-") {
				datas.Links = append(datas.Links, line)
			}
			continue
		}
	}
	// Si doubleEnd/doubleStart est resté false, c'est qu'aucun start ou aucun end n'a été rencontré.
	if !doubleEnd {
		datas.Errors = append(datas.Errors, errors.New("No end"))
	}
	if !doubleStart {
		datas.Errors = append(datas.Errors, errors.New("No start"))
	}
	return datas
}
