package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type coord struct {
	x float64
	y float64
}

type asteroid struct {
	sightlines map[float64][]coord
	x          float64
	y          float64
}

func (a *asteroid) recordSightline(other asteroid) {
	if a.y == other.y && a.x == other.x {
		return // don't compare against self
	}

	// ys are inverted
	slope := (a.y - other.y) / (other.x - a.x)

	a.sightlines[slope] = append(a.sightlines[slope], coord{x: other.x, y: other.y})
}

func isLeftOf(a coord, b coord) bool {
	return a.x < b.x
}

func isBelow(a coord, b coord) bool {
	return a.y > b.y
}

func countVisible(aster coord, line []coord) int {
	if len(line) < 1 {
		return 0
	} else if len(line) == 1 {
		return 1
	}

	var hasAsteroidsToTheLeft bool
	var hasAsteroidsToTheRight bool
	var hasAsteroidsAbove bool
	var hasAsteroidsBelow bool

	for _, c := range line {
		if isLeftOf(aster, c) {
			hasAsteroidsToTheRight = true
		} else {
			hasAsteroidsToTheLeft = true
		}
		if isBelow(aster, c) {
			hasAsteroidsAbove = true
		} else {
			hasAsteroidsBelow = true
		}

		if (hasAsteroidsToTheLeft && hasAsteroidsToTheRight) || (hasAsteroidsAbove && hasAsteroidsBelow) {
			return 2
		}
	}

	return 1

}

func (a *asteroid) countVisible() int {
	count := 0

	for _, line := range a.sightlines {
		position := coord{x: a.x, y: a.y}
		count += countVisible(position, line)
	}

	return count
}

/*Returns the x and y coordinates of the asteroid with the
best view, and the number of other asteroids it can see.*/
func findBestSight(asterMap []asteroid, display bool) (asteroid, int) {
	var bestAsterIndex int
	var bestSight int = 0

	for i, asteroid := range asterMap {
		for _, other := range asterMap {
			asteroid.recordSightline(other)
		}
		sight := asteroid.countVisible()
		if sight > bestSight {
			bestSight = sight
			bestAsterIndex = i
		}
	}

	if display {
		for _, sightline := range asterMap[bestAsterIndex].sightlines {
			fmt.Println(sightline)
		}
	}

	return asterMap[bestAsterIndex], int(bestSight)
}

func getDirection(prevDir int, prevSlope float64, currentSlope float64) int {
	/*
			Quadrants:
					|
				4 | 1
					|
			---------
					|
				3 | 2
		      |
	*/

	if prevDir == 1 {
		if prevSlope >= 0 && currentSlope < 0 {
			return -1 // was in quadrant 1, moved to quadrant 2
		}
		// if previous slope is <0 then we are in the 4th quadrant
		return 1
	}
	if prevDir == -1 {
		if prevSlope != math.Inf(-1) && prevSlope >= 0 && currentSlope < 0 {
			// was in quadrant 3, moved to quadrant 4
			return 1
		}
		return -1
	}

	return 0
}

func norm2(origin coord, x coord) float64 {
	base := x.x - origin.x
	height := x.y - origin.y
	return math.Sqrt(base*base + height*height)
}

func sortKeys(sightlines map[float64][]coord) []float64 {
	var keys []float64
	for k := range sightlines {
		keys = append(keys, k)
	}

	sort.Float64s(keys)
	return keys
}

func stableDelete(slice []float64, s int) []float64 {
	return append(slice[:s], slice[s+1:]...)
}

func zapOne(dir int, origin coord, key float64, asteroids *map[float64][]coord) (coord, bool) {

	var bestDistance float64 = math.Inf(1)
	var bestIndex int = -1

	for i, target := range (*asteroids)[key] {

		if dir == 1 {
			// only zap asteroids that have a lower y
			// or the same y but a higher x (ie. quadrants 4 and 1)
			if target.y > origin.y {
				continue
			}
			if target.y == origin.y && target.x < origin.x {
				continue
			}
		} else {
			// only zap asteroids that have a higher y
			// or the same y but a lower x (ie. quadrants 2 and 3)
			if target.y < origin.y {
				continue
			}
			if target.y == origin.y && target.x > origin.x {
				continue
			}
		}

		// the candidate asteroid is actually in the direction our laser
		// is pointing to. But we can only zap the nearest one
		distance := norm2(origin, target)
		if distance < bestDistance {
			bestDistance = distance
			bestIndex = i
		}
	}

	zapped := coord{}
	exists := false

	if bestIndex > -1 {
		// remove the nearest asteroid
		last := len((*asteroids)[key]) - 1
		zapped = coord{x: (*asteroids)[key][bestIndex].x, y: (*asteroids)[key][bestIndex].y}
		(*asteroids)[key][bestIndex] = (*asteroids)[key][last]
		(*asteroids)[key] = (*asteroids)[key][:last]
		exists = true
	}

	return zapped, exists
}

