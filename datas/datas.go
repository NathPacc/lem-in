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

// GetDatas reads the input file and returns its lines as a slice of strings.
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

// SaveDatas parses the file content and fills a Datas struct with ants, rooms, links, and errors.
func SaveDatas(filecontent []string) modules.Datas {
	var datas modules.Datas
	var err error
	isRoom := true
	isStart := false
	doubleStart := false
	doubleEnd := false
	isEnd := false
	for i, line := range filecontent {
		if line == "" {
			continue
		}
		// Parse the number of ants
		if i == 0 {
			datas.NbAnts, err = strconv.Atoi(line)
			if err != nil {
				datas.Errors = append(datas.Errors, errors.New("Bad format for number of ants"))
				return datas
			}
			continue
		}

		// Check if we are adding the start or end room
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

		// Check for start/end markers
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
		if strings.Contains(line, "-") {
			isRoom = false
		}
		if rune(line[0]) == '#' {
			continue
		}
		// Parse rooms and links
		if isRoom {
			if checkRoomFormat(line) != "" {
				datas.Errors = append(datas.Errors, errors.New(checkRoomFormat(line)))
				continue
			}
			if checkRoomFormat(line) == "comment" {
				continue
			}
			datas.Rooms = append(datas.Rooms, line)
			continue
		} else if !isRoom {
			if strings.Contains(line, "-") {
				datas.Links = append(datas.Links, line)
			}
			continue
		}
	}
	if !doubleEnd {
		datas.Errors = append(datas.Errors, errors.New("No end"))
	}
	if !doubleStart {
		datas.Errors = append(datas.Errors, errors.New("No start"))
	}
	return datas
}
