package main

import (
	"fmt"
	"math"
	"sort"
)

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

/**The AA tile*/
var entrance coord

/**The ZZ tile*/
var exit coord

/**Splits the grid left to right*/
var verticalMedian int

/**Splits the grid approximately top to bottom
(inexact, but good enough to tell which side a portal is on)*/
var horizontalMedian int

/**Locations of all tiles connected to portals on the inside of the map*/
var insidePortals map[string]coord = make(map[string]coord)

/**Locations of all tiles connected to portals on the outside of the map*/
var outsidePortals map[string]coord = make(map[string]coord)

type coord struct {
	x int
	y int
	z int
}

type location struct {
	x                int
	y                int
	depth            int
	dijkstraDistance int
	dijkstraParent   coord
}

/*Euclidian distance from the goal node,
with extra weight assigned to the recursive depth*/
func (loc *location) distance(goal coord) float64 {
	zWeight := abs(loc.depth) * 100
	x := float64(goal.x - loc.x)
	y := float64(goal.y - loc.y)
	return math.Sqrt(x*x+y*y) + float64(zWeight)
}

/**Returns the location corresponding to moving upward from the
current location. A location with x and y equal to -1 indicates a wall.*/
func (loc *location) up(grid *[][]string) location {
	return makeMove(loc.x, loc.y, 0, -1, loc.depth, grid)
}

/**Returns the location corresponding to moving downward from the
current location. A location with x and y equal to -1 indicates a wall.*/
func (loc *location) down(grid *[][]string) location {
	return makeMove(loc.x, loc.y, 0, 1, loc.depth, grid)
}

/**Returns the location corresponding to moving leftward from the
current location. A location with x and y equal to -1 indicates a wall.*/
func (loc *location) left(grid *[][]string) location {
	return makeMove(loc.x, loc.y, -1, 0, loc.depth, grid)
}

/**Returns the location corresponding to moving rightward from the
current location. A location with x and y equal to -1 indicates a wall.*/
func (loc *location) right(grid *[][]string) location {
	return makeMove(loc.x, loc.y, 1, 0, loc.depth, grid)
	// row := (*grid)[loc.y]

	// // we're already at the rightmost edge
	// if loc.x >= len(row)-1 {
	// 	return location{x: -1, y: -1}
	// }

	// // we're somewhere inside the map
	// char := row[loc.x+1]
	// if isUppercase(char[0]) {
	// 	// the current tile is a named tile, and most likely a portal

	// 	portalCode := char + row[loc.x+2]
	// 	if portalCode == "AA" {
	// 		// the entrance never functions as a portal
	// 		return location{x: -1, y: -1}
	// 	} else if portalCode == "ZZ" && loc.depth == 0 {
	// 		fmt.Println("Whoa nelly, we've reached the exit!")
	// 		return location{x: -1, y: -1}
	// 	} else if portalCode == "ZZ" && loc.depth != 0 {
	// 		// at any depth other than 0, ZZ functions as a wall
	// 		return location{x: -1, y: -1}
	// 	} else if loc.x > verticalMedian {
	// 		// travel through an outside portal, arrive at the equivalent location via an inside portal
	// 		wormhole := insidePortals[portalCode]
	// 		return location{x: wormhole.x, y: wormhole.y, depth: loc.depth - 1}
	// 	} else {
	// 		// travel through an inside portal, arrive at the equivalent location via an outside portal
	// 		wormhole := outsidePortals[portalCode]
	// 		return location{x: wormhole.x, y: wormhole.y, depth: loc.depth + 1}
	// 	}

	// } else if char == "." {
	// 	// tile on the same  depth
	// 	return location{x: loc.x + 1, y: loc.y, depth: loc.depth}
	// }

	// // wall
	// return location{x: -1, y: -1}
}

/**Returns a slice of locations corresponding to all the possible
places that can be reached from the current location*/
func (loc *location) getMoves(grid *[][]string) []location {
	retval := make([]location, 0)
	u := loc.up(grid)
	d := loc.down(grid)
	r := loc.right(grid)
	l := loc.left(grid)

	if u.x != -1 && u.y != -1 {
		retval = append(retval, u)
	}
	if d.x != -1 && d.y != -1 {
		retval = append(retval, d)
	}
	if r.x != -1 && r.y != -1 {
		retval = append(retval, r)
	}
	if l.x != -1 && l.y != -1 {
		retval = append(retval, l)
	}
	for x := range retval {
		// all edges have weight 1
		retval[x].dijkstraDistance = loc.dijkstraDistance + 1
	}
	return retval
}

