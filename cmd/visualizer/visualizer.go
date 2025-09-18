package main

import (
	"bufio"
	"fmt"
	"image/color"
	"lem-in/colony"
	"lem-in/datas"
	"lem-in/modules"
	"log"
	"math"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var antSprite *ebiten.Image
var dirt *ebiten.Image

func init() {
	img, _, err := ebitenutil.NewImageFromFile("cmd/visualizer/assets/Ants.png")
	if err != nil {
		log.Fatal("Erreur de chargement du sprite de fourmi :", err)
	}
	antSprite = img
	img, _, err = ebitenutil.NewImageFromFile("cmd/visualizer/assets/dirt.png")
	if err != nil {
		log.Fatal("Erreur de chargement du sprite de dirt :", err)
	}
	dirt = img
}

type Visualization struct {
	Rooms       []*modules.Room
	Ants        []*modules.Ant
	Turns       [][]*modules.Ant
	CurrentTurn int
}

func (g *Visualization) Update() error {
	if g.CurrentTurn >= len(g.Turns) {
		return nil
	}

	allDone := true
	for _, ant := range g.Turns[g.CurrentTurn] {
		if ant.T < 1.0 {
			ant.T += 0.02
			allDone = false
		}
	}

	if allDone {
		g.CurrentTurn++
	}

	return nil
}

func ApplyMovements(lines []string, rooms []*modules.Room) ([][]*modules.Ant, []*modules.Ant) {
	var turns [][]*modules.Ant
	antMap := make(map[string]*modules.Ant)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var turn []*modules.Ant
		parts := strings.Fields(line)

		for _, part := range parts {
			if strings.HasPrefix(part, "L") {
				move := strings.Split(part[1:], "-")
				if len(move) == 2 {
					id := move[0]
					destName := move[1]
					dest := colony.GetRoomByName(destName, rooms)

					ant, exists := antMap[id]
					if !exists {
						ant = &modules.Ant{
							Id:          id,
							LastRoom:    rooms[0],
							CurrentRoom: dest,
							T:           0.0,
							Active:      true,
						}
						antMap[id] = ant
					} else {
						ant.LastRoom = ant.CurrentRoom
						ant.CurrentRoom = dest
						ant.T = 0.0
						ant.Active = true
					}

					// Cloner pour ce tour
					cloned := &modules.Ant{
						Id:          ant.Id,
						LastRoom:    ant.LastRoom,
						CurrentRoom: ant.CurrentRoom,
						T:           0.0,
						Active:      true,
					}
					turn = append(turn, cloned)

				}
			}
		}
		turns = append(turns, turn)
	}

	// Extraire la slice globale des fourmis
	var ants []*modules.Ant
	for _, ant := range antMap {
		ants = append(ants, ant)
	}

	return turns, ants
}

func AntMovement(screen *ebiten.Image, from, to *modules.Room, progress float64, sprite *ebiten.Image) {
	if from == nil || to == nil || sprite == nil {
		return
	}

	// Position interpolée
	x := (1-progress)*float64(from.Coordinates.X) + progress*float64(to.Coordinates.X)
	y := (1-progress)*float64(from.Coordinates.Y) + progress*float64(to.Coordinates.Y)

	// Angle de direction
	dx := float64(to.Coordinates.X - from.Coordinates.X)
	dy := float64(to.Coordinates.Y - from.Coordinates.Y)
	angle := math.Atan2(dy, dx)

	// Oscillation douce
	oscillation := math.Sin(progress*10*math.Pi) * (5 * math.Pi / 180) // ±5°

	// Rotation finale
	totalRotation := angle + oscillation

	// Dessin avec rotation
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sprite.Bounds().Dx())/2, -float64(sprite.Bounds().Dy())/2)
	op.GeoM.Rotate(totalRotation)
	op.GeoM.Translate(x, y)
	screen.DrawImage(sprite, op)
}

