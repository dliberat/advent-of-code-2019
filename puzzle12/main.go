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

func (b *body) print() string {
	return fmt.Sprintf("<body x: %d, y: %d, z: %d; vX: %d, vY: %d, vZ: %d>", b.x, b.y, b.z, b.velX, b.velY, b.velZ)
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
		//fmt.Println("Time step: ", ts)

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
			// fmt.Println(system[i].print())
		}
		//fmt.Println("")

	}

	ttlEnergy := 0
	for _, b := range system {
		ttlEnergy += b.getPotentialEnergy() * b.getKineticEnergy()
	}

	fmt.Println("[Part1] Total Energy: ", ttlEnergy)
}

func getPositionHash(a body, b body, c body, d body) string {
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d%d",
		a.x, a.y, a.z,
		b.x, b.y, b.z,
		c.x, c.y, c.z,
		d.x, d.y, d.z)
}
func getVelocityHash(a body, b body, c body, d body) string {
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d%d",
		a.velX, a.velY, a.velZ,
		b.velX, b.velY, b.velZ,
		c.velX, c.velY, c.velZ,
		d.velX, d.velY, d.velZ)
}

func part2() {
	/* Brute force algorithm. Won't work. */

	system := getInput()

	var ts int64 = 1
	var phash string
	var vhash string
	var posMap map[string]map[string]int64
	posMap = make(map[string]map[string]int64)

	for true {

		phash = getPositionHash(system[0], system[1], system[2], system[3])
		vhash = getVelocityHash(system[0], system[1], system[2], system[3])

		if posMap[phash] == nil {
			posMap[phash] = make(map[string]int64)
			posMap[phash][vhash] = ts
		} else {
			if posMap[phash][vhash] == 0 {
				posMap[phash][vhash] = ts
			} else {
				fmt.Println("[Part 2] Found repeated state at time step", ts-1)
				return
			}
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
}

func main() {
	part1()
	part2()
}
