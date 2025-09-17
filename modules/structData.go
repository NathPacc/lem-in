// Package modules defines core data structures for the lem-in project.
package modules

// Datas holds all parsed input data for the colony, including ants, rooms, links, and errors.
type Datas struct {
	NbAnts int      // Number of ants
	Start  string   // Start room definition
	End    string   // End room definition
	Rooms  []string // List of room definitions
	Links  []string // List of link definitions
	Errors []error  // List of errors found during parsing/validation
}
