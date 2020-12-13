/*
--- Day 4: Passport Processing ---
You arrive at the airport only to realize that you grabbed your North Pole Credentials instead of your passport. While
these documents are extremely similar, North Pole Credentials aren't issued by a country and therefore aren't actually
valid documentation for travel in most of the world.

It seems like you're not the only one having problems, though; a very long line has formed for the automatic passport
scanners, and the delay could upset your travel itinerary.

Due to some questionable network security, you realize you might be able to solve both of these problems at the same time.

The automatic passport scanners are slow because they're having trouble detecting which passports have all required
fields. The expected fields are as follows:

byr (Birth Year)
iyr (Issue Year)
eyr (Expiration Year)
hgt (Height)
hcl (Hair Color)
ecl (Eye Color)
pid (Passport ID)
cid (Country ID)
Passport data is validated in batch files (your puzzle input). Each passport is represented as a sequence of key:value
pairs separated by spaces or newlines. Passports are separated by blank lines.

Here is an example batch file containing four passports:

ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in
The first passport is valid - all eight fields are present. The second passport is invalid - it is missing hgt
(the Height field).

The third passport is interesting; the only missing field is cid, so it looks like data from North Pole Credentials,
not a passport at all! Surely, nobody would mind if you made the system temporarily ignore missing cid fields. Treat
this "passport" as valid.

The fourth passport is missing two fields, cid and byr. Missing cid is fine, but missing any other field is not, so
this passport is invalid.

According to the above rules, your improved system would report 2 valid passports.

Count the number of valid passports - those that have all required fields. Treat cid as optional. In your batch file,
how many passports are valid?

--- Part Two ---
The line is moving more quickly now, but you overhear airport security talking about how passports with invalid data
are getting through. Better add some data validation, quick!

You can continue to ignore the cid field, but each other field has strict rules about what values are valid for
automatic validation:

byr (Birth Year) - four digits; at least 1920 and at most 2002.
iyr (Issue Year) - four digits; at least 2010 and at most 2020.
eyr (Expiration Year) - four digits; at least 2020 and at most 2030.
hgt (Height) - a number followed by either cm or in:
If cm, the number must be at least 150 and at most 193.
If in, the number must be at least 59 and at most 76.
hcl (Hair Color) - a # followed by exactly six characters 0-9 or a-f.
ecl (Eye Color) - exactly one of: amb blu brn gry grn hzl oth.
pid (Passport ID) - a nine-digit number, including leading zeroes.
cid (Country ID) - ignored, missing or not.

Your job is to count the passports where all required fields are both present and valid according to the above rules.
Here are some example values:

byr valid:   2002
byr invalid: 2003

hgt valid:   60in
hgt valid:   190cm
hgt invalid: 190in
hgt invalid: 190

hcl valid:   #123abc
hcl invalid: #123abz
hcl invalid: 123abc

ecl valid:   brn
ecl invalid: wat

pid valid:   000000001
pid invalid: 0123456789
Here are some invalid passports:

eyr:1972 cid:100
hcl:#18171d ecl:amb hgt:170 pid:186cm iyr:2018 byr:1926

iyr:2019
hcl:#602927 eyr:1967 hgt:170cm
ecl:grn pid:012533040 byr:1946

hcl:dab227 iyr:2012
ecl:brn hgt:182cm pid:021572410 eyr:2020 byr:1992 cid:277

hgt:59cm ecl:zzz
eyr:2038 hcl:74454a iyr:2023
pid:3556412378 byr:2007
Here are some valid passports:

pid:087499704 hgt:74in ecl:grn iyr:2012 eyr:2030 byr:1980
hcl:#623a2f

eyr:2029 ecl:blu cid:129 byr:1989
iyr:2014 pid:896056539 hcl:#a97842 hgt:165cm

hcl:#888785
hgt:164cm byr:2001 iyr:2015 cid:88
pid:545766238 ecl:hzl
eyr:2022

iyr:2010 hgt:158cm hcl:#b6652a ecl:blu byr:1944 eyr:2021 pid:093154719

Count the number of valid passports - those that have all required fields and valid values. Continue to treat cid as
optional. In your batch file, how many passports are valid?
*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

//type Passport map[string]string
type PassportResult struct {
	passport    sync.Map
	valid       bool
	validStrict bool
}

func main() {

	input, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var passports []sync.Map
	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	for {
		passport, eof := scanPassport(scanner)
		passports = append(passports, passport)
		if eof {
			break
		}
	}

	var wg sync.WaitGroup
	passportChannel := make(chan PassportResult, len(passports))
	wg.Add(len(passports))
	for _, passport := range passports {
		go check(passportChannel, passport, &wg)
	}

	go func() {
		wg.Wait()
		close(passportChannel)
	}()

	count := 0
	valid := 0
	strict := 0
	for result := range passportChannel {
		count++
		if result.valid {
			valid++
			if result.validStrict {
				strict++
			}
		}
	}

	fmt.Printf("Count: %d, Total valid: %d, Total strictly valid: %d\n", count, valid, strict)
}

func check(prc chan PassportResult, passport sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	type valueValidation func(string) bool
	type fieldReqs struct {
		field     string
		validator valueValidation
	}

	// this got bloated...
	requiredFields := []fieldReqs{
		{"byr", func(s string) bool {
			i, err := strconv.Atoi(s)
			if err != nil {
				return false
			}
			return i >= 1920 && i <= 2020
		}},
		{"iyr", func(s string) bool {
			i, err := strconv.Atoi(s)
			if err != nil {
				return false
			}
			return i >= 2010 && i <= 2020
		}},
		{"eyr", func(s string) bool {
			i, err := strconv.Atoi(s)
			if err != nil {
				return false
			}
			return i >= 2020 && i <= 2030
		}},
		{"hgt", func(s string) bool {
			r := regexp.MustCompile("^[0-9]+(cm|in)$")
			p := r.FindStringSubmatch(s)
			if len(p) != 2 {
				return false
			}
			i, err := strconv.Atoi(p[0][:len(p[0])-len(p[1])])
			if err != nil {
				return false
			}
			if p[1] == "in" {
				return i >= 59 && i <= 76

			} else if p[1] == "cm" {
				return i >= 150 && i <= 193
			} else {
				return false
			}

		}},
		{"hcl", func(s string) bool {
			matched, err := regexp.MatchString("^#[0-9a-f]{6}$", s)
			if !matched || err != nil {
				return false
			}
			return true
		}},
		{"ecl", func(s string) bool {
			matched, err := regexp.MatchString("^(amb|blu|brn|gry|grn|hzl|oth)$", s)
			if !matched || err != nil {
				return false
			}
			return true
		}},
		{"pid", func(s string) bool {
			matched, err := regexp.MatchString("^[0-9]{9}$", s)
			if !matched || err != nil {
				return false
			}
			return true
		}},
		//		{"cid", func(s string) bool { return true }},
	}

	validStrict := true

	for _, req := range requiredFields {
		valInt, ok := passport.Load(req.field)
		if !ok {
			prc <- PassportResult{passport, false, false}
			return
		}
		validStrict = validStrict && req.validator(valInt.(string))
		/*		if !validStrict {
					panic(valInt.(string))
				}
		*/
	}

	prc <- PassportResult{passport, true, validStrict}

}

func scanPassport(scanner *bufio.Scanner) (sync.Map, bool) {
	var passport sync.Map
	var eof bool
	var err error
	for {
		eof = !scanner.Scan()
		if eof {
			break
		}
		line := scanner.Text()
		if strings.Trim(line, " ") == "" {
			break
		}
		passport, err = buildPassport(passport, parsePassportLine(line))
		if err != nil {
			panic(err)
		}
	}
	return passport, eof
}

func buildPassport(passport sync.Map, data sync.Map) (sync.Map, error) {
	data.Range(func(key, value interface{}) bool {
		_, ok := passport.Load(key)
		if ok {
			panic(fmt.Errorf("Duplicate key in passport: %s", key))
		}
		passport.Store(key, value)
		return true
	})

	return passport, nil
}

func parsePassportLine(line string) sync.Map {
	var p sync.Map
	attribs := strings.Split(line, " ")
	for _, attribStr := range attribs {
		attrib := strings.Split(attribStr, ":")
		p.Store(attrib[0], attrib[1])
	}
	return p
}
