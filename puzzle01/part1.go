/*
Question:

Fuel required to launch a given module is based on its mass. Specifically, to find the fuel required for a module, take its mass, divide by three, round down, and subtract 2.

For example:

    For a mass of 12, divide by 3 and round down to get 4, then subtract 2 to get 2.
    For a mass of 14, dividing by 3 and rounding down still yields 4, so the fuel required is also 2.
    For a mass of 1969, the fuel required is 654.
    For a mass of 100756, the fuel required is 33583.

The Fuel Counter-Upper needs to know the total fuel requirement. To find it, individually calculate the fuel needed for the mass of each module (your puzzle input), then add together all the fuel values.

What is the sum of the fuel requirements for all of the modules on your spacecraft?

To build and run:
  $ go build part1.go
  $ cat input.txt | ./part1
*/

package main

import "bufio"
import "fmt"
import "os"
import "strconv"

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func fuel_cost(w int) int {
	return max(0, (w/3) - 2)
}


func main() {
	var ttl_fuel_cost int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		mass, err := strconv.Atoi(scanner.Text())

		if err != nil {
			fmt.Println("ERROR: Cannot parse input: " + scanner.Text())
			return
		}

		ttl_fuel_cost += fuel_cost(mass)
	}

	fmt.Println(ttl_fuel_cost)
}

