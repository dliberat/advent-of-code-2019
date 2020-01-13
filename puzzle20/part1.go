package main

import (
	"fmt"
)

type tile struct {
	x       int
	y       int
	visited bool
	name    string
}

var allTiles []tile
var edges map[*tile][]*tile

/*Part1 runs the code for part 1 of this puzzle*/
func Part1(fname string) {
	allTiles = make([]tile, 0)
	edges = make(map[*tile][]*tile)

	maze := GetInput(fname)
	t := parseInput(maze)

	fmt.Println("Parsed maze! ", t)

	depth := bfs(t)
	fmt.Println("[Part 1] The maze exit can be reached in", depth, "steps")
}

func bfs(t *tile) int {

	depth := 0
	queue := make([]*tile, 1)
	queue[0] = t
	nextLayer := make([]*tile, 0)

	for len(queue) > 0 {
		currentTile := queue[0]

		if currentTile.name == "ZZ" {
			return depth
		}

		neighbors := edges[currentTile]
		for _, neighbor := range neighbors {
			if !neighbor.visited {
				nextLayer = append(nextLayer, neighbor)
			}
		}

		currentTile.visited = true

		if len(queue) == 1 {
			queue = nextLayer
			nextLayer = nil
			depth++
		} else {
			queue = queue[1:]
		}
	}

	panic("No exit found!")
}

func parseInput(input [][]string) *tile {

	var entrance *tile

	// first, find all the tiles
	for row := range input {
		for col := range input[row] {
			if input[row][col] == "." {
				allTiles = append(allTiles, tile{x: col, y: row})
			}
		}
	}

	for i := range allTiles {
		t := &allTiles[i]

		right := getRight(t, &input)
		left := getLeft(t, &input)
		up := getUp(t, &input)
		down := getDown(t, &input)

		// names
		if len(right) == 2 {
			t.name = right
			if right != "AA" && right != "ZZ" {
				neighbor := getNamedTile(right)
				addEdge(t, neighbor)
			}
		} else if len(left) == 2 {
			t.name = left
			if left != "AA" && left != "ZZ" {
				neighbor := getNamedTile(left)
				addEdge(t, neighbor)
			}
		} else if len(up) == 2 {
			t.name = up
			if up != "AA" && up != "ZZ" {
				neighbor := getNamedTile(up)
				addEdge(t, neighbor)
			}
		} else if len(down) == 2 {
			t.name = down
			if down != "AA" && down != "ZZ" {
				neighbor := getNamedTile(down)
				addEdge(t, neighbor)
			}
		}

		if t.name == "AA" {
			entrance = t
		}

		if right == "." {
			neighbor := getTile(t.x+1, t.y)
			addEdge(t, neighbor)
		}
		if left == "." {
			neighbor := getTile(t.x-1, t.y)
			addEdge(t, neighbor)
		}
		if up == "." {
			neighbor := getTile(t.x, t.y-1)
			addEdge(t, neighbor)
		}
		if down == "." {
			neighbor := getTile(t.x, t.y+1)
			addEdge(t, neighbor)
		}
	}

	return entrance
}

func addEdge(from *tile, to *tile) {
	if from == nil || to == nil {
		return
	}

	if edges[from] == nil {
		edges[from] = make([]*tile, 0)
	}
	if edges[to] == nil {
		edges[to] = make([]*tile, 0)
	}

	isRegistered := false
	for _, neighbor := range edges[from] {
		if neighbor == to {
			isRegistered = true
			break
		}
	}
	if !isRegistered {
		edges[from] = append(edges[from], to)
	}

	isRegistered = false
	for _, neighbor := range edges[to] {
		if neighbor == from {
			isRegistered = true
			break
		}
	}
	if !isRegistered {
		edges[to] = append(edges[to], from)

	}
}

func getTile(x int, y int) *tile {
	for i := range allTiles {
		if allTiles[i].x == x && allTiles[i].y == y {
			return &allTiles[i]
		}
	}
	panic("Tried to get pointer to nonexistent tile")
}
func getNamedTile(name string) *tile {
	for i := range allTiles {
		if allTiles[i].name == name {
			return &allTiles[i]
		}
	}
	return nil
}

func getRight(t *tile, maze *[][]string) string {
	row := (*maze)[t.y]

	if t.x >= len(row)-1 {
		// we're already at the rightmost edge
		return "#"
	}

	char := row[t.x+1]
	if char == "." || char == "#" {
		return char
	}

	// otherwise, it's a two letter symbol indicating a portal
	return char + row[t.x+2]
}
func getLeft(t *tile, maze *[][]string) string {
	row := (*maze)[t.y]

	if t.x == 0 {
		// we're already at the leftmost edge
		return "#"
	}

	char := row[t.x-1]
	if char == "." || char == "#" {
		return char
	}

	// otherwise, it's a two letter symbol indicating a portal
	return row[t.x-2] + char
}
func getUp(t *tile, maze *[][]string) string {
	if t.y == 0 {
		return "#"
	}

	char := (*maze)[t.y-1][t.x]
	if char == "." || char == "#" {
		return char
	}

	// otherwise, it's a two letter symbol indicating a portal
	return (*maze)[t.y-2][t.x] + char
}
func getDown(t *tile, maze *[][]string) string {
	if t.y >= len(*maze)-1 {
		return "#"
	}

	char := (*maze)[t.y+1][t.x]
	if char == "." || char == "#" {
		return char
	}

	// otherwise, it's a two letter symbol indicating a portal
	return char + (*maze)[t.y+2][t.x]
}
