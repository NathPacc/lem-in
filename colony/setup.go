// Package colony provides setup utilities for building the colony structure from input data.
package colony

import (
	"lem-in/modules"
	"strconv"
	"strings"
)

// Créer le lien entre deux salles déjà existantes
func addLink(room1, room2 *modules.Room) {
	if !containsRoom(room1.Neighbours, room2) {
		room1.Neighbours = append(room1.Neighbours, room2)
		room2.Neighbours = append(room2.Neighbours, room1)
	}
}

// Vérifie si une salle est présente dans un chemin/une slice de Rooms
func containsRoom(rooms []*modules.Room, target *modules.Room) bool {
	for _, r := range rooms {
		if r.Name == target.Name {
			return true
		}
	}
	return false
}

// Créer l'ensemble des salles, sans les liens, à partir des datas.
func CreatRooms(datas modules.Datas) []*modules.Room {
	var rooms []*modules.Room
	var tempX int
	var tempY int
	// On créé la salle d'entrée
	tempX, _ = strconv.Atoi(strings.Fields(datas.Start)[1])
	tempY, _ = strconv.Atoi(strings.Fields(datas.Start)[2])
	entry := modules.Room{
		Name: strings.Fields(datas.Start)[0],
		Coordinates: modules.Point{
			X: tempX,
			Y: tempY,
		},
	}

	// On créé la salle de sortie
	tempX, _ = strconv.Atoi(strings.Fields(datas.End)[1])
	tempY, _ = strconv.Atoi(strings.Fields(datas.End)[2])
	exit := modules.Room{
		Name: strings.Fields(datas.End)[0],
		Coordinates: modules.Point{
			X: tempX,
			Y: tempY,
		},
	}
	// On place l'entrée au début de la slice qu'on va retourner
	rooms = append(rooms, &entry)
	// On créé et ajoute toutes les salles intermédiaires
	for _, room := range datas.Rooms {
		tempName := strings.Fields(room)[0]
		tempX, _ = strconv.Atoi(strings.Fields(room)[1])
		tempY, _ = strconv.Atoi(strings.Fields(room)[2])
		temp := modules.Room{
			Name: tempName,
			Coordinates: modules.Point{
				X: tempX,
				Y: tempY,
			},
		}
		rooms = append(rooms, &temp)
	}
	// On ajoute la salle de sortie à la fin de la liste
	rooms = append(rooms, &exit)
	return rooms
}

// Créer l'ensemble des liens d'une colonie (dont les salles ont été créés)
func CreatColony(datas modules.Datas, rooms []*modules.Room) {
	for _, link := range datas.Links {
		left, right, _ := strings.Cut(link, "-")
		var leftRoom *modules.Room
		var rightRoom *modules.Room
		for _, room := range rooms {
			if room.Name == left {
				leftRoom = room
			} else if room.Name == right {
				rightRoom = room
			}
		}
		if leftRoom != nil && rightRoom != nil {
			addLink(leftRoom, rightRoom)
		}
	}
}

// Trouve une salle par son nom (utilsé pour visualizer uniquement)
func GetRoomByName(name string, rooms []*modules.Room) *modules.Room {
	for _, r := range rooms {
		if r.Name == name {
			return r
		}
	}
	return nil
}
