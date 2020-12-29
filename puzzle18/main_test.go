package main

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/dliberat/algoimpl/go/graph"
)

func TestCoordToString(t *testing.T) {
	c := coord{x: 1, y: 2, val: []rune("a")[0]}
	s := c.toString()
	if s != "(1,2,a)" {
		t.Errorf("Incorrect string output for (1, 2, a): %s", s)
	}
}

func TestIsKey(t *testing.T) {
	keys := "abcdefghijklmnopqrstuvwxyz"
	for _, k := range keys {
		result := isKey(k)
		if !result {
			t.Errorf("'%c' should be a key.", k)
		}
	}

	nonKeys := "ABCDEFGHIJKLMNOPQRSTUVWXYZ#@!1234567890-_'[]{}|?()*&"
	for _, k := range nonKeys {
		result := isKey(k)
		if result {
			t.Errorf("'%c' should NOT be a key.", k)
		}
	}
}

func TestHaveKey(t *testing.T) {
	runeA := []rune("a")[0]
	runeB := []rune("b")[0]
	runeC := []rune("c")[0]
	runeD := []rune("d")[0]
	runeE := []rune("e")[0]
	keyA := coord{x: 0, y: 0, val: runeA}
	keyB := coord{x: 1, y: 1, val: runeB}
	keyC := coord{x: 2, y: 2, val: runeC}
	keys := []coord{keyA, keyB, keyC}

	result := haveKey(runeA, keys)
	if !result {
		t.Errorf("Key A is in the keyset.")
	}
	result = haveKey(runeB, keys)
	if !result {
		t.Errorf("Key B is in the keyset.")
	}
	result = haveKey(runeC, keys)
	if !result {
		t.Errorf("Key C is in the keyset.")
	}
	result = haveKey(runeD, keys)
	if result {
		t.Errorf("Key D is NOT in the keyset.")
	}
	result = haveKey(runeE, keys)
	if result {
		t.Errorf("Key E is NOT in the keyset.")
	}
}

func TestFindStartSingle(t *testing.T) {
	maze := `#########
#b.A.@.a#
#########`

	runes := map2runes(maze)
	start := findStart(&runes)[0]
	if start.x != 5 || start.y != 1 || string(start.val) != "@" {
		t.Error("Failed to locate start coord.")
	}
}

func TestFindStartMultiple(t *testing.T) {
	maze := `#######
#a.#Cd#
##@#@##
#######
##@#@##
#cB#Ab#
#######`

	runes := map2runes(maze)
	startLocations := findStart(&runes)
	if len(startLocations) != 4 {
		t.Errorf("Should have found 4 start locations. Only found %d", len(startLocations))
	}

	foundTopLeft := false
	foundTopRight := false
	foundBottomLeft := false
	foundBottomRight := false
	for _, c := range startLocations {
		if c.x == 2 && c.y == 2 {
			foundTopLeft = true
		}
		if c.x == 4 && c.y == 2 {
			foundTopRight = true
		}
		if c.x == 2 && c.y == 4 {
			foundBottomLeft = true
		}
		if c.x == 4 && c.y == 4 {
			foundBottomRight = true
		}
	}
	if !(foundTopLeft && foundTopRight && foundBottomLeft && foundBottomRight) {
		locs := make([]string, 0)
		for _, c := range startLocations {
			locs = append(locs, c.toString())
		}
		t.Errorf("Failed to find all 4 locations. %v", locs)
	}
}

func TestFindKeys(t *testing.T) {
	maze := `#########
####d####
####.####
#b.A.@.a#
###.#####
###c#####
#########`

	runes := map2runes(maze)
	keys := findKeys(&runes)
	if len(keys) != 4 {
		t.Errorf("Expected to find 4 keys, found %d", len(keys))
	}
	foundA := false
	foundB := false
	foundC := false
	foundD := false
	for _, key := range keys {
		isA := key.x == 7 && key.y == 3 && string(key.val) == "a"
		isB := key.x == 1 && key.y == 3 && string(key.val) == "b"
		isC := key.x == 3 && key.y == 5 && string(key.val) == "c"
		isD := key.x == 4 && key.y == 1 && string(key.val) == "d"
		if isA {
			foundA = true
		}
		if isB {
			foundB = true
		}
		if isC {
			foundC = true
		}
		if isD {
			foundD = true
		}
	}
	if !(foundA && foundB && foundC && foundD) {
		t.Errorf("Failed to find all keys in maze. %v", keys)
	}
}

