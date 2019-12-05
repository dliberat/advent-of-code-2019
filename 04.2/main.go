/*https://adventofcode.com/2019/day/4*/
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
	seqLength := 1

	for i := 1; i < len(code); i++ {
		if code[i] == code[i-1] {
			seqLength++
		} else {
			if seqLength == 2 {
				return true
			}
			seqLength = 1
		}
	}
	return seqLength == 2
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
