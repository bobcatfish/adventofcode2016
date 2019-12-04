package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func load() ([][]string, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, fmt.Errorf("couldn't open file: %v", err)
	}
	defer file.Close()

	vals := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := scanner.Text()
		v := strings.Split(val, ",")
		vals = append(vals, v)
	}
	return vals, err
}

func up(p point) point {
	return point{
		x: p.x,
		y: p.y - 1,
	}
}

func down(p point) point {
	return point{
		x: p.x,
		y: p.y + 1,
	}
}

func left(p point) point {
	return point{
		x: p.x - 1,
		y: p.y,
	}
}

func right(p point) point {
	return point{
		x: p.x + 1,
		y: p.y,
	}
}

type goFunc = func(point) point

func follow(c string, curr point) ([]point, point, error) {
	points := []point{}
	cc := c[0]
	n, err := strconv.Atoi(c[1:])
	if err != nil {
		return points, curr, err
	}
	var f goFunc
	switch cc {
	case 'U':
		f = up
	case 'D':
		f = down
	case 'L':
		f = left
	case 'R':
		f = right
	default:
		return points, curr, fmt.Errorf("unknown command %s", cc)
	}
	for i := 0; i < n; i++ {
		pp := f(curr)
		curr = pp
		points = append(points, pp)
	}
	return points, curr, nil
}

func linePoints(line []string) (map[point]struct{}, error) {
	points := map[point]struct{}{}
	curr := point{0, 0}
	for _, c := range line {
		var p []point
		var err error
		p, curr, err = follow(c, curr)
		if err != nil {
			return points, err
		}
		for _, pp := range p {
			points[pp] = struct{}{}
		}
	}
	return points, nil
}

func getCrossings(l1 map[point]struct{}, l2 map[point]struct{}) []point {
	p := []point{}
	for pp, _ := range l1 {
		if _, ok := l2[pp]; ok {
			p = append(p, pp)
		}
	}
	return p
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func getClosest(p []point) (point, int) {
	c := p[0]
	d := abs(p[0].x) + abs(p[0].y)

	for _, pp := range p {
		dd := abs(pp.x) + abs(pp.y)
		if dd < d {
			c = pp
			d = dd
		}
	}
	return c, d
}

func main() {
	vals, err := load()
	if err != nil {
		log.Fatalf("Couldn't load nums from file: %v", err)
	}
	fmt.Println(vals)
	l1, err := linePoints(vals[0])
	if err != nil {
		log.Fatalf("oh no %v", err)
	}
	l2, err := linePoints(vals[1])
	if err != nil {
		log.Fatalf("oh no %v", err)
	}

	fmt.Println(l1)
	fmt.Println(l2)

	crossings := getCrossings(l1, l2)

	fmt.Println(crossings)

	closest, distance := getClosest(crossings)

	fmt.Println(closest, distance)
}
