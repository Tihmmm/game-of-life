package game

import (
	"log"
	"slices"

	"github.com/tihmmm/game-of-life/set"
)

type Game struct {
	Generations []set.Set
	GridSizeX   int
	GridSizeY   int
	Board       *set.Set
}

func NewGame(initialState []set.Point, gridSizeX, gridSizeY int) Game {
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

	return game
}

func (g *Game) NextGen() set.Set {
	currentGen := g.GetLastGen()
	survivingPoints := GetSurvivingPoints(&currentGen, g.Board)
	currentGen = set.Set{Points: survivingPoints}
	g.Generations = append(g.Generations, currentGen)

	return currentGen
}

func (g *Game) GetNthGeneration(n int) *set.Set {
	return &g.Generations[n]
}

func (g *Game) GetLastGen() set.Set {
	return g.Generations[len(g.Generations)-1]
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
