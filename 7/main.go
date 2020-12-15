/*

--- Day 7: Handy Haversacks ---
You land at the regional airport in time for your next flight. In fact, it looks like you'll even have time to grab
some food: all flights are currently delayed due to issues in luggage processing.

Due to recent aviation regulations, many rules (your puzzle input) are being enforced about bags and their contents;
bags must be color-coded and must contain specific quantities of other color-coded bags. Apparently, nobody
responsible for these regulations considered how long they would take to enforce!

For example, consider the following rules:

light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.

These rules specify the required contents for 9 bag types. In this example, every faded blue bag is empty, every
vibrant plum bag contains 11 bags (5 faded blue and 6 dotted black), and so on.

You have a shiny gold bag. If you wanted to carry it in at least one other bag, how many different bag colors would
be valid for the outermost bag? (In other words: how many colors can, eventually, contain at least one shiny gold bag?)

In the above rules, the following options would be available to you:

A bright white bag, which can hold your shiny gold bag directly.
A muted yellow bag, which can hold your shiny gold bag directly, plus some other bags.
A dark orange bag, which can hold bright white and muted yellow bags, either of which could then hold your shiny gold bag.
A light red bag, which can hold bright white and muted yellow bags, either of which could then hold your shiny gold bag.

So, in this example, the number of bag colors that can eventually contain at least one shiny gold bag is 4.

How many bag colors can eventually contain at least one shiny gold bag? (The list of rules is quite long; make sure you
get all of it.)


--- Part Two ---
It's getting pretty expensive to fly these days - not because of ticket prices, but because of the ridiculous number
of bags you need to buy!

Consider again your shiny gold bag and the rules from the above example:

faded blue bags contain 0 other bags.
dotted black bags contain 0 other bags.
vibrant plum bags contain 11 other bags: 5 faded blue bags and 6 dotted black bags.
dark olive bags contain 7 other bags: 3 faded blue bags and 4 dotted black bags.

So, a single shiny gold bag must contain 1 dark olive bag (and the 7 bags within it) plus 2 vibrant plum bags
(and the 11 bags within each of those): 1 + 1*7 + 2 + 2*11 = 32 bags!

Of course, the actual rules have a small chance of going several levels deeper than this example; be sure to count all
of the bags, even if the nesting becomes topologically impractical!

Here's another example:

shiny gold bags contain 2 dark red bags.
dark red bags contain 2 dark orange bags.
dark orange bags contain 2 dark yellow bags.
dark yellow bags contain 2 dark green bags.
dark green bags contain 2 dark blue bags.
dark blue bags contain 2 dark violet bags.
dark violet bags contain no other bags.

In this example, a single shiny gold bag must contain 126 other bags.

How many individual bags are required inside your single shiny gold bag?


*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

type FitsInRule struct {
	bag    string
	amount int
}

type ContainsRule struct {
	bag    string
	amount int
}

type ParsedRule struct {
	bag       string
	container string
	amount    int
}

type BagResult struct {
	bag    string
	amount int
	total  int
}

func main() {

	input, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var rules sync.Map
	var containerRules sync.Map
	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	for scanner.Scan() {
		ruleText := scanner.Text()
		ParsedRules := parseRule(ruleText)

		for _, rule := range ParsedRules {
			// first the fits in map:
			val, existed := rules.LoadOrStore(rule.bag, []FitsInRule{{rule.container, rule.amount}})
			if existed {
				fir := val.([]FitsInRule)
				fir = append(fir, FitsInRule{rule.container, rule.amount})
				rules.Store(rule.bag, fir)
			}
			// next the "contains" map:
			val, existed = containerRules.LoadOrStore(rule.container, []ContainsRule{{rule.bag, rule.amount}})
			if existed {
				cr := val.([]ContainsRule)
				cr = append(cr, ContainsRule{rule.bag, rule.amount})
				containerRules.Store(rule.container, cr)
			}
		}
	}

	ContainingBags := count(rules, "shiny gold bag", 1)
	Contains := countContaining(containerRules, "shiny gold bag") - 1

	fmt.Print("BAGS\n")
	for Bag := range ContainingBags {
		fmt.Printf("%s\n", Bag)
	}

	fmt.Printf("Bag types containing it: %d\n", len(ContainingBags))
	fmt.Printf("Bags it contains: %d\n", Contains)
}

func DebugRule(bag string, rule []FitsInRule) {
	fmt.Printf("%s\n", bag)
	for _, r := range rule {
		fmt.Printf("\t%s", r.bag)
	}
	fmt.Print("\n")
}

func count(rules sync.Map, bag string, amount int) map[string]struct{} {

	// There's probably a better way to deal with typing in
	// sync maps (and I'm not using goroutines in this puzzle anyway)
	UniqueBags := make(map[string]struct{})
	IFitsIn, ok := rules.Load(bag)
	if !ok {
		return UniqueBags
	}
	fir := IFitsIn.([]FitsInRule)
	if !ok {
		panic(fmt.Errorf("bag not found: %s", bag))
	}
	for _, rule := range fir {
		bags := count(rules, rule.bag, rule.amount)
		UniqueBags[rule.bag] = struct{}{}
		for k := range bags {
			UniqueBags[k] = struct{}{}
		}
	}

	return UniqueBags

}

func countContaining(rules sync.Map, bag string) int {
	IContains, ok := rules.Load(bag)
	sum := 1
	if !ok {
		return sum
	}
	cr := IContains.([]ContainsRule)

	for _, rule := range cr {
		sum += rule.amount * countContaining(rules, rule.bag)
	}

	return sum
}

func parseRule(text string) []ParsedRule {
	ruleBody := strings.Split(text, " contain ")
	ContainerBag := removePlural(ruleBody[0])
	var rules []ParsedRule
	// trim out the '.' at the end and split on commas
	contains := strings.Split(strings.ReplaceAll(ruleBody[1], ".", ""), ", ")
	for _, ContainRule := range contains {
		if ContainRule == "no other bags" {
			continue
		}

		// pull the integer out
		ContainDetails := strings.SplitN(ContainRule, " ", 2)
		amount, err := strconv.Atoi(ContainDetails[0])
		if err != nil {
			panic(err)
		}

		bag := ContainDetails[1]

		if amount != 1 {
			bag = removePlural(bag)
		}

		rules = append(rules, ParsedRule{bag, ContainerBag, amount})

	}

	return rules
}

func removePlural(bagText string) string {
	return bagText[:len(bagText)-1]
}
