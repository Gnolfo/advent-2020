/*
--- Day 2: Password Philosophy ---
Your flight departs in a few days from the coastal airport; the easiest way down to the coast from here is via
toboggan.

The shopkeeper at the North Pole Toboggan Rental Shop is having a bad day. "Something's wrong with our computers;
we can't log in!" You ask if you can take a look.

Their password database seems to be a little corrupted: some of the passwords wouldn't have been allowed by the
Official Toboggan Corporate Policy that was in effect when they were chosen.

To try to debug the problem, they have created a list (your puzzle input) of passwords (according to the corrupted
database) and the corporate policy when that password was set.

For example, suppose you have the following list:

1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc
Each line gives the password policy and then the password. The password policy indicates the lowest and highest
number of times a given letter must appear for the password to be valid. For example, 1-3 a means that the password
must contain a at least 1 time and at most 3 times.

In the above example, 2 passwords are valid. The middle password, cdefg, is not; it contains no instances of b, but
needs at least 1. The first and third passwords are valid: they contain one a or nine c, both within the limits of
their respective policies.

How many passwords are valid according to their policies?

--- Part Two ---
While it appears you validated the passwords correctly, they don't seem to be what the Official Toboggan Corporate
Authentication System is expecting.

The shopkeeper suddenly realizes that he just accidentally explained the password policy rules from his old job at
the sled rental place down the street! The Official Toboggan Corporate Policy actually works a little differently.

Each policy actually describes two positions in the password, where 1 means the first character, 2 means the second
character, and so on. (Be careful; Toboggan Corporate Policies have no concept of "index zero"!) Exactly one of these
positions must contain the given letter. Other occurrences of the letter are irrelevant for the purposes of policy
enforcement.

Given the same example list from above:

1-3 a: abcde is valid: position 1 contains a and position 3 does not.
1-3 b: cdefg is invalid: neither position 1 nor position 3 contains b.
2-9 c: ccccccccc is invalid: both position 2 and position 9 contain c.
How many passwords are valid according to the new interpretation of the policies?
*/

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

type PasswordChecker struct {
	password string
	char     rune
	min, max int
}

type NewPasswordChecker struct {
	password string
	char     rune
	pos      [2]int
}

func main() {

	input, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	rows := csv.NewReader(strings.NewReader(string(input)))
	rows.Comma = ':'
	rows.TrimLeadingSpace = true

	var arr []PasswordChecker
	var arr2 []NewPasswordChecker

	for {
		row, err := rows.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		char, min, max := parseRule(row[0])
		password := row[1]

		arr = append(arr, PasswordChecker{password, char, min, max})
		arr2 = append(arr2, NewPasswordChecker{password, char, [2]int{min - 1, max - 1}})
	}

	countChan := make(chan bool, len(arr))
	newCountChan := make(chan bool, len(arr))

	validCount := 0
	newValidCount := 0
	var wg1, wg2 sync.WaitGroup
	wg1.Add(len(arr))
	wg2.Add(len(arr))
	for i := 0; i < len(arr); i++ {
		go arr[i].check(countChan, &wg1)
		go arr2[i].check(newCountChan, &wg2)
	}

	go func() {
		wg1.Wait()
		close(countChan)
	}()

	go func() {
		wg2.Wait()
		close(newCountChan)
	}()

	var wgTop sync.WaitGroup
	wgTop.Add(2)
	go func() {
		defer wgTop.Done()
		for result := range countChan {
			if result {
				validCount++
			}
		}
		fmt.Printf("part 1: %d\n", validCount)
	}()

	go func() {
		defer wgTop.Done()
		for result := range newCountChan {
			if result {
				newValidCount++
			}
		}
		fmt.Printf("part 2: %d\n", newValidCount)
	}()

	wgTop.Wait()
}

func parseRule(rule string) (rune, int, int) {
	var min, max int
	var err error
	rules := strings.Split(rule, " ")
	minmax := strings.Split(rules[0], "-")
	min, err = strconv.Atoi(minmax[0])
	if err != nil {
		panic(err)
	}
	max, err = strconv.Atoi(minmax[1])
	if err != nil {
		panic(err)
	}
	return rune(rules[1][0]), min, max
}

func (pc PasswordChecker) check(c chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	count := strings.Count(pc.password, string(pc.char))
	if count >= pc.min && count <= pc.max {
		c <- true
	} else {
		c <- false
	}

}

func (pc NewPasswordChecker) check(c chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	// password is too short
	if len(pc.password) < pc.pos[0] {
		c <- false
		return
	}
	// short enough for only first posiion
	if len(pc.password) < pc.pos[1] {
		c <- rune(pc.password[pc.pos[0]]) == pc.char
		return
	}

	if (rune(pc.password[pc.pos[0]]) == pc.char && rune(pc.password[pc.pos[1]]) != pc.char) ||
		(rune(pc.password[pc.pos[0]]) != pc.char && rune(pc.password[pc.pos[1]]) == pc.char) {
		c <- true
	} else {
		c <- false
	}
}
