/*
https://adventofcode.com/2019/day/18

Objective: Determine the minimum number of steps needed to pick up
all they keys (lowercase letters) in the maze.
Doors (uppercase letters) cannot be traversed until the corresponding
key has been picked up.

In part 2, the maze is transformed by mutating the central pattern as follows:
...    @#@
.@. => ###
...    @#@
This splits the maze into four disjoint subsections.
In order to pick up all the keys, we deploy one robot to each subsection.
The robots move sequentially (i.e., not in parallel),
but otherwise the problem remains unchanged.

About the solution:
This solution employs Dijkstra's Shortest Path algorithm and dynamic programming.

Step 1. Convert the puzzle input into an undirected graph representation.
Step 2. Precalculate the shortest path from each node to every key in the maze.
Step 3. Recursively calculate the sequence of keys that results in
		the shortest path starting from the current puzzle state.

Each path that is calculated is saved in a memoization dictionary whose
keys represent a distinct state of the puzzle.
*/
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"unicode"

	"github.com/dliberat/algoimpl/go/graph"
)

// WALL cannot be traversed
var WALL = []rune("#")[0]

// PERIOD empty space on map
var PERIOD = []rune(".")[0]

// START entry point of map
var START = []rune("@")[0]

// nodes contained in the graph
var nodes = make(map[coord]graph.Node, 0)

// allPaths contains, for each coord, the shortest path
// to each key in the maze.
var allPaths = make(map[coord][]graph.Path, 0)

// memo records the smallest amount of steps required
// to collect all the keys in the maze given an initial state
var memo memoDict

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

type coord struct {
	x   int
	y   int
	val rune
}

func (c *coord) toString() string {
	return fmt.Sprintf("(%d,%d,%s)", c.x, c.y, string(c.val))
}

func (c *coord) equals(other coord) bool {
	return (*c).x == other.x && (*c).y == other.y && (*c).val == other.val
}

type pathSearchResult struct {
	foundPath           bool
	badTargetNode       bool
	keyAlongPath        string
	lockedDoorAlongPath string
}

type memoDict struct {
	dict map[string]int
}

// getHashKey generates the unique identifier for the given combination of elements
func (c *memoDict) getHashKey(keyset []coord, botPositions []coord, activeBot int) string {
	keyStrings := make([]string, 0)
	for _, key := range keyset {
		keyStrings = append(keyStrings, string(key.val))
	}
	sort.Strings(keyStrings)

	botStrings := make([]string, len(botPositions))
	for i := range botPositions {
		botStrings[i] = botPositions[i].toString()
	}

	hashKey := strings.Join(botStrings, "") + "-" + fmt.Sprintf("%d", activeBot) + "-" + strings.Join(keyStrings, "")
	return hashKey
}

// getValue returns the memoized value for the shortest path that picks up
// all keys starting from the given coord and holding the given keyset
// A return value of 0 indicates nothing has been memoized yet for
// the given parameters.
func (c *memoDict) getValue(keyset []coord, botPositions []coord, activeBot int) int {
	hashKey := c.getHashKey(keyset, botPositions, activeBot)
	return c.dict[hashKey]
}

func (c *memoDict) setValue(keyset []coord, botPositions []coord, activeBot, value int) {
	hashKey := c.getHashKey(keyset, botPositions, activeBot)
	c.dict[hashKey] = value
}

func haveKey(key rune, keys []coord) bool {
	for _, e := range keys {
		if e.val == key {
			return true
		}
	}
	return false
}

// map2runes converts string data from input file to a 2-D array of runes
func map2runes(mapData string) [][]rune {
	m := make([][]rune, 0)
	lines := strings.Split(mapData, "\n")
	for _, line := range lines {
		m = append(m, []rune(line))
	}
	return m
}

// findStart identifies the coordinates containing the "@" start locations for the maze
func findStart(mapData *[][]rune) []coord {
	maxX := len((*mapData)[0]) - 1
	maxY := len(*mapData) - 1
	startLocations := make([]coord, 0)

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			r := (*mapData)[y][x]
			if r == START {
				start := coord{x: x, y: y, val: r}
				startLocations = append(startLocations, start)
			}
		}
	}
	if len(startLocations) == 0 {
		panic("No start location")
	}
	return startLocations
}

// isKey returns true if the provided rune represents a key in the maze
func isKey(r rune) bool {
	return r != START && r != PERIOD && r != WALL && unicode.IsLower(r)
}

// isDoor returns true if the provided rune represents a door in the maze
func isDoor(r rune) bool {
	return r != PERIOD && r != START && unicode.IsUpper(r)
}