func (loc *location) getCoord() coord {
	return coord{x: loc.x, y: loc.y, z: loc.depth}
}

/*Part2 solves part 2 of this puzzle*/
func Part2(fname string) {
	grid := GetInput(fname)

	// split the grid up into 4 quarters
	yEnd := len(grid) - 1
	xEnd := len(grid[0]) - 1
	verticalMedian = xEnd / 2
	horizontalMedian = yEnd / 2

	parsePortals(&grid)

	dijkstra(&grid)
}

func isUppercase(b byte) bool {
	return b >= 65 && b <= 90
}

/**Find all the portal locations on the grid provided as text input to the puzzle*/
func parsePortals(grid *[][]string) {

	yEnd := len(*grid) - 1
	xEnd := len((*grid)[0]) - 1

	for y := 0; y < yEnd; y++ {
		for x := 0; x < xEnd; x++ {

			c := (*grid)[y][x]
			if isUppercase(c[0]) {

				r := (*grid)[y][x+1]
				d := (*grid)[y+1][x]

				// look right
				if isUppercase(r[0]) {
					portalCode := c + r
					if x == 0 { // outside portal on the left side

						position := coord{x: x + 2, y: y}
						if portalCode == "AA" {
							entrance = position
						} else if portalCode == "ZZ" {
							exit = position
						} else {
							outsidePortals[portalCode] = position
						}

					} else if x+1 >= xEnd { // outside portal on the right side
						position := coord{x: x - 1, y: y}
						if portalCode == "AA" {
							entrance = position
						} else if portalCode == "ZZ" {
							exit = position
						} else {
							outsidePortals[portalCode] = position
						}

					} else if x < verticalMedian { // inside portal on the left side
						insidePortals[portalCode] = coord{x: x - 1, y: y}
					} else { // inside portal on the right side
						insidePortals[portalCode] = coord{x: x + 2, y: y}
					}

					// look down
				} else if isUppercase(d[0]) {
					portalCode := c + d

					if y == 0 { //outside portal on the top of the map
						position := coord{x: x, y: y + 2}
						if portalCode == "AA" {
							entrance = position
						} else if portalCode == "ZZ" {
							exit = position
						} else {
							outsidePortals[portalCode] = position
						}
					} else if y+1 >= yEnd { // outside portal on the bottom of the map
						position := coord{x: x, y: y - 1}
						if portalCode == "AA" {
							entrance = position
						} else if portalCode == "ZZ" {
							exit = position
						} else {
							outsidePortals[portalCode] = position
						}
					} else if y < horizontalMedian { // inside portal on the top half
						insidePortals[portalCode] = coord{x: x, y: y - 1}
					} else { // inside portal on the bottom half
						insidePortals[portalCode] = coord{x: x, y: y + 2}
					}

				}
				// otherwise, we've already found this portal, so move on

			}
		}
	}

}

func isInsidePortal(x, y, xDiff, yDiff int) bool {
	if xDiff > 0 {
		return x < verticalMedian
	} else if xDiff < 0 {
		return x > verticalMedian
	} else if yDiff > 0 {
		return y < horizontalMedian
	}
	return y > horizontalMedian
}

