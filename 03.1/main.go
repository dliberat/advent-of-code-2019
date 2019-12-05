/*  https://adventofcode.com/2019/day/3  */

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type step struct {
	dir  byte
	dist int
}

type coord struct {
	x int
	y int
}

func (c *coord) equals(other *coord) bool {
	return c.x == other.x && c.y == other.y
}

func (c *coord) manhattanDistance() int {
	return abs(c.x) + abs(c.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

func newStep(txt string) (step, error) {

	s := step{dir: txt[0]}
	dist, err := strconv.Atoi(txt[1:])

	if err != nil {
		return s, err
	}

	s.dist = dist
	return s, nil
}

func stringToCoord(txt string) coord {
	xy := strings.Split(txt, ",")
	x, _ := strconv.Atoi(xy[0])
	y, _ := strconv.Atoi(xy[1])
	return coord{x: x, y: y}
}

func getRoute(input string) map[string]bool {
	// Convert a line of inputs into steps
	txt := strings.Split(input, ",")
	steps := make([]step, len(txt))
	for i, v := range txt {
		s, err := newStep(v)
		if err != nil {
			fmt.Println("ERROR: Cannot parse input: ", v)
			return make(map[string]bool)
		}
		steps[i] = s
	}

	// trace the steps to figure out every single coordinate
	// that gets touched on the map
	route := make(map[string]bool)
	currentpos := coord{x: 0, y: 0}

	var nextX int
	var nextY int

	for _, v := range steps {

		if v.dir == 85 { // up
			nextX = currentpos.x
			nextY = currentpos.y + v.dist
		} else if v.dir == 68 { // down
			nextX = currentpos.x
			nextY = currentpos.y - v.dist
		} else if v.dir == 76 { // left
			nextX = currentpos.x - v.dist
			nextY = currentpos.y
		} else if v.dir == 82 { // right
			nextX = currentpos.x + v.dist
			nextY = currentpos.y
		} else {
			fmt.Println("Unexpected direction: ", v.dir)
		}

		for q := currentpos.x; q <= nextX; q++ {
			for p := currentpos.y; p <= nextY; p++ {
				pos := fmt.Sprintf("%v,%v", q, p)
				route[pos] = true
			}
		}

		currentpos.x = nextX
		currentpos.y = nextY
	}

	return route
}

func getIntersections(a map[string]bool, b map[string]bool) []coord {
	intersection := make([]coord, max(len(a), len(b)))
	n := 0

	for key := range a {
		if b[key] {
			intersection[n] = stringToCoord(key)
			n++
		}
	}

	return intersection[0:n]
}

func main() {

	var routeA map[string]bool
	var routeB map[string]bool
	var intersections []coord
	var shortestDistance int = math.MaxInt64
	var distance int = -1

	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Bad input.")
		return
	}
	routeA = getRoute(strings.TrimSuffix(line, "\n"))

	line, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Bad input.")
		return
	}
	routeB = getRoute(strings.TrimSuffix(line, "\n"))

	// all the points where routeA and routeB cross
	intersections = getIntersections(routeA, routeB)

	// minimum manhattan distance from 0,0 in the set of intersections
	for _, inter := range intersections {
		distance = inter.manhattanDistance()
		if distance < shortestDistance {
			shortestDistance = distance
		}
	}

	fmt.Println(shortestDistance)

}