// isLockedDoor returns true if the coordinate c is a door and
// its corresponding key is not present in the currentKeys set.
func isLockedDoor(c coord, currentKeys []coord) bool {
	return isDoor(c.val) && !haveKey(unicode.ToLower(c.val), currentKeys)
}

// findKeys returns the coordinates of all keys in the maze
func findKeys(mapData *[][]rune) []coord {
	maxX := len((*mapData)[0]) - 1
	maxY := len(*mapData) - 1
	keys := make([]coord, 0)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			r := (*mapData)[y][x]
			if isKey(r) {
				keys = append(keys, coord{x: x, y: y, val: r})
			}
		}
	}
	return keys
}

func getNeighbors(mapData [][]rune, c coord) []coord {
	maxX := len(mapData[0]) - 1
	maxY := len(mapData) - 1
	neighbors := make([]coord, 0)
	if c.y > 0 && mapData[c.y-1][c.x] != WALL {
		neighbors = append(neighbors, coord{x: c.x, y: c.y - 1, val: mapData[c.y-1][c.x]})
	}
	if c.y < maxY && mapData[c.y+1][c.x] != WALL {
		neighbors = append(neighbors, coord{x: c.x, y: c.y + 1, val: mapData[c.y+1][c.x]})
	}
	if c.x > 0 && mapData[c.y][c.x-1] != WALL {
		neighbors = append(neighbors, coord{x: c.x - 1, y: c.y, val: mapData[c.y][c.x-1]})
	}
	if c.x < maxX && mapData[c.y][c.x+1] != WALL {
		neighbors = append(neighbors, coord{x: c.x + 1, y: c.y, val: mapData[c.y][c.x+1]})
	}
	return neighbors
}

func makeGraph(mapData [][]rune, nodes *map[coord]graph.Node) *graph.Graph {
	G := graph.New(graph.Undirected)

	for y := 0; y < len(mapData); y++ {
		for x := 0; x < len(mapData[y]); x++ {
			currentNode := coord{x: x, y: y, val: mapData[y][x]}
			if currentNode.val == WALL {
				continue
			}

			_, exists := (*nodes)[currentNode]
			if !exists {
				(*nodes)[currentNode] = G.MakeNode()
			}

			neighbors := getNeighbors(mapData, currentNode)
			for _, neighbor := range neighbors {

				if mapData[neighbor.y][neighbor.x] == WALL {
					continue
				}

				_, exists = (*nodes)[neighbor]
				if !exists {
					(*nodes)[neighbor] = G.MakeNode()
				}

				G.MakeEdgeWeight((*nodes)[currentNode], (*nodes)[neighbor], 1)
			}
		}
	}

	for key, node := range *nodes {
		*node.Value = key
	}

	return G
}

func isPathToNonPeriodNode(path graph.Path) bool {
	if len(path.Path) < 1 {
		return false
	}
	e := path.Path[len(path.Path)-1]
	targetNodeRune := (*e.End.Value).(coord).val
	return targetNodeRune != PERIOD
}

func isPathToKey(path graph.Path) bool {
	if len(path.Path) < 1 {
		return false
	}
	e := path.Path[len(path.Path)-1]
	targetNodeRune := (*e.End.Value).(coord).val
	return isKey(targetNodeRune)
}

// pathContainsOnlyPeriods returns true if all the nodes
// in between the start and end points of the path are PERIODs.
// The starting and ending nodes are not checked.
// Example: The path @....a would return true
// Example: The path b....D would return true
// Example: The path X..... would return true
func pathContainsOnlyPeriods(path graph.Path) bool {
	for i, edge := range path.Path {
		// First node should never be a period
		if i == 0 {
			continue
		}
		if (*edge.Start.Value).(coord).val != PERIOD {
			return false
		}
	}
	return true
}