func makeMove(x int, y int, xDiff int, yDiff int, currentDepth int, grid *[][]string) location {
	g := *grid
	char := g[y+yDiff][x+xDiff]
	var portalCode string

	if isUppercase(char[0]) {
		// the current tile is a named tile, and most likely a portal

		if yDiff < 0 || xDiff < 0 {
			// if looking upward or leftward, we need to look at the further letter first
			portalCode = g[y+2*yDiff][x+2*xDiff] + char
		} else {
			portalCode = char + g[y+2*yDiff][x+2*xDiff]
		}

		isOutsidePortal := !isInsidePortal(x, y, xDiff, yDiff)

		if portalCode == "AA" {
			// the entrance never functions as a portal
			return location{x: -1, y: -1}
		} else if portalCode == "ZZ" && currentDepth == 0 {
			fmt.Println("Whoa nelly, we've reached the exit!")
			return location{x: -1, y: -1}
		} else if portalCode == "ZZ" && currentDepth != 0 {
			// at any depth other than 0, ZZ functions as a wall
			return location{x: -1, y: -1}
		} else if isOutsidePortal && currentDepth != 0 {
			// travel through an outside portal, arrive at the equivalent location via an inside portal
			wormhole := insidePortals[portalCode]
			return location{x: wormhole.x, y: wormhole.y, depth: currentDepth - 1}
		} else if isOutsidePortal && currentDepth == 0 {
			// when at the outermost level, only the outer labels AA and ZZ function
			// (as the start and end, respectively); all other outer labeled tiles are effectively walls.
			return location{x: -1, y: -1}
		} else {
			// travel through an inside portal, arrive at the equivalent location via an outside portal
			wormhole := outsidePortals[portalCode]
			return location{x: wormhole.x, y: wormhole.y, depth: currentDepth + 1}
		}

	} else if char == "." {
		// tile on the same  depth
		return location{x: x + xDiff, y: y + yDiff, depth: currentDepth}
	}

	// wall
	return location{x: -1, y: -1}
}

func inSlice(l location, s *[]location) int {
	for i, v := range *s {
		if v.x == l.x && v.y == l.y && v.depth == l.depth {
			return i
		}
	}
	return -1
}

func coordInLocSlice(c coord, s *[]location) int {
	for i, v := range *s {
		if v.x == c.x && v.y == c.y && v.depth == c.z {
			return i
		}
	}
	return -1
}

func printGrid(grid *[][]string, currentPos coord) {
	fmt.Println("\n      ====== DEPTH", currentPos.z, "=======")
	for y := range *grid {
		for x := range (*grid)[0] {
			if y == currentPos.y && x == currentPos.x {
				fmt.Print("O")
			} else {
				fmt.Print((*grid)[y][x])
			}
		}
		fmt.Print("\n")
	}
}

func tracePath(end location, visited *[]location, grid *[][]string) {

	path := make([]location, 1)
	path[0] = end

	currentLoc := end
	parent := end.dijkstraParent

	for parent.x != 0 && parent.y != 0 {
		index := coordInLocSlice(parent, visited)
		if index > -1 {
			currentLoc = (*visited)[index]
			path = append(path, currentLoc)
			parent = currentLoc.dijkstraParent
		} else {
			fmt.Println("Missing parent?")
			break
		}
	}

	steps := -1
	for i := len(path) - 1; i > -1; i-- {
		// printGrid(grid, path[i].getCoord())
		steps++
	}
	fmt.Println(steps)
}

func dijkstra(grid *[][]string) {

	var visited []location = make([]location, 0)

	// initialize the queue with only the start position
	var queue []location = make([]location, 1)
	queue[0] = location{x: entrance.x, y: entrance.y, depth: 0, dijkstraDistance: 0}

	for true {
		// pop the first item off of the queue
		currentLocation := queue[0]
		currentCoord := currentLocation.getCoord()
		queue = queue[1:]

		// termination condition. Wait until the exit is at the top of the queue
		// to ensure that there are no shorter paths left to explore in the queue
		if currentCoord == exit {
			fmt.Println("Reached exit coordinate.")
			tracePath(currentLocation, &visited, grid)
			return
		}

		neighbors := currentLocation.getMoves(grid)
		for _, neighbor := range neighbors {

			if inSlice(neighbor, &visited) >= 0 {
				continue
			}

			index := inSlice(neighbor, &queue)
			if index == -1 { // this is the first time visiting this node
				neighbor.dijkstraParent = currentCoord
				queue = append(queue, neighbor)
			} else {
				if queue[index].dijkstraDistance > neighbor.dijkstraDistance {
					queue[index].dijkstraDistance = neighbor.dijkstraDistance
					queue[index].dijkstraParent = currentCoord
				}
			}
		}

		sort.Slice(queue, func(i, j int) bool { return queue[i].distance(exit) < queue[j].distance(exit) })

		visited = append(visited, currentLocation)
	}

}