func (g *Visualization) Draw(screen *ebiten.Image) {
	// Fond marron
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(
		float64(screen.Bounds().Dx())/float64(dirt.Bounds().Dx()),
		float64(screen.Bounds().Dy())/float64(dirt.Bounds().Dy()),
	)
	screen.DrawImage(dirt, op)

	// Dessiner les liens
	for _, room := range g.Rooms {
		for _, neighbor := range room.Neighbours {
			ebitenutil.DrawLine(
				screen,
				float64(room.Coordinates.X),
				float64(room.Coordinates.Y),
				float64(neighbor.Coordinates.X),
				float64(neighbor.Coordinates.Y),
				color.RGBA{124, 180, 50, 255},
			)
		}
	}

	// Dessiner les salles
	for i, room := range g.Rooms {
		var col color.Color = color.White
		if i == 0 {
			col = color.RGBA{0, 255, 0, 255} // Vert
		} else if i == len(g.Rooms)-1 {
			col = color.RGBA{255, 0, 0, 255} // Rouge
		}
		ebitenutil.DrawRect(screen, float64(room.Coordinates.X)-10, float64(room.Coordinates.Y)-10, 20, 20, col)
		ebitenutil.DebugPrintAt(screen, room.Name, room.Coordinates.X+12, room.Coordinates.Y-10)
	}

	if g.CurrentTurn < len(g.Turns) {
		for _, ant := range g.Turns[g.CurrentTurn] {
			AntMovement(screen, ant.LastRoom, ant.CurrentRoom, ant.T, antSprite)
		}
	}

	tourText := fmt.Sprintf("Tour : %d / %d", g.CurrentTurn+1, len(g.Turns)+1)
	ebitenutil.DebugPrintAt(screen, tourText, 700, 10)
}

func (g *Visualization) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func centerColony(rooms []*modules.Room, screenWidth, screenHeight int, margin int) {
	minX, minY := 99999, 99999
	maxX, maxY := 0, 0

	for _, r := range rooms {
		if r.Coordinates.X < minX {
			minX = r.Coordinates.X
		}
		if r.Coordinates.Y < minY {
			minY = r.Coordinates.Y
		}
		if r.Coordinates.X > maxX {
			maxX = r.Coordinates.X
		}
		if r.Coordinates.Y > maxY {
			maxY = r.Coordinates.Y
		}
	}

	colonyWidth := maxX - minX
	colonyHeight := maxY - minY

	// Appliquer la marge en réduisant l’espace disponible
	availableWidth := screenWidth - 2*margin
	availableHeight := screenHeight - 2*margin

	// Calcul du facteur d’échelle pour que tout rentre dans la fenêtre
	scaleX := float64(availableWidth) / float64(colonyWidth)
	scaleY := float64(availableHeight) / float64(colonyHeight)
	scale := math.Min(scaleX, scaleY)

	// Appliquer le centrage et la marge
	for _, r := range rooms {
		r.Coordinates.X = int((float64(r.Coordinates.X-minX) * scale)) + margin
		r.Coordinates.Y = int((float64(r.Coordinates.Y-minY) * scale)) + margin
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var instructions []string
	var mouvements []string
	instructionsDone := false

	for scanner.Scan() {
		line := scanner.Text()
		if !instructionsDone && strings.HasPrefix(line, "L") {
			instructionsDone = true
		}

		if instructionsDone {
			mouvements = append(mouvements, line)
		} else {
			instructions = append(instructions, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	filedatas := datas.SaveDatas(instructions)
	datas.CheckErrors(&filedatas)
	if len(filedatas.Errors) != 0 {
		for _, err := range filedatas.Errors {
			fmt.Println(fmt.Errorf("error : %w", err))
		}
		return
	}
	rooms := colony.CreatRooms(filedatas)
	colony.CreatColony(filedatas, rooms)
	centerColony(rooms, 800, 600, 50)
	turns, ants := ApplyMovements(mouvements, rooms)
	Visualization := &Visualization{
		Rooms:       rooms,
		Turns:       turns,
		Ants:        ants,
		CurrentTurn: 0,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Lem-in Visualizer")
	if err := ebiten.RunGame(Visualization); err != nil {
		log.Fatal(err)
	}
}
