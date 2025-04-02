package game

import (
	"github.com/tihmmm/game-of-life/set"
	"log"
	"slices"
)

type Game struct {
	Generations []set.Set
	GridSizeX   int
	GridSizeY   int
	Board       *set.Set
}

func NewGameWithNGenerations(initialState []set.Point, gridSizeX, gridSizeY, generationsNum int) *Game {
	game := Game{
		GridSizeX: gridSizeX,
		GridSizeY: gridSizeY,
	}

	n, m := uint(gridSizeX), uint(gridSizeY)
	game.Board = set.NewZNByM(n, m)

	if initialState == nil {
		in := set.NewZSubsetRand(n, m)
		game.Generations = append(game.Generations, in)
	} else {
		game.Generations = append(game.Generations, set.Set{Points: initialState})
	}

	log.Printf("board:%dx%d\ninitial state:\n%v\n", gridSizeX, gridSizeY, game.Generations[0])

	game.NextNGens(generationsNum, &game.Generations[0], game.Board)

	return &game
}

func (g *Game) NextNGens(n int, currentGen, board *set.Set) {
	gen := &set.Set{Points: currentGen.Points}
	for range n {
		gen = g.NextGen(gen, board)
	}
}

func (g *Game) NextGen(currentGen, board *set.Set) *set.Set {
	survivingPoints := GetSurvivingPoints(currentGen, board)
	currentGen = &set.Set{Points: survivingPoints}
	g.Generations = append(g.Generations, *currentGen)

	return currentGen
}

func GetSurvivingPoints(currentGen, board *set.Set) []set.Point {
	var survivingPoints []set.Point

	for _, boardPoint := range board.Points {
		neighbours := set.FindMooreOneNeighborhood(&boardPoint, currentGen)
		neighboursNum := len(neighbours.Points)
		if neighboursNum > 3 || neighboursNum < 2 {
			continue
		} else if slices.Contains(currentGen.Points, boardPoint) {
			survivingPoints = append(survivingPoints, boardPoint)
		} else {
			if neighboursNum == 2 {
				continue
			} else if neighboursNum == 3 {
				survivingPoints = append(survivingPoints, boardPoint)
			}
		}
	}

	return survivingPoints
}
