// Package colony provides setup utilities for building the colony structure from input data.
package colony

import (
	"lem-in/modules"
	"strconv"
	"strings"
)

// addLink creates a bidirectional link between two rooms if not already linked.
func addLink(room1, room2 *modules.Room) {
	if !containsRoom(room1.Neighbours, room2) {
		room1.Neighbours = append(room1.Neighbours, room2)
		room2.Neighbours = append(room2.Neighbours, room1)
	}
}

// containsRoom checks if a room slice contains a room with the same name as target.
func containsRoom(rooms []*modules.Room, target *modules.Room) bool {
	for _, r := range rooms {
		if r.Name == target.Name {
			return true
		}
	}
	return false
}

// CreatRooms creates Room structs for the start, end, and all intermediate rooms from input data.
func CreatRooms(datas modules.Datas) []*modules.Room {
	var rooms []*modules.Room
	var tempX int
	var tempY int
	// Create entry room (start)
	tempX, _ = strconv.Atoi(strings.Fields(datas.Start)[1])
	tempY, _ = strconv.Atoi(strings.Fields(datas.Start)[2])
	entry := modules.Room{
		Name: strings.Fields(datas.Start)[0],
		Coordinates: modules.Point{
			X: tempX,
			Y: tempY,
		},
	}

	// Create exit room (end)
	tempX, _ = strconv.Atoi(strings.Fields(datas.End)[1])
	tempY, _ = strconv.Atoi(strings.Fields(datas.End)[2])
	exit := modules.Room{
		Name: strings.Fields(datas.End)[0],
		Coordinates: modules.Point{
			X: tempX,
			Y: tempY,
		},
	}
	rooms = append(rooms, &entry)
	// Create intermediate rooms
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
	rooms = append(rooms, &exit)
	return rooms
}

// CreatColony links rooms together based on the links in the input data.
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