func TestGetNeighborsLeftWall(t *testing.T) {
	maze := `#####
#.a.#
#b.c#
#.d.#
#####
`
	runes := map2runes(maze)
	b := coord{x: 1, y: 2, val: []rune("b")[0]}
	neighbors := getNeighbors(runes, b)

	if len(neighbors) != 3 {
		t.Errorf("Expected 3 neighbors, found %d. %v", len(neighbors), neighbors)
	}

	foundUp := false
	foundRight := false
	foundDown := false
	for _, n := range neighbors {
		isUp := n.x == b.x && n.y == b.y-1 && string(n.val) == "."
		isRight := n.x == b.x+1 && n.y == b.y && string(n.val) == "."
		isDown := n.x == b.x && n.y == b.y+1 && string(n.val) == "."
		if isUp {
			foundUp = true
		}
		if isRight {
			foundRight = true
		}
		if isDown {
			foundDown = true
		}
	}
	if !(foundUp && foundRight && foundDown) {
		t.Errorf("Failed to find all neighbors. %v", neighbors)
	}
}

func TestGetNeighborsRightWall(t *testing.T) {
	maze := `#####
#.a.#
#b.c#
#.d.#
#####
`
	runes := map2runes(maze)
	c := coord{x: 3, y: 2, val: []rune("c")[0]}
	neighbors := getNeighbors(runes, c)

	if len(neighbors) != 3 {
		t.Errorf("Expected 3 neighbors, found %d. %v", len(neighbors), neighbors)
	}

	foundUp := false
	foundLeft := false
	foundDown := false
	for _, n := range neighbors {
		isUp := n.x == c.x && n.y == c.y-1 && string(n.val) == "."
		isLeft := n.x == c.x-1 && n.y == c.y && string(n.val) == "."
		isDown := n.x == c.x && n.y == c.y+1 && string(n.val) == "."
		if isUp {
			foundUp = true
		}
		if isLeft {
			foundLeft = true
		}
		if isDown {
			foundDown = true
		}
	}
	if !(foundUp && foundLeft && foundDown) {
		t.Errorf("Failed to find all neighbors. %v", neighbors)
	}
}

func TestGetNeighborsTopWall(t *testing.T) {
	maze := `#####
#.a.#
#b.c#
#.d.#
#####
`
	runes := map2runes(maze)
	a := coord{x: 2, y: 1, val: []rune("a")[0]}
	neighbors := getNeighbors(runes, a)

	if len(neighbors) != 3 {
		t.Errorf("Expected 3 neighbors, found %d. %v", len(neighbors), neighbors)
	}

	foundLeft := false
	foundRight := false
	foundDown := false
	for _, n := range neighbors {
		isLeft := n.x == a.x-1 && n.y == a.y && string(n.val) == "."
		isRight := n.x == a.x+1 && n.y == a.y && string(n.val) == "."
		isDown := n.x == a.x && n.y == a.y+1 && string(n.val) == "."
		if isLeft {
			foundLeft = true
		}
		if isRight {
			foundRight = true
		}
		if isDown {
			foundDown = true
		}
	}
	if !(foundLeft && foundRight && foundDown) {
		t.Errorf("Failed to find all neighbors. %v", neighbors)
	}
}

