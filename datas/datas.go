package datas

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

type Datas struct {
	NbAnts int
	Start  string
	End    string
	Rooms  []string
	Links  []string
	Errors []error
}

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

func SaveDatas(filecontent []string) Datas {
	var datas Datas
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
		// On récupère le nombre de fourmis
		if i == 0 {
			datas.NbAnts, err = strconv.Atoi(line)
			if err != nil {
				datas.Errors = append(datas.Errors, errors.New("Bad format for number of ants"))
				return datas
			}
			continue
		}

		// On vérifie si l'on est en train d'ajouter le start ou le end
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
		// On vérifie si la prochain ligne est un start ou un end

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
