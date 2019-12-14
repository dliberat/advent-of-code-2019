package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type node struct {
	name     string
	parent   *node
	children map[string]*node
	d        int
	visited  bool
}

func (b *node) depth() int {
	// memoized value
	if b.d > -1 {
		return b.d
	}

	if b.parent != nil {
		b.d = b.parent.depth() + 1
	} else {
		b.d = 0
	}
	return b.d
}

func (b *node) addParent(parent *node) {
	b.parent = parent
	parent.children[b.name] = b
}

func (b *node) neighbors() []*node {
	count := len(b.children)
	if b.parent != nil {
		count++
	}

	neighs := make([]*node, count)
	if b.parent != nil {
		neighs = append(neighs, b.parent)
	}

	for _, v := range b.children {
		neighs = append(neighs, v)
	}

	return neighs
}

func (b *node) visit() {
	b.visited = true
}

func bfs(queue []*node, depth int) int {
	var nextQueue []*node

	for _, candidate := range queue {

		neighs := candidate.neighbors()
		for _, n := range neighs {
			if n == nil || n.visited {
				continue
			}
			if n.name == "SAN" {
				return depth
			}
			nextQueue = append(nextQueue, n)
		}

		candidate.visit()

	}

	return bfs(nextQueue, depth+1)
}

func partOne(orbits map[string]*node) {
	// sum over all node depths
	ttl := 0
	for _, v := range orbits {
		d := v.depth()
		ttl += d
	}
	fmt.Println("[Part 1] Total depths:", ttl)
}

func partTwo(orbits map[string]*node) {
	// min distance between parents of node YOU
	// and node SAN. Breadth-first-search

	var queue []*node

	queue = append(queue, orbits["YOU"].parent)

	distance := bfs(queue, 0)
	fmt.Println("[Part 2] Single-source shortest path: ", distance)

}

func main() {

	var orbits = make(map[string]*node)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		// builds a tree with node "COM" at the root.
		// Every node has 1 in-edge, and an arbitrary number of out-edges

		relation := strings.Split(scanner.Text(), ")")
		eeName := relation[0]
		erName := relation[1]

		_, ok := orbits[eeName]
		if !ok {
			orbits[eeName] = &node{name: eeName, children: make(map[string]*node), d: -1}
		}
		_, ok = orbits[erName]
		if !ok {
			orbits[erName] = &node{name: erName, children: make(map[string]*node), d: -1}
		}

		orbits[erName].addParent(orbits[eeName])
	}

	partOne(orbits)
	partTwo(orbits)

}