func TestGetNeighborsBottomWall(t *testing.T) {
	maze := `#####
#.a.#
#b.c#
#.d.#
#####
`
	runes := map2runes(maze)
	d := coord{x: 2, y: 3, val: []rune("d")[0]}
	neighbors := getNeighbors(runes, d)

	if len(neighbors) != 3 {
		t.Errorf("Expected 3 neighbors, found %d. %v", len(neighbors), neighbors)
	}

	foundUp := false
	foundLeft := false
	foundRight := false
	for _, n := range neighbors {
		isUp := n.x == d.x && n.y == d.y-1 && string(n.val) == "."
		isLeft := n.x == d.x-1 && n.y == d.y && string(n.val) == "."
		isRight := n.x == d.x+1 && n.y == d.y && string(n.val) == "."
		if isUp {
			foundUp = true
		}
		if isLeft {
			foundLeft = true
		}
		if isRight {
			foundRight = true
		}
	}
	if !(foundUp && foundLeft && foundRight) {
		t.Errorf("Failed to find all neighbors. %v", neighbors)
	}
}

func TestGetNeighborsNoWalls(t *testing.T) {
	maze := `#####
#.a.#
#b.c#
#.d.#
#####
`
	runes := map2runes(maze)
	center := coord{x: 2, y: 2, val: PERIOD}
	neighbors := getNeighbors(runes, center)

	if len(neighbors) != 4 {
		t.Errorf("Expected 4 neighbors, found %d. %v", len(neighbors), neighbors)
	}

	foundUp := false
	foundLeft := false
	foundRight := false
	foundDown := false
	for _, n := range neighbors {
		isUp := n.x == 2 && n.y == 1 && string(n.val) == "a"
		isLeft := n.x == 1 && n.y == 2 && string(n.val) == "b"
		isRight := n.x == 3 && n.y == 2 && string(n.val) == "c"
		isDown := n.x == 2 && n.y == 3 && string(n.val) == "d"
		if isUp {
			foundUp = true
		}
		if isLeft {
			foundLeft = true
		}
		if isRight {
			foundRight = true
		}
		if isDown {
			foundDown = true
		}
	}
	if !(foundUp && foundLeft && foundRight && foundDown) {
		t.Errorf("Failed to find all neighbors. %v", neighbors)
	}
}

func isIn(n graph.Node, seq []graph.Node) bool {
	for _, node := range seq {
		if n == node {
			return true
		}
	}
	return false
}

func TestMakeGraph(t *testing.T) {
	maze := `#####
#a.b#
##c##
#####
`
	var a graph.Node
	var b graph.Node
	var c graph.Node
	var dot graph.Node

	runes := map2runes(maze)
	nodes := make(map[coord]graph.Node, 0)
	g := makeGraph(runes, &nodes)

	for _, node := range nodes {
		tmp := (*node.Value).(coord)
		val := string(tmp.val)
		if val == "a" {
			a = node
		}
		if val == "b" {
			b = node
		}
		if val == "c" {
			c = node
		}
		if val == "." {
			dot = node
		}
	}

	aNeighbors := g.Neighbors(a)
	if len(aNeighbors) != 1 {
		t.Error("'a' should have exactly 1 neighbor")
	}
	if !isIn(dot, aNeighbors) {
		t.Errorf("'a' should have '.' as a neighbor")
	}

	bNeighbors := g.Neighbors(b)
	if len(bNeighbors) != 1 {
		t.Error("'b' should have exactly 1 neighbor")
	}
	if !isIn(dot, bNeighbors) {
		t.Errorf("'b' should have '.' as b neighbor")
	}

	cNeighbors := g.Neighbors(c)
	if len(cNeighbors) != 1 {
		t.Error("'c' should have exactly 1 neighbor")
	}
	if !isIn(dot, cNeighbors) {
		t.Errorf("'c' should have '.' as c neighbor")
	}

	dotNeighbors := g.Neighbors(dot)
	if len(dotNeighbors) != 3 {
		t.Error("'.' should have exactly 3 neighbors")
	}
	if !isIn(a, dotNeighbors) || !isIn(b, dotNeighbors) || !isIn(c, dotNeighbors) {
		t.Error("'.' should have 'a', 'b', and 'c' as neighbors")
	}
}

func getNodeAtCoord(seq map[coord]graph.Node, x, y int) (graph.Node, error) {
	for c, node := range seq {
		if c.x == x && c.y == y {
			return node, nil
		}
	}
	return graph.Node{}, errors.New("node does not exist at the specified coordinates")
}