func zappingOrder(a asteroid) []coord {

	zapped := make([]coord, 0)
	sortedKeys := sortKeys(a.sightlines) // zap in clockwise order
	pos := sort.SearchFloat64s(sortedKeys, math.Inf(1))

	startIndex := pos // start zapping upward
	if startIndex >= len(sortedKeys) {
		// if there are no asteroids directly above the base station
		startIndex = (startIndex + len(sortedKeys) - 1) % len(sortedKeys)
	}

	prevKey := sortedKeys[startIndex]
	dir := 1
	origin := coord{x: a.x, y: a.y}

	for len(sortedKeys) > 0 {

		k := sortedKeys[startIndex]
		dir = getDirection(dir, prevKey, k)

		target, exists := zapOne(dir, origin, k, &(a.sightlines))
		if exists {
			zapped = append(zapped, target)
			if len(zapped) == 200 {
				return zapped
			}
		}

		if len(a.sightlines[k]) == 0 {

			if len(sortedKeys) == 1 {
				break
			}

			// remove element
			sortedKeys = stableDelete(sortedKeys, startIndex)
		}

		prevKey = k
		startIndex = (startIndex + len(sortedKeys) - 1) % len(sortedKeys)
	}

	return zapped
}

func main() {
	var asterMap []asteroid
	var a asteroid
	var n int
	var display bool = false

	asterMap = makeMap()
	a, n = findBestSight(asterMap, display)
	fmt.Println("Asteroid at", a.x, ",", a.y, " can see ", n, " other asteroids.")
	zapped := zappingOrder(a)
	fmt.Println("The 200th asteroid to be zapped is at", zapped[199].x, ",", zapped[199].y)
}

func makeMap() []asteroid {
	asteroids := make([]string, 0)

	asteroids = append(asteroids, ".##.#.#....#.#.#..##..#.#.")
	asteroids = append(asteroids, "#.##.#..#.####.##....##.#.")
	asteroids = append(asteroids, "###.##.##.#.#...#..###....")
	asteroids = append(asteroids, "####.##..###.#.#...####..#")
	asteroids = append(asteroids, "..#####..#.#.#..#######..#")
	asteroids = append(asteroids, ".###..##..###.####.#######")
	asteroids = append(asteroids, ".##..##.###..##.##.....###")
	asteroids = append(asteroids, "#..#..###..##.#...#..####.")
	asteroids = append(asteroids, "....#.#...##.##....#.#..##")
	asteroids = append(asteroids, "..#.#.###.####..##.###.#.#")
	asteroids = append(asteroids, ".#..##.#####.##.####..#.#.")
	asteroids = append(asteroids, "#..##.#.#.###.#..##.##....")
	asteroids = append(asteroids, "#.#.##.#.##.##......###.#.")
	asteroids = append(asteroids, "#####...###.####..#.##....")
	asteroids = append(asteroids, ".#####.#.#..#.##.#.#...###")
	asteroids = append(asteroids, ".#..#.##.#.#.##.#....###.#")
	asteroids = append(asteroids, ".......###.#....##.....###")
	asteroids = append(asteroids, "#..#####.#..#..##..##.#.##")
	asteroids = append(asteroids, "##.#.###..######.###..#..#")
	asteroids = append(asteroids, "#.#....####.##.###....####")
	asteroids = append(asteroids, "..#.#.#.########.....#.#.#")
	asteroids = append(asteroids, ".##.#.#..#...###.####..##.")
	asteroids = append(asteroids, "##...###....#.##.##..#....")
	asteroids = append(asteroids, "..##.##.##.#######..#...#.")
	asteroids = append(asteroids, ".###..#.#..#...###..###.#.")
	asteroids = append(asteroids, "#..#..#######..#.#..#..#.#")

	asterMap := make([]asteroid, 0)
	for y, line := range asteroids {
		arr := strings.Split(line, "")
		for x, mark := range arr {
			if mark == "#" {
				asterMap = append(asterMap, asteroid{x: float64(x), y: float64(y), sightlines: make(map[float64][]coord)})
			}
		}
	}
	return asterMap
}
