package colony

import (
	"lem-in/datas"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Room struct {
	Name        string
	Neighbours  []*Room
	Coordinates Point
}

func addLink(room1, room2 *Room) {
	if !containsRoom(room1.Neighbours, room2) {
		room1.Neighbours = append(room1.Neighbours, room2)
		room2.Neighbours = append(room2.Neighbours, room1)
	}
}

func containsRoom(rooms []*Room, target *Room) bool {
	for _, r := range rooms {
		if r.Name == target.Name {
			return true
		}
	}
	return false
}

func CreatRooms(datas datas.Datas) []*Room {
	var rooms []*Room
	var tempX int
	var tempY int
	tempX, _ = strconv.Atoi(strings.Fields(datas.Start)[1])
	tempY, _ = strconv.Atoi(strings.Fields(datas.Start)[2])
	entry := Room{
		Name: strings.Fields(datas.Start)[0],
		Coordinates: Point{
			X: tempX,
			Y: tempY,
		},
	}

	tempX, _ = strconv.Atoi(strings.Fields(datas.End)[1])
	tempY, _ = strconv.Atoi(strings.Fields(datas.End)[2])
	exit := Room{
		Name: strings.Fields(datas.End)[0],
		Coordinates: Point{
			X: tempX,
			Y: tempY,
		},
	}
	rooms = append(rooms, &entry)
	for _, room := range datas.Rooms {
		tempName := strings.Fields(room)[0]
		tempX, _ = strconv.Atoi(strings.Fields(room)[1])
		tempY, _ = strconv.Atoi(strings.Fields(room)[2])
		temp := Room{
			Name: tempName,
			Coordinates: Point{
				X: tempX,
				Y: tempY,
			},
		}
		rooms = append(rooms, &temp)
	}
	rooms = append(rooms, &exit)
	return rooms
}

func CreatColony(datas datas.Datas, rooms []*Room) {
	for _, link := range datas.Links {
		left, right, _ := strings.Cut(link, "-")
		var leftRoom *Room
		var rightRoom *Room
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