// createPath takes a sequence of nodes separated by
// edges with weight = 1 in a single-line string
func createPath(s string) graph.Path {
	nodeCount := len(s)
	edgeCount := nodeCount - 1

	line := "#" + s + "#"
	border := strings.Repeat("#", utf8.RuneCountInString(line))
	maze := border + "\n" + line + "\n" + border
	runes := map2runes(maze)
	nodes := make(map[coord]graph.Node, 0)
	g := makeGraph(runes, &nodes)

	start, err := getNodeAtCoord(nodes, 1, 1)
	if err != nil {
		panic("cannot construct path. Bad test data.")
	}
	paths := g.DijkstraSearch(start)
	for _, path := range paths {
		// if all edges have weight = 1, this will be
		// the path from the first node to the last
		if path.Weight == edgeCount {
			return path
		}
	}

	return graph.Path{Weight: 0}
}

func TestIsPathToNonPeriodNode(t *testing.T) {
	// true cases
	s := "abcd"
	p := createPath(s)
	if !isPathToNonPeriodNode(p) {
		t.Errorf("'%s' ends in a non-period node", s)
	}

	s = "abcD"
	p = createPath(s)
	if !isPathToNonPeriodNode(p) {
		t.Errorf("'%s' ends in a non-period node", s)
	}

	s = "abc@"
	p = createPath(s)
	if !isPathToNonPeriodNode(p) {
		t.Errorf("'%s' ends in a non-period node", s)
	}

	s = "a...b"
	p = createPath(s)
	if !isPathToNonPeriodNode(p) {
		t.Errorf("'%s' ends in a non-period node", s)
	}

	// false cases
	s = "a"
	p = createPath(s)
	if isPathToNonPeriodNode(p) {
		t.Errorf("'%s' ends in a period", s)
	}

	s = "abc."
	p = createPath(s)
	if isPathToNonPeriodNode(p) {
		t.Errorf("'%s' ends in a period", s)
	}

	s = "a..."
	p = createPath(s)
	if isPathToNonPeriodNode(p) {
		t.Errorf("'%s' ends in a period", s)
	}
}

func TestIsPathToKey(t *testing.T) {
	// true cases
	s := "abcd"
	p := createPath(s)
	if !isPathToKey(p) {
		t.Errorf("'%s' ends in a key", s)
	}

	s = "a...f"
	p = createPath(s)
	if !isPathToKey(p) {
		t.Errorf("'%s' ends in a key", s)
	}

	// false cases
	s = "."
	p = createPath(s)
	if isPathToKey(p) {
		t.Errorf("'%s' does not end in a key", s)
	}

	s = "a...."
	p = createPath(s)
	if isPathToKey(p) {
		t.Errorf("'%s' does not end in a key", s)
	}

	s = "a...A"
	p = createPath(s)
	if isPathToKey(p) {
		t.Errorf("'%s' does not end in a key", s)
	}
}

func TestPathContainsOnlyPeriods(t *testing.T) {
	// true cases
	s := "a.d"
	p := createPath(s)
	if !pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' only contains periods between the start and end nodes", s)
	}

	s = "p.........q"
	p = createPath(s)
	if !pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' only contains periods between the start and end nodes", s)
	}

	s = "A.....M"
	p = createPath(s)
	if !pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' only contains periods between the start and end nodes", s)
	}

	s = "f..J"
	p = createPath(s)
	if !pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' only contains periods between the start and end nodes", s)
	}

	s = "ab"
	p = createPath(s)
	if !pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' only contains periods between the start and end nodes", s)
	}

	// false cases
	s = "a.b...c"
	p = createPath(s)
	if pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' contains a non-period node between the start and end nodes", s)
	}

	s = "a.Y...c"
	p = createPath(s)
	if pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' contains a non-period node between the start and end nodes", s)
	}

	s = "a@c"
	p = createPath(s)
	if pathContainsOnlyPeriods(p) {
		t.Errorf("'%s' contains a non-period node between the start and end nodes", s)
	}
}

