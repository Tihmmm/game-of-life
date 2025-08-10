package set

import (
	"math/rand"
	"slices"

	"github.com/tihmmm/game-of-life/util"
)

type Set struct {
	Points []Point
}

type Point struct {
	X, Y uint
}

func NewZNByM(n, m uint) *Set {
	n, m = n+1, m+1
	set := new(Set)

	for i := range n {
		for j := range m {
			point := Point{i, j}
			set.Points = append(set.Points, point)
		}
	}

	return set
}

// NewZSubsetRand
// new subset of Z^{2} with randomized x and y values.
// x is in [0, n],
// y is in [0, m].
func NewZSubsetRand(n, m uint) Set {
	n, m = n+1, m+1
	set := Set{}

	numPoints := rand.Intn(int(n * m))
	xSlice := util.CreateRandomUintSlice(numPoints, n)
	ySlice := util.CreateRandomUintSlice(numPoints, m)

	for i, x := range xSlice {
		point := Point{x, ySlice[i]}
		set.Points = append(set.Points, point)
	}

	set.Points = util.RemoveDuplicate(set.Points)
	sortSetByXAsc(&set)
	return set
}

// FindMooreOneNeighborhood
// find Moore neighbours of `point` in the `target` set
func FindMooreOneNeighborhood(point *Point, target *Set) *Set {
	targetNeighbourhood := new(Set)
	allNeighbours := *calculateMooreOneNeighbours(point)
	for _, np := range allNeighbours {
		if slices.Contains(target.Points, np) {
			targetNeighbourhood.Points = append(targetNeighbourhood.Points, np)
		}
	}

	return targetNeighbourhood
}

func calculateMooreOneNeighbours(p *Point) *[]Point {
	x, y := p.X, p.Y
	return &[]Point{
		{x + 1, y}, {x - 1, y}, {x, y + 1}, {x, y - 1},
		{x + 1, y + 1}, {x + 1, y - 1}, {x - 1, y + 1}, {x - 1, y - 1},
	}
}

func sortSetByXAsc(set *Set) {
	slices.SortFunc(
		set.Points,
		func(a, b Point) int {
			if a.X < b.X {
				return -1
			} else if a.X > b.X {
				return 1
			} else {
				if a.Y < b.Y {
					return -1
				} else if a.Y > b.Y {
					return 1
				} else {
					return 0
				}
			}
		},
	)
}