// pruneGraph removes all PERIOD nodes and sets edge weights between
// non-PERIOD nodes accordingly.
func pruneGraph(g *graph.Graph, oldNodes map[coord]graph.Node) (*graph.Graph, map[coord]graph.Node) {
	newGraph := graph.New(graph.Undirected)
	newNodes := make(map[coord]graph.Node, 0)

	// first pass: Copy nodes into the new graph.
	// Nodes need to be in the graph before we can create edges
	for _, from := range oldNodes {
		currentCoord := (*from.Value).(coord)
		if currentCoord.val == PERIOD {
			continue
		}
		newNodes[currentCoord] = newGraph.MakeNode()
		*newNodes[currentCoord].Value = currentCoord
	}

	// Second pass: create edges between non-period nodes
	for _, node := range oldNodes {
		currentCoord := (*node.Value).(coord)
		if currentCoord.val == PERIOD {
			continue
		}

		from := newNodes[currentCoord]

		paths := g.DijkstraSearch(node)
		for _, path := range paths {
			if path.Weight == 0 {
				continue // Ignore paths to self
			}
			if len(path.Path) == 0 {
				continue // Ignore disojint sections
			}
			if !isPathToNonPeriodNode(path) {
				continue // Only consider paths that end in something other than .
			}
			if !pathContainsOnlyPeriods(path) {
				continue // cannot skip over keys or doors
			}

			lastEdge := path.Path[len(path.Path)-1]
			targetCoord := (*lastEdge.End.Value).(coord)
			to := newNodes[targetCoord]

			newGraph.MakeEdgeWeight(from, to, path.Weight)
		}
	}

	return newGraph, newNodes
}

// precalcDijkstra calculates the shortest path from each key and start node
// to every other node in the graph.
func precalcDijkstra(g *graph.Graph, keys []coord, startLocations []coord) {
	for _, startLoc := range startLocations {
		allPaths[startLoc] = g.DijkstraSearch(nodes[startLoc])
	}

	for _, key := range keys {
		currentNode := nodes[key]
		allPaths[key] = g.DijkstraSearch(currentNode)
	}
}

// pathToNode determines if the provided path can be traveled to reach the node n.
// There are two conditions that must be met in order for a path to be traversible:
// 1. The end node of the provided path must be the target node n
// 2. If there are any doors along the path, the corresponding key must be present in currentKeys.
//
// A third condition is included as an optimization, since paths that do not
// meet this condition will never result in an optimal route:
// 3. There are no keys along the path which are not present in currentKeys.
//
// The return value provides a boolean value indicating whether the path is
// valid or not, and in the case that the path is invalid, additional information
// explaining why it was determined to be invalid. This additional information
// is useful in debugging.
func pathToNode(path graph.Path, n coord, currentKeys []coord) pathSearchResult {
	e := path.Path[len(path.Path)-1]
	pathEndCoord := (*e.End.Value).(coord)

	if !pathEndCoord.equals(n) {
		return pathSearchResult{foundPath: false, badTargetNode: true}
	}

	for i, e := range path.Path {
		c := (*e.End.Value).(coord)

		// If there is a new key along the path, it's always optimal to pick it up
		// rather than come back for it later.
		if i != len(path.Path)-1 && isKey(c.val) && !haveKey(c.val, currentKeys) {
			return pathSearchResult{foundPath: false, keyAlongPath: string(c.val)}
		}

		if isLockedDoor(c, currentKeys) {
			return pathSearchResult{foundPath: false, lockedDoorAlongPath: string(c.val)}
		}
	}

	return pathSearchResult{foundPath: true}
}

// getPathToNode selects the path that leads to n from the provided slice of paths.
// If node n cannot be reached using the given paths, an error is thrown.
func getPathToNode(paths []graph.Path, n coord, currentKeys []coord) (graph.Path, error) {
	for _, path := range paths {

		if path.Weight == 0 {
			continue // Path to self
		}

		if len(path.Path) == 0 {
			continue // disjoint
		}

		res := pathToNode(path, n, currentKeys)
		if res.foundPath {
			return path, nil
		}
		if res.keyAlongPath != "" {
			return graph.Path{}, errors.New("Path is suboptimal. Better pick up " + res.keyAlongPath + " first.")
		}
		if res.lockedDoorAlongPath != "" {
			return graph.Path{}, errors.New("Cannot reach key. Need to unlock " + res.lockedDoorAlongPath + " first.")
		}
	}

	return graph.Path{}, errors.New("there are no paths that can reach the target node")
}

// rec determines the minimum number of steps required to pick up all the keys
// in the maze starting from the given conditions. keys is the set of all keys available in the maze.
// activeBot is an index into the botLocations slice indicating which was the last robot to move.
func rec(g *graph.Graph, keys []coord, currentKeys []coord, botLocations []coord, activeBot int) int {

	memoizedValue := memo.getValue(currentKeys, botLocations, activeBot)
	if memoizedValue > 0 {
		return memoizedValue
	}

	if len(currentKeys) == len(keys) {
		return 0
	}

	localBest := maxInt

	for newActiveBot, currentPos := range botLocations {

		for _, key := range keys {
			if haveKey(key.val, currentKeys) {
				continue
			}

			// determine the shortest path from the active robot's
			// current position to the new target key
			paths := allPaths[currentPos]
			path, err := getPathToNode(paths, key, currentKeys)
			if err != nil {
				continue // target key is unreachable from the current state
			}

			weightCurrentBranch := maxInt
			newCurrentKeys := append(currentKeys, key)
			newBots := make([]coord, len(botLocations))
			copy(newBots, botLocations)
			newBots[newActiveBot] = key // active bot moves to the target key's location

			res := rec(g, keys, newCurrentKeys, newBots, newActiveBot)
			if res != maxInt {
				weightCurrentBranch = path.Weight + res
			}

			if localBest > weightCurrentBranch {
				localBest = weightCurrentBranch
			}
		}
	}

	memo.setValue(currentKeys, botLocations, activeBot, localBest)

	return localBest
}

