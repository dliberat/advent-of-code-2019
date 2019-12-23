package main

import (
	"fmt"
)

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

type node struct {
	val      int
	children map[int]*node
}

type trie struct {
	root node
	f    vectorMaker
}

type vectorMaker func(...body) []int

type body struct {
	x    int
	y    int
	z    int
	velX int
	velY int
	velZ int
}

func (b *body) applyGravity(other body) {
	if other.x > b.x {
		b.velX++
	} else if other.x < b.x {
		b.velX--
	}

	if other.y > b.y {
		b.velY++
	} else if other.y < b.y {
		b.velY--
	}

	if other.z > b.z {
		b.velZ++
	} else if other.z < b.z {
		b.velZ--
	}
}

func (b *body) update() {
	b.x += b.velX
	b.y += b.velY
	b.z += b.velZ
}

func (b *body) getPotentialEnergy() int {
	return abs(b.x) + abs(b.y) + abs(b.z)
}

func (b *body) getKineticEnergy() int {
	return abs(b.velX) + abs(b.velY) + abs(b.velZ)
}

func makeBody(x int, y int, z int) body {
	return body{x: x, y: y, z: z, velX: 0, velY: 0, velZ: 0}
}

func getInput() []body {
	/*
			Puzzle input
		   <x=14, y=4, z=5>
		   <x=12, y=10, z=8>
		   <x=1, y=7, z=-10>
		   <x=16, y=-5, z=3>
	*/

	system := make([]body, 0)
	system = append(system, makeBody(14, 4, 5))
	system = append(system, makeBody(12, 10, 8))
	system = append(system, makeBody(1, 7, -10))
	system = append(system, makeBody(16, -5, 3))

	/*
			Sample input 1
		   <x=-1, y=0, z=2>
		   <x=2, y=-10, z=-7>
		   <x=4, y=-8, z=8>
		   <x=3, y=5, z=-1>
	*/

	// system := make([]body, 0)
	// system = append(system, makeBody(-1, 0, 2))
	// system = append(system, makeBody(2, -10, -7))
	// system = append(system, makeBody(4, -8, 8))
	// system = append(system, makeBody(3, 5, -1))

	/*
		Sample input 2
		<x=-8, y=-10, z=0>
		<x=5, y=5, z=10>
		<x=2, y=-7, z=3>
		<x=9, y=-8, z=-3>
	*/

	// system := make([]body, 0)
	// system = append(system, makeBody(-8, -10, 0))
	// system = append(system, makeBody(5, 5, 10))
	// system = append(system, makeBody(2, -7, 3))
	// system = append(system, makeBody(9, -8, -3))

	return system
}

func part1() {
	system := getInput()

	// time steps
	for ts := 0; ts < 1000; ts++ {

		for i := range system {
			for j, o := range system {
				if i == j {
					continue
				}

				system[i].applyGravity(o)
			}
		}

		for i := range system {
			system[i].update()
		}

	}

	ttlEnergy := 0
	for _, b := range system {
		ttlEnergy += b.getPotentialEnergy() * b.getKineticEnergy()
	}

	fmt.Println("[Part1] Total Energy: ", ttlEnergy)
}

func makeNode(val int) *node {
	return &node{val: val, children: make(map[int]*node)}
}

/*vectorX creates a vector out of the x position
and x velocity components of the system state*/
func vectorX(bodies ...body) []int {
	v := make([]int, len(bodies)*2)
	for i := range bodies {
		v[i] = bodies[i].x
		v[i+1] = bodies[i].velX
	}
	return v
}

/*vectorY creates a vector out of the y position
and y velocity components of the system state*/
func vectorY(bodies ...body) []int {
	v := make([]int, len(bodies)*2)
	for i := range bodies {
		v[i] = bodies[i].y
		v[i+1] = bodies[i].velY
	}
	return v
}

