package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	TOTAL_ROUND int = 10000
)

const (
	MULTI Operator = iota
	ADD
	SUB
	UNKNOWN
)

type Operator int

type Outcome struct {
	Success int
	Failure int
}

type Operation struct {
	Action Operator
	Value  string
}

func (o Operator) string() string {
	switch o {
	case MULTI:
		return "multiplied by"
	case ADD:
		return "increased by"
	case SUB:
		return "substracted by"
	default:
		return "UNKNOWN"
	}
}

func GetOperator(op string) Operator {
	switch op {
	case "*":
		return MULTI
	case "+":
		return ADD
	case "-":
		return SUB
	default:
		return UNKNOWN
	}
}

type MonkeyTeam struct {
	Monkey         []Monkey
	CompletedRound int
}

// Return top X monkey inspected times
func (mt *MonkeyTeam) getTopMonkey(topX int) []int {
	inspectedTimes := []int{}
	for _, mk := range mt.Monkey {
		inspectedTimes = append(inspectedTimes, mk.itemInspected)
	}

	sort.Ints(inspectedTimes)

	return inspectedTimes[len(inspectedTimes)-topX:]
}

func (mt *MonkeyTeam) print(round int) {
	fmt.Printf("== After round %v ==\n", round)
	for i, mk := range mt.Monkey {
		fmt.Printf("Monkey %v inspected items %v times\n", i, mk.itemInspected)
	}
}

type Monkey struct {
	// lists your worry level for each item the monkey is currently holding in the order they will be inspected.
	Items []*big.Int

	// Operation shows how your worry level changes as that monkey inspects an item.
	// (An operation like new = old * 5 means that your worry level after the monkey inspected the item is five times whatever your worry level was before inspection.)
	Op Operation

	// Test shows how the monkey uses your worry level to decide where to throw an item next.
	Test *big.Int

	// shows what happens with an item if the Test was true.
	// shows what happens with an item if the Test was false.
	Outcome Outcome

	// Counter Item inspected
	itemInspected int
}

func NewMonkey(input []string) Monkey {
	// Expecting a valid raw input like this:
	//
	// Monkey 0:
	// 	Starting items: 79, 98
	// 	Operation: new = old * 19
	// 	Test: divisible by 23
	// 		If true: throw to monkey 2
	// 		If false: throw to monkey 3

	rawItems := strings.Split(input[1], ":")[1]
	rawOperations := strings.Split(input[2], ":")[1]
	rawTest := strings.Split(input[3], ":")[1]
	rawOutcomeTrue := strings.Split(input[4], ":")[1]
	rawOutcomeFalse := strings.Split(input[5], ":")[1]

	// Parse items as int arrays
	items := []*big.Int{}
	for _, item := range strings.Split(rawItems, ",") {
		it, _ := new(big.Int).SetString(strings.TrimSpace(item), 0) // strconv.Atoi(strings.TrimSpace(item))
		items = append(items, it)
	}

	// Parse operations
	re := regexp.MustCompile(`new = old ([*+]) (\w+)`)
	opMatch := re.FindStringSubmatch(rawOperations)
	// opValue, _ := strconv.Atoi(opMatch[2])

	// Parse test
	re = regexp.MustCompile(`divisible by (\d+)`)
	testMatch := re.FindStringSubmatch(rawTest)
	testValue, _ := new(big.Int).SetString(testMatch[1], 0)

	// Parse Outcomes
	re = regexp.MustCompile(`throw to monkey (\d+)`)
	outcomeTrueMatch := re.FindStringSubmatch(rawOutcomeTrue)
	outcomeTrueMonkeyValue, _ := strconv.Atoi(outcomeTrueMatch[1])

	outcomeFalseMatch := re.FindStringSubmatch(rawOutcomeFalse)
	outcomeFalseMonkeyValue, _ := strconv.Atoi(outcomeFalseMatch[1])

	monkey := Monkey{
		Items: items,
		Op: Operation{
			Action: GetOperator(opMatch[1]),
			Value:  opMatch[2],
		},
		Test: testValue,
		Outcome: Outcome{
			Success: outcomeTrueMonkeyValue,
			Failure: outcomeFalseMonkeyValue,
		},
		itemInspected: 0,
	}

	return monkey
}

