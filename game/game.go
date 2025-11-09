package game

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/tihmmm/game-of-life/set"
)

type Game struct {
	IsDeterministic bool
	GridSizeX       int
	GridSizeY       int
	Generations     []set.Set
	Board           *set.Set
}

func NewGame(initialState []set.Point, gridSizeX, gridSizeY int, isDeterministic bool) Game {
	game := Game{
		IsDeterministic: isDeterministic,
		GridSizeX:       gridSizeX,
		GridSizeY:       gridSizeY,
	}

	n, m := uint(gridSizeX), uint(gridSizeY)
	game.Board = set.NewZNByM(n, m)

	if initialState == nil || len(initialState) == 0 {
		in := set.NewZSubsetRand(n, m)
		game.Generations = append(game.Generations, in)
	} else {
		game.Generations = append(game.Generations, set.Set{Points: initialState})
	}

	return game
}

func (g *Game) NextGen() set.Set {
	currentGen := g.GetLastGen()
	survivingPoints := GetSurvivingPoints(&currentGen, g.Board)
	currentGen = set.Set{Points: survivingPoints}
	if !g.IsDeterministic {
		deleteRandomPoints(-1, &currentGen)
	}

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

// deleteRandomPoints deletes random n points from s.
// If n < 0, n is generated randomly.
// Returns the number of deleted points and an error
func deleteRandomPoints(n int, s *set.Set) (int, error) {
	if len(s.Points) <= 1 {
		return 0, nil
	}

	if n > len(s.Points) {
		return n, fmt.Errorf("points out of range")
	} else if n < 0 {
		n = rand.Intn(len(s.Points) - 1)
	}

	for range n {
		i := rand.Intn(len(s.Points))
		// remove element from Points slice at index i
		s.Points[i] = s.Points[len(s.Points)-1]
		s.Points = s.Points[:len(s.Points)-1]
	}

	return n, nil
}