/*vectorZ creates a vector out of the z position
and z velocity components of the system state*/
func vectorZ(bodies ...body) []int {
	v := make([]int, len(bodies)*2)
	for i := range bodies {
		v[i] = bodies[i].z
		v[i+1] = bodies[i].velZ
	}
	return v
}

func (t *trie) insert(system []body) bool {
	current := &(t.root)

	// everything but the last body because
	// on the last body we need to do some special checking
	for _, b := range system[:len(system)-1] {
		vector := t.f(b)

		for _, v := range vector {
			if current.children[v] == nil {
				current.children[v] = makeNode(v)
			}
			current = current.children[v]

		}
	}

	// final body. Check if the leaf node already exists
	b := system[len(system)-1]

	vector := t.f(b)

	for _, v := range vector[:len(vector)-1] {
		if current.children[v] == nil {
			current.children[v] = makeNode(v)
		}
		current = current.children[v]
	}

	if current.children[vector[len(vector)-1]] == nil {
		current.children[vector[len(vector)-1]] = makeNode(vector[len(vector)-1])
		return false
	}
	return true
}

/*Greatest common divisor*/
func gcd(x int, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

/*Least-common multiple*/
func lcm(nums ...int) int {

	if len(nums) > 2 {
		return lcm(nums[0], lcm(nums[1:]...))
	}

	return nums[0] * nums[1] / gcd(nums[0], nums[1])

}

func findPeriod(f vectorMaker) int {
	system := getInput()

	// The trie will keep track of every single state of the system
	// with respect to the specified axis (x, y, or z, as determined
	// by the vectorMaker function) as it evolves
	t := trie{root: node{val: 0, children: make(map[int]*node)}, f: f}

	t.insert(system)

	var timeStepAtRepeat int
	var repeatedState []int

	ts := 0
	for true {

		for i := range system {
			for j, o := range system {
				if i == j {
					continue
				}
				system[i].applyGravity(o)
			}
		}

		for i := range system {
			system[i].update()
		}

		if t.insert(system) {
			// Reverted back to a previously visited state.
			// From this point forward, the system's evolution
			// will repeat in a cycle. If the cycle includes the initial
			// state of the system, the period would be equal to ts
			timeStepAtRepeat = ts
			repeatedState = f(system...)
			break
		}

		ts++
	}

	// it's possible that the cycle doesn't start at the
	// initial condition, so we need to spin it up again and
	// make sure about the period length
	system = getInput()

	ts = 0
	for true {

		currentState := f(system...)
		for v := range currentState {
			if currentState[v] != repeatedState[v] {
				break
			}
			// found the initial occurence of the repeated state!
			return timeStepAtRepeat - ts + 1
		}

		for i := range system {
			for j, o := range system {
				if i == j {
					continue
				}
				system[i].applyGravity(o)
			}
		}
		for i := range system {
			system[i].update()
		}

		ts++
	}

	return -1 // this should never happen
}

func part2() {

	/* Because the x, y, and z components of state (positions and velocities)
	all evolve independently of each other, we can find the period of the
	functions on each axis without needing to consider the other two axes.*/

	xAxisPeriod := findPeriod(vectorX)
	fmt.Println("The x axis cycles every", xAxisPeriod, "time steps")

	yAxisPeriod := findPeriod(vectorY)
	fmt.Println("The y axis cycles every", yAxisPeriod, "time steps")

	zAxisPeriod := findPeriod(vectorZ)
	fmt.Println("The z axis cycles every", zAxisPeriod, "time steps")

	/*Then, assuming that the initial state of each axis is part of
	the stable periodic funcion, the lowest common multiple of the
	three periods will give us the first time index at which all
	components of state return to their initial conditions.
	(Note the assumption: If the planets take some amount of time
	to settle into a periodic pattern, the lcm method would not work)*/
	period := lcm(xAxisPeriod, yAxisPeriod, zAxisPeriod)

	fmt.Println("[Part2]: The period of the moons' orbits is", period)
}

func main() {
	part1()
	part2()
}