func (op *Operation) apply(old *big.Int) *big.Int {
	var opValue *big.Int
	if op.Value == "old" {
		opValue = old
	} else {
		opValue, _ = new(big.Int).SetString(op.Value, 0)
	}

	new := big.NewInt(0)
	switch op.Action {
	case MULTI:
		new.Mul(old, opValue)
	case ADD:
		new.Add(old, opValue)
	case SUB:
		new.Sub(old, opValue)
	}

	// fmt.Printf("\t\tWorry level is %s %v to %v.\n", op.Action.string(), opValue, new)
	return new
}

func (mt *MonkeyTeam) add(m Monkey) {
	mt.Monkey = append(mt.Monkey, m)
}

func (mt *MonkeyTeam) playRound() {
	// Fully Process each Monkey per Round
	for i, mk := range mt.Monkey {
		// fmt.Printf("Monkey %v\n", i)
		for _, item := range mk.Items {
			// Monkey Inspect by applying its Operation
			// fmt.Printf("\tMonkey inspects an item with a worry level of %v.\n", item)
			mt.Monkey[i].itemInspected++
			item = mk.Op.apply(item)

			// Worry level divided by 3 and round down to nearest integer
			// item = int(math.Trunc(float64(item) / float64(3)))
			// fmt.Printf("\t\tMonkey gets bored with item. Worry level is divided by 3 to %v.\n", item)

			// Test Outcome
			// _, m1 := new(big.Int).DivMod(item, big.NewInt(1), new(big.Int))
			// if len(m1.Bits()) == 0 {
			// 	fmt.Printf("it works 1")
			// }

			_, m := new(big.Int).DivMod(item, mk.Test, new(big.Int))
			if len(m.Bits()) == 0 {
				// fmt.Printf("\t\tCurrent worry level is divisible by %v.\n", mk.Test)
				mt.Monkey[mk.Outcome.Success].Items = append(mt.Monkey[mk.Outcome.Success].Items, item)
				// fmt.Printf("\t\tItem with worry level %v is thrown to monkey %v.\n", item, mk.Outcome.Success)
			} else {
				// fmt.Printf("\t\tCurrent worry level is not divisible by %v.\n", mk.Test)
				mt.Monkey[mk.Outcome.Failure].Items = append(mt.Monkey[mk.Outcome.Failure].Items, item)
				// fmt.Printf("\t\tItem with worry level %v is thrown to monkey %v.\n", item, mk.Outcome.Failure)
			}
		}
		// Empty all items from the monkey
		mt.Monkey[i].Items = []*big.Int{}
	}
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	monkeyRaw := []string{}
	mt := MonkeyTeam{
		Monkey:         []Monkey{},
		CompletedRound: 0,
	}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			monkey := NewMonkey(monkeyRaw)
			mt.add(monkey)
			monkeyRaw = []string{}
		} else {
			monkeyRaw = append(monkeyRaw, line)
		}
	}
	monkey := NewMonkey(monkeyRaw)
	mt.add(monkey)

	roundCompleted := 0
	fullRunStartTime := time.Now()
	for roundCompleted < TOTAL_ROUND {
		// fmt.Printf("monky 0 , item 0 , base10: %v", mt.Monkey[0].Items[0].Text(10))
		startRound := time.Now()
		mt.playRound()
		endRound := time.Now()
		roundDuration := endRound.Sub(startRound)
		totalDuration := endRound.Sub(fullRunStartTime)
		roundOutDuration := time.Time{}.Add(roundDuration)
		totalOutDuration := time.Time{}.Add(totalDuration)
		roundCompleted++
		fmt.Printf("Round %v completed in: %s ; total time at: %s\n", roundCompleted, roundOutDuration.Format("15:04:05"), totalOutDuration.Format("15:04:05"))
		if roundCompleted == 1 || roundCompleted == 20 || roundCompleted%1000 == 0 {
			mt.print(roundCompleted)
		}

	}

	// Part1
	topInspected := mt.getTopMonkey(2)
	fmt.Println("TopInspected = ", topInspected)
	fmt.Printf("What is the level of monkey business after 20 rounds of stuff-slinging simian shenanigans?: %v", topInspected[0]*topInspected[1])
}
