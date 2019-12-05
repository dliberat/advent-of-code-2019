/* https://adventofcode.com/2019/day/4

You arrive at the Venus fuel depot only to discover it's protected by a password. The Elves had written the password on a sticky note, but someone threw it out.

However, they do remember a few key facts about the password:

    It is a six-digit number.
    The value is within the range given in your puzzle input.
    Two adjacent digits are the same (like 22 in 122345).
    Going from left to right, the digits never decrease; they only ever increase or stay the same (like 111123 or 135679).

Other than the range rule, the following are true:

    111111 meets these criteria (double 11, never decreases).
    223450 does not meet these criteria (decreasing pair of digits 50).
    123789 does not meet these criteria (no double).

How many different passwords within the range given in your puzzle input meet these criteria?

Your puzzle input is 264360-746325.

*/
package main

import (
	"fmt"
	"math"
)

func flat(code []int) int {
	ttl := 0
	for i := 0; i < len(code); i++ {
		pow := len(code) - 1 - i
		ttl += code[i] * int(math.Pow10(pow))
	}
	return ttl
}

func hasDouble(code []int) bool {
	for i := 1; i < len(code); i++ {
		if code[i] == code[i-1] {
			return true
		}
	}
	return false
}

func countPerms(code []int, pos int, minAtPos int, min int, max int) int {
	perms := 0

	// base case. Last digit
	if pos == len(code)-1 {
		for i := minAtPos; i <= 9; i++ {
			code[pos] = i
			flattened := flat(code)
			if flattened < min {
				continue
			}
			if flattened > max {
				return perms
			}
			if hasDouble(code) {
				fmt.Println("Valid combo found.", code)
				perms++
			}
		}
		return perms
	}

	// recursive case
	for i := minAtPos; i <= 9; i++ {
		tmp := make([]int, len(code))
		copy(tmp, code)
		tmp[pos] = i
		if flat(tmp) > max {
			break
		}
		perms += countPerms(tmp, pos+1, i, min, max)
	}
	return perms
}

func main() {
	min := 264360
	max := 746325
	code := []int{0, 0, 0, 0, 0, 0}

	perms := countPerms(code, 0, 2, min, max)
	fmt.Println(perms)
}
