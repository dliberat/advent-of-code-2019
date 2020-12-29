package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type reaction struct {
	outputName string
	outputAmt  int
	inputs     map[string]int
}

/* ========================= PRIORITY QUEUE ========================= */

type queueItem struct {
	r           reaction
	amtRequired int
}

type fifoQueue struct {
	items []queueItem
}

func (q *fifoQueue) push(r reaction, amt int) {
	item := queueItem{r: r, amtRequired: amt}
	q.items = append(q.items, item)
}

func (q *fifoQueue) pop() (reaction, int) {
	i := q.items[0]
	q.items = q.items[1:]
	return i.r, i.amtRequired
}

func (q *fifoQueue) size() int {
	return len(q.items)
}

/* =========================  ========================= */

func scalarMultiplyMap(m map[string]int, n int) map[string]int {
	output := make(map[string]int)
	for k, v := range m {
		output[k] = v * n
	}
	return output
}

func findReaction(name string, reactions []reaction) reaction {
	for _, r := range reactions {
		if r.outputName == name {
			return r
		}
	}
	panic("Reaction " + name + " does not exist.")
}

func part1(reactions []reaction, targetAmt int) int {

	totalReqdOreAmt := 0
	onHand := make(map[string]int)
	fuel := findReaction("FUEL", reactions)
	queue := fifoQueue{}
	queue.push(fuel, targetAmt)

	for queue.size() > 0 {
		target, amt := queue.pop()

		// in order to make the target, we need this many of each input element
		required := scalarMultiplyMap(target.inputs, amt)

		// before trying to make more, use up anything we have on hand
		for requirementName, requiredAmt := range required {
			if onHand[requirementName] > 0 {
				diff := requiredAmt - onHand[requirementName]
				if diff > 0 {
					// we can cover part of the requirement using on-hand materials.
					// Therefore, reduce the amount required by however many we have on hand
					required[requirementName] = requiredAmt - onHand[requirementName]
					// and use up anything we have on hand
					onHand[requirementName] = 0
				} else {
					// we can cover the entire requirement using on-hand materials
					onHand[requirementName] = onHand[requirementName] - requiredAmt
					required[requirementName] = 0
				}
			}
		}

		// Whatever is left in 'required' needs to be created
		for requirementName, requiredAmt := range required {
			if requirementName == "ORE" {
				// ore has no constituent elements (it is the raw material)
				totalReqdOreAmt += requiredAmt
				continue
			} else if requiredAmt == 0 {
				continue
			}

			r := findReaction(requirementName, reactions)
			needToMake := 1
			for r.outputAmt*needToMake < requiredAmt {
				needToMake++
			}

			queue.push(r, needToMake)

			// depending on the nature of the reaction, we may need to make
			// more than we really need (e.g. if the reaction makes element
			// X in increments of 5, but we only need 3, we will have 2 left
			// over). Keep this amount on hand since we might use those
			// leftovers for other reactions.
			leftOver := r.outputAmt*needToMake - requiredAmt
			onHand[requirementName] = onHand[requirementName] + leftOver
		}

		// repeat until the queue is empty
	}

	return totalReqdOreAmt
}

func part2(reactions []reaction, totalOre int) int {
	targetFuel := 2_000_000
	increment := 1_000_000

	for increment >= 1 {
		reqOre := 0

		for reqOre < totalOre {
			targetFuel += increment
			reqOre = part1(reactions, targetFuel)
			fmt.Println(reqOre, "ORE =>", targetFuel, "FUEL")
		}

		targetFuel -= increment

		if increment == 1 {
			break
		}
		increment /= 10
	}

	return targetFuel
}

func main() {
	input := getInput("input.txt")
	reactions := parseReactionList(input)

	p1 := part1(reactions, 1)
	fmt.Println("[PART 1]", p1, "ORE are required in order to make 1 FUEL")

	p2 := part2(reactions, 1_000_000_000_000)
	fmt.Println("[PART 2] Up to", p2, "FUEL can be generated with 1 trillion ORE")
}

/* ========================= INPUT PROCESSING ========================= */

func getInput(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Cannot read input program.")
	}
	return string(data)
}

/**parseReactionLeftHandSide reads the left side of a reaction
string and returns a map whose keys are the elements and values are the
amounts associated with each element.
Ex:
  parseReactionLeftHandSide("3 A, 4 B") == {"A": 3, "B": 4}
*/
func parseReactionLeftHandSide(txt string) map[string]int {
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

/**parseReactionRightHandSide reads the right side of a reaction string and
returns a reaction struct with its outputName and outputAmt set. The
inputs to the reaction must be processed separately using the left
hand side of the reaction string.*/
func parseReactionRightHandSide(txt string) reaction {
	s := strings.Split(txt, " ")
	amt, err := strconv.Atoi(s[0])
	if err != nil {
		panic("Invalid right-hand side for the reaction.")
	}
	r := reaction{outputAmt: amt, outputName: s[1]}
	return r
}

/**parseReactionList reads the entire puzzle input and returns a slice
of reaction structs representing the same data.*/
func parseReactionList(txt string) []reaction {
	reactions := make([]reaction, 0)

	for _, line := range strings.Split(txt, "\n") {

		arrowIndex := strings.Index(line, "=>")
		if arrowIndex < 0 {
			// not a reaction
			continue
		}

		lhs := strings.Trim(line[0:arrowIndex], " ")
		lhMap := parseReactionLeftHandSide(lhs)

		rhs := strings.Trim(line[arrowIndex+2:], " ")
		r := parseReactionRightHandSide(rhs)
		r.inputs = lhMap

		reactions = append(reactions, r)
	}

	return reactions
}
