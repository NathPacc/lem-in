// Package modules defines core data structures for the lem-in project.
package modules

// Point represents a 2D coordinate (X, Y) for a room.
type Point struct {
	X int
	Y int
}

// Room represents a room in the colony, with its name, neighbours, and coordinates.
type Room struct {
	Name        string
	Neighbours  []*Room
	Coordinates Point
}
