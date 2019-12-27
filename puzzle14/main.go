package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

type reaction struct {
	reactants map[string]int
	outputs   map[string]int
}

func min(nums []int) int {
	least := math.Inf(1)
	for _, n := range nums {
		if float64(n) < least {
			least = float64(n)
		}
	}
	return int(least)
}

func getInput() string {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read input program.")
	}
	return string(data)
}

func parseReactionSide(txt string) map[string]int {
	m := make(map[string]int)

	elements := strings.Split(txt, ",")
	for _, e := range elements {
		element := strings.Trim(e, " ")
		elementArr := strings.Split(element, " ")
		count, err := strconv.Atoi(strings.Trim(elementArr[0], " "))
		if err != nil {
			panic("Unable to parse number.")
		}
		name := strings.Trim(elementArr[1], " ")
		m[name] = count
	}
	return m
}

func parseReactionList(txt string) []reaction {
	reactions := make([]reaction, 0)

	lst := strings.Split(txt, "\n")
	for _, line := range lst {

		arrowIndex := strings.Index(line, "=>")
		if arrowIndex < 0 {
			fmt.Println("Error. Cannot find arrow in reaction:", line)
			continue
		}

		lhs := strings.Trim(line[0:arrowIndex], " ")
		lhMap := parseReactionSide(lhs)

		rhs := strings.Trim(line[arrowIndex+2:], " ")
		rhMap := parseReactionSide(rhs)

		react := reaction{reactants: lhMap, outputs: rhMap}

		reactions = append(reactions, react)
	}

	return reactions
}

func mapHash(a *map[string]int) string {
	arr := make([]string, len(*a))
	i := 0
	for k, v := range *a {
		arr[i] = fmt.Sprintf("%s%d", k, v)
		i++
	}
	sort.Strings(arr)
	return strings.Join(arr, "")
}

func cloneMap(a *map[string]int) map[string]int {
	b := make(map[string]int)
	for k, v := range *a {
		b[k] = v
	}
	return b
}

func findReactionWithOutput(k string, reactions *[]reaction) reaction {
	for i := range *reactions {
		for key := range (*reactions)[i].outputs {
			if key == k {
				return (*reactions)[i]
			}
		}
	}

	panic("No reaction with the requested output")
}

func checkCompletion(a *map[string]int) bool {
	for k, v := range *a {
		if k == "ORE" || v <= 0 {
			continue
		}
		return false
	}
	return true
}

func getCandidateReactions(a *map[string]int, reactions *[]reaction) []reaction {
	candidates := make([]reaction, 0)
	for k, v := range *a {
		if k != "ORE" && v > 0 {
			r := findReactionWithOutput(k, reactions)
			candidates = append(candidates, r)
		}
	}
	return candidates
}

func applyOneReduction(a *map[string]int, r reaction) {

	for reactant, amt := range r.reactants {
		(*a)[reactant] = (*a)[reactant] + amt
	}
	for output, amt := range r.outputs {
		(*a)[output] = (*a)[output] - amt
	}

}

func nonLossyReductions(a *map[string]int, reactions *[]reaction) {

	didReactionsFlag := true // exit condition

	for didReactionsFlag {

		b := make(map[string]int)

		didReactionsFlag = false

		// k = element
		// v = amount needed
		for k, v := range *a {

			if v == 0 {
				continue
			}

			if b[k] > 0 {
				v += b[k]
				b[k] = 0
			}

			if k == "ORE" { // ORE is not produced by any reaction
				b[k] = v
				continue
			}

			r := findReactionWithOutput(k, reactions)

			// while the output is still less than or equal to the entire needed amount,
			// do as many reactions as we can without going over.
			// this can be optimized later
			for r.outputs[k] <= v {
				didReactionsFlag = true

				b[k] = v

				for reactant, amt := range r.reactants {
					b[reactant] = b[reactant] + amt
				}
				for output, amt := range r.outputs {
					// it's possible to have more on hand than we need. Therefore, negative numbers
					b[output] = b[output] - amt
				}

				v = b[k]
				b[k] = 0
			}

			b[k] = b[k] + v
		}

		*a = b
	}

}

var visited map[string]int = make(map[string]int)
var best int = math.MaxInt64
var limit int = 500

func dfs(a map[string]int, reactions *[]reaction, depth int) int {

	if a["ORE"] >= best {
		return best
	}

	if checkCompletion(&a) {
		fmt.Println("hit completion at 209 with depth", depth, "and ore:", a["ORE"])
		if a["ORE"] < best {
			best = a["ORE"]
		}
		return best
	}

	limit--
	if limit <= 0 {
		return best
	}

	nonLossyReductions(&a, reactions)

	hash := mapHash(&a)
	if visited[hash] > 0 {
		return visited[hash]
	}

	if a["ORE"] >= best {
		visited[hash] = a["ORE"]
		return best
	}

	if checkCompletion(&a) {
		if a["ORE"] < best {
			best = a["ORE"]
		}
		return best
	}

	candidates := getCandidateReactions(&a, reactions)
	scores := make([]int, len(candidates))

	for i := range candidates {
		clone := cloneMap(&a)
		applyOneReduction(&clone, candidates[i])

		scores[i] = dfs(clone, reactions, depth+1)

	}

	score := min(scores)
	if score < best {
		best = score
	}

	visited[hash] = score
	return score
}

func part1() {
	input := getInput()
	reactions := parseReactionList(input)

	a := make(map[string]int)
	a["FUEL"] = 1

	// stabilizes to the correct answer, but takes too long to traverse
	// the entire graph of possible combinations. Fortunately, we can set
	// a recursion limit and be fairly confident that it will
	// return the correct result
	dfs(a, &reactions, 0)

	fmt.Println("[Part 1] At least", best, "ore is needed to generate 1 fuel.")
}

func main() {
	part1()
}
