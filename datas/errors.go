package datas

import (
	"errors"
	"lem-in/modules"
	"slices"
	"strconv"
	"strings"
)

// CheckErrors runs all error checks on the provided datas structure.
// It checks ants count, room and link formats, and duplicates.
func CheckErrors(datas *modules.Datas) {
	if datas.NbAnts <= 0 {
		datas.Errors = append(datas.Errors, errors.New("Bad format for number of ant"))
	}
	checkExtrimities(datas)
	if len(datas.Rooms) != 0 {
		checkRooms(datas)
		checkLinks(datas)
		checkDuplicates(datas)
	}
}

// checkExtrimities checks the format of the start and end rooms.
func checkExtrimities(datas *modules.Datas) {
	if checkRoomFormat(datas.Start) != "" && datas.Start != "" {
		datas.Errors = append(datas.Errors, errors.New("Bad format for start"))
	}
	if checkRoomFormat(datas.End) != "" && datas.End != "" {
		datas.Errors = append(datas.Errors, errors.New("Bad format for end"))
	}
}

// checkRooms checks the format of each room string in datas.Rooms.
// Each room should have three fields, and the last two should be integers.
func checkRooms(datas *modules.Datas) {
	for _, roomstr := range datas.Rooms {
		roomtab := strings.Fields(roomstr)
		if len(roomtab) != 3 {
			datas.Errors = append(datas.Errors, errors.New("Bad format for the following room : "+roomstr))
			continue
		}
		_, err1 := strconv.Atoi(roomtab[1])
		_, err2 := strconv.Atoi(roomtab[2])
		if err1 != nil || err2 != nil {
			datas.Errors = append(datas.Errors, errors.New("Bad format for the following room : "+roomstr))
			continue
		}
	}
}

// checkLinks checks the format and validity of each link in datas.Links.
// It ensures links are between two different, existing rooms.
func checkLinks(datas *modules.Datas) {
	if datas.Start == "" || datas.End == "" {
		return
	}
	for _, link := range datas.Links {
		left, right, found := strings.Cut(link, "-")
		if !found {
			datas.Errors = append(datas.Errors, errors.New("Bad format for the following link : "+link))
			continue
		}
		if left == right {
			datas.Errors = append(datas.Errors, errors.New("Bad format for the following link : "+link))
			continue
		}
		leftExist := false
		rightExist := false
		// Check if both rooms in the link exist in the rooms list or as start/end
		for _, roomstr := range datas.Rooms {
			if strings.Fields(roomstr)[0] == left || strings.Fields(datas.Start)[0] == left || strings.Fields(datas.End)[0] == left {
				leftExist = true
			}
			if strings.Fields(roomstr)[0] == right || strings.Fields(datas.Start)[0] == right || strings.Fields(datas.End)[0] == right {
				rightExist = true
			}
		}
		if !leftExist || !rightExist {
			datas.Errors = append(datas.Errors, errors.New("Bad format for the following link : "+link))
			continue
		}
	}
}

// checkRoomFormat checks if a room line is valid (3 fields, last two are integers).
// Returns an error string if invalid, or empty string if valid.
func checkRoomFormat(line string) string {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return "Bad format : comment without # : " + line
	}
	_, err1 := strconv.Atoi(parts[1])
	_, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return "Bad format for room : " + line
	}
	return ""
}

// checkDuplicates checks for duplicate room names in datas.Rooms, Start, and End.
// Adds errors for any duplicates found.
func checkDuplicates(datas *modules.Datas) {
	var duplicatesIndex []int
	for i, room := range datas.Rooms {
		roompart := strings.Fields(room)
		for j, comparative := range datas.Rooms {
			if i != j && !slices.Contains(duplicatesIndex, j) {
				comparativepart := strings.Fields(comparative)
				if roompart[0] == comparativepart[0] {
					duplicatesIndex = append(duplicatesIndex, i)
					datas.Errors = append(datas.Errors, errors.New("Duplicate for rooms "+room+" and "+comparative))
				}
			}
		}
		// Check for duplicates with start room
		if datas.Start != "" {
			startroom := strings.Fields(datas.Start)[0]
			if roompart[0] == startroom {
				datas.Errors = append(datas.Errors, errors.New("Duplicate for rooms "+room+" and "+datas.Start))
			}
		}
		// Check for duplicates with end room
		if datas.End != "" {
			endroom := strings.Fields(datas.End)[0]
			if roompart[0] == endroom {
				datas.Errors = append(datas.Errors, errors.New("Duplicate for rooms "+room+" and "+datas.End))
			}
		}
	}
}
