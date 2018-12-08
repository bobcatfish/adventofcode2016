package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type startingPoint struct {
	// Adjadcent points, indexed by their distance away.
	// At each distance, we have a map where the keys are the
	// adjacent points.
	adjByDist map[int]map[point]bool
	area      int
	p         point

	id rune
}

// adjCount can be used to updtae a startingPoint if another
// adj point is found at the same location
type adjCount struct {
	adj          map[point]bool
	startingArea *int
	startPoint   point
	id           rune
	// shouldnt matter cuz each new one we encounter is either less close or the same
	dist int
}

var tombstone = adjCount{dist: -1, id: '.'}

func getAdj(p point) []point {
	return []point{
		// above
		{x: p.x, y: p.y + 1},

		// beside
		{x: p.x - 1, y: p.y},
		{x: p.x + 1, y: p.y},

		// below
		{x: p.x, y: p.y - 1},
	}

}

func stepOut(start *startingPoint, dist int, adj map[point]bool, found map[point]adjCount) {
	start.adjByDist[dist] = map[point]bool{}
	for p := range adj {
		for _, adjP := range getAdj(p) {
			start.adjByDist[dist][adjP] = true
			prev, ok := found[adjP]
			if ok {
				if prev.dist != tombstone.dist && prev.startPoint != start.p && prev.dist >= dist {
					*prev.startingArea--
					found[adjP] = tombstone
				}
			} else {
				start.area++
				found[adjP] = adjCount{
					startPoint:   start.p,
					adj:          start.adjByDist[dist],
					startingArea: &start.area,
					dist:         dist,
					id:           start.id,
				}
			}
		}
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	defer file.Close()

	vals := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		vals = append(vals, scanner.Text())
	}

	starts := []startingPoint{}
	for i, v := range vals {
		vs := strings.Split(v, ", ")
		x, err := strconv.Atoi(vs[0])
		if err != nil {
			log.Fatalf("Error converting %q to point: %v", v, err)
		}
		y, err := strconv.Atoi(vs[1])
		if err != nil {
			log.Fatalf("Error converting %q to point: %v", v, err)
		}
		startPoint := point{
			x: x,
			y: y,
		}
		start := startingPoint{
			p: startPoint,
			adjByDist: map[int]map[point]bool{
				0: map[point]bool{
					startPoint: true,
				},
			},
			area: 1,
			id:   rune('A' + i),
		}
		starts = append(starts, start)
	}

	// The keys here are points that are adjacent to other points.
	// The values are the maps from that starting point that is adjacent,
	// so that if we find another adjacent point, we can remove them
	found := map[point]adjCount{}

	// seed the known points in
	for i := range starts {
		start := starts[i]
		found[start.p] = adjCount{
			startPoint:   start.p,
			adj:          start.adjByDist[0],
			startingArea: &start.area,
			dist:         0,
			id:           start.id,
		}
	}
	// this will be used to detect values that aren't infinitely increasing
	prevCounts := map[point]int{}
	for i := 1; ; i++ {
		for j := range starts {
			start := &starts[j]
			stepOut(start, i, start.adjByDist[i-1], found)
		}

		// This is the "best" idea I have for detecting infinity
		fmt.Printf("%d...", i)
		if i%10 == 0 {
			fmt.Println()
			candidates := []int{}
			for _, start := range starts {
				if prevCounts[start.p] == start.area {
					candidates = append(candidates, start.area)
				}
			}
			if len(candidates) > 0 {
				sort.Ints(candidates)
				fmt.Println(candidates[len(candidates)-1])
			}
		}

		for _, start := range starts {
			prevCounts[start.p] = start.area
		}

		/*
			if i == 10 {
				grid := map[point]rune{}

				for x := 0; x < 10; x++ {
					for y := 0; y < 10; y++ {
						grid[point{x: x, y: y}] = '.'
					}
				}
				for adj, f := range found {
					grid[adj] = f.id
				}

				for y := 0; y < 10; y++ {
					for x := 0; x < 10; x++ {
						fmt.Printf("%s ", string(grid[point{x: x, y: y}]))
					}
					fmt.Println()
				}
				break
			}
		*/
	}
}