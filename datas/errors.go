package datas

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

func CheckErrors(datas *Datas) {
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

func checkExtrimities(datas *Datas) {
	if checkRoomFormat(datas.Start) != "" && datas.Start != "" {
		datas.Errors = append(datas.Errors, errors.New("Bad format for start"))
	}
	if checkRoomFormat(datas.End) != "" && datas.End != "" {
		datas.Errors = append(datas.Errors, errors.New("Bad format for end"))
	}
}

func checkRooms(datas *Datas) {
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

func checkLinks(datas *Datas) {
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

func checkRoomFormat(line string) string {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return "Bad format : comment without # : " + line
	}
	_, err1 := strconv.Atoi(parts[1])
	_, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return "Bad format for room : " + line
	}
	return ""
}

func checkDuplicates(datas *Datas) {
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
		if datas.Start != "" {
			startroom := strings.Fields(datas.Start)[0]
			if roompart[0] == startroom {
				datas.Errors = append(datas.Errors, errors.New("Duplicate for rooms "+room+" and "+datas.Start))
			}
		}
		if datas.End != "" {
			endroom := strings.Fields(datas.End)[0]
			if roompart[0] == endroom {
				datas.Errors = append(datas.Errors, errors.New("Duplicate for rooms "+room+" and "+datas.Start))
			}
		}
	}
}