func TestPruneGraph(t *testing.T) {
	maze := `#######
##...c#
##.####
#a...b#
#######
`

	runes := map2runes(maze)
	nodes := make(map[coord]graph.Node, 0)
	g := makeGraph(runes, &nodes)
	aNode, err := getNodeAtCoord(nodes, 1, 3)
	if err != nil {
		panic("Failed to find node 'a'")
	}

	neighbors := g.Neighbors(aNode)
	if len(neighbors) != 1 {
		t.Error("'a' should have only 1 neighbor.")
	}

	prunedG, newNodes := pruneGraph(g, nodes)
	aNode, err = getNodeAtCoord(newNodes, 1, 3)
	paths := prunedG.DijkstraSearch(aNode)
	if len(paths) != 3 {
		t.Errorf("'a' should have 3 paths (to self, to b, to c), but has %d", len(paths))
	}
	foundB := false
	foundC := false
	for _, path := range paths {
		if path.Weight == 4 {
			foundB = true
		}
		if path.Weight == 6 {
			foundC = true
		}
	}
	if !(foundB && foundC) {
		t.Error("a' should have direct paths to 'b' and 'c'")
	}
}

func TestPruneGraphWithDoor(t *testing.T) {
	maze := `#######
##B..c#
##.####
#a...b#
#######
`

	runes := map2runes(maze)
	nodes := make(map[coord]graph.Node, 0)
	g := makeGraph(runes, &nodes)
	aNode, err := getNodeAtCoord(nodes, 1, 3)
	if err != nil {
		panic("Failed to find node 'a'")
	}

	neighbors := g.Neighbors(aNode)
	if len(neighbors) != 1 {
		t.Error("'a' should have only 1 neighbor.")
	}

	prunedG, newNodes := pruneGraph(g, nodes)
	aNode, err = getNodeAtCoord(newNodes, 1, 3)
	paths := prunedG.DijkstraSearch(aNode)
	if len(paths) != 4 {
		t.Errorf("'a' should have 4 paths (to self, to B, to b, to c), but has %d", len(paths))
	}
	foundBdoor := false
	foundB := false
	foundC := false
	for _, path := range paths {
		if path.Weight == 3 {
			foundBdoor = true
		}
		if path.Weight == 4 {
			foundB = true
		}
		if path.Weight == 6 {
			foundC = true
		}
	}
	if !(foundB && foundC && foundBdoor) {
		t.Error("a' should have direct paths to 'b' and 'c' and 'B")
	}
}

func integrationTest(maze string, expected int, t *testing.T) {
	result := part1(maze)
	if result != expected {
		t.Errorf("Expected %d steps, got %d", expected, result)
	}
}

func TestIntegration01(t *testing.T) {
	maze := `########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`

	integrationTest(maze, 81, t)
}

func TestIntegration02(t *testing.T) {
	maze := `#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`

	integrationTest(maze, 136, t)
}

func TestIntegration03(t *testing.T) {
	maze := `#########
#b.A.@.a#
#########`

	integrationTest(maze, 8, t)
}

func TestIntegration04(t *testing.T) {
	maze := `#######
#a.#Cd#
##@#@##
#######
##@#@##
#cB#Ab#
#######`

	integrationTest(maze, 8, t)
}

func TestIntegration05(t *testing.T) {
	maze := `###############
#d.ABC.#.....a#
######@#@######
###############
######@#@######
#b.....#.....c#
###############`

	integrationTest(maze, 24, t)
}

func TestIntegration06(t *testing.T) {
	maze := `#############
#DcBa.#.GhKl#
#.###@#@#I###
#e#d#####j#k#
###C#@#@###J#
#fEbA.#.FgHi#
#############`

	integrationTest(maze, 32, t)
}

func TestIntegration07(t *testing.T) {
	maze := `#############
#g#f.D#..h#l#
#F###e#E###.#
#dCba@#@BcIJ#
#############
#nK.L@#@G...#
#M###N#H###.#
#o#m..#i#jk.#
#############`

	integrationTest(maze, 72, t)
}