// modifyMapForPart2 modifies the central 9 tiles
// in the maze as described in the problem description
func modifyMapForPart2(runes [][]rune) [][]rune {
	start := findStart(&runes)[0]
	runes[start.y-1][start.x-1] = START
	runes[start.y-1][start.x] = WALL
	runes[start.y-1][start.x+1] = START
	runes[start.y][start.x-1] = WALL
	runes[start.y][start.x] = WALL
	runes[start.y][start.x+1] = WALL
	runes[start.y+1][start.x-1] = START
	runes[start.y+1][start.x] = WALL
	runes[start.y+1][start.x+1] = START
	return runes
}

func precalcPathFilter(allPaths *map[coord][]graph.Path) {
	for key, paths := range *allPaths {
		filtered := make([]graph.Path, 0)
		for _, path := range paths {
			if len(path.Path) == 0 {
				continue // disjoint
			}
			if path.Weight == 0 {
				continue // path to self
			}
			if !isPathToKey(path) {
				continue
			}
			filtered = append(filtered, path)
		}
		(*allPaths)[key] = filtered
	}
}

func part1(mapData string) int {
	// reset globals
	memo = memoDict{dict: make(map[string]int, 0)}
	nodes = make(map[coord]graph.Node, 0)
	allPaths = make(map[coord][]graph.Path, 0)

	runes := map2runes(mapData)
	botLocations := findStart(&runes)
	keys := findKeys(&runes)

	tmpNodes := make(map[coord]graph.Node, 0)
	g := makeGraph(runes, &tmpNodes)
	g, nodes = pruneGraph(g, tmpNodes)

	precalcDijkstra(g, keys, botLocations)

	stepCount := rec(g, keys, []coord{}, botLocations, 0)

	return stepCount
}

func part2(mapData string) int {
	// reset globals
	memo = memoDict{dict: make(map[string]int, 0)}
	nodes = make(map[coord]graph.Node, 0)
	allPaths = make(map[coord][]graph.Path, 0)

	runes := map2runes(mapData)
	runes = modifyMapForPart2(runes)
	botLocations := findStart(&runes)
	keys := findKeys(&runes)

	tmpNodes := make(map[coord]graph.Node, 0)
	g := makeGraph(runes, &tmpNodes)
	g, nodes = pruneGraph(g, tmpNodes)

	precalcDijkstra(g, keys, botLocations)
	precalcPathFilter(&allPaths)

	stepCount := rec(g, keys, []coord{}, botLocations, 0)
	return stepCount
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read input file.")
	}

	result := part1(string(data))
	fmt.Println("Part 1: Required", result, "steps.")

	result = part2(string(data))
	fmt.Println("Part 2: Required", result, "steps.")
}

func printMap(mmap [][]rune) {
	for _, lst := range mmap {
		for _, r := range lst {
			fmt.Print(string(r))
		}
		fmt.Println("")
	}
}

// printGraph displays the graph as an adjacency list
func printGraph(G *graph.Graph, nodeList *map[coord]graph.Node) {
	nodes := *nodeList
	for c, node := range nodes {
		val := string(c.val)
		edges := G.Edges(node)
		fmt.Print("(", val, ") -> ")
		for _, edge := range edges {
			i := *(edge.End.Value)
			if i == nil {
				fmt.Print("(nil)")
				continue
			}
			c := i.(coord)
			v := string(c.val)
			fmt.Print("(", v, ": ", edge.Weight, ")")
		}
		fmt.Print("\n")
	}
}

func printPath(p graph.Path) {
	var val rune
	for i, edge := range p.Path {
		if i == 0 {
			val = (*edge.Start.Value).(coord).val
			fmt.Print(string(val), "->")
		}
		val = (*edge.End.Value).(coord).val
		fmt.Print(string(val))
		if i != len(p.Path)-1 {
			fmt.Print("->")
		}
	}
	fmt.Println(" (", p.Weight, ")")
}
