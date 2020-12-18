/*
--- Day 12: Rain Risk ---
Your ferry made decent progress toward the island, but the storm came in faster than anyone expected. The ferry needs
to take evasive actions!

Unfortunately, the ship's navigation computer seems to be malfunctioning; rather than giving a route directly to safety,
it produced extremely circuitous instructions. When the captain uses the PA system to ask if anyone can help, you quickly
volunteer.

The navigation instructions (your puzzle input) consists of a sequence of single-character actions paired with integer
input values. After staring at them for a few minutes, you work out what they probably mean:

Action N means to move north by the given value.
Action S means to move south by the given value.
Action E means to move east by the given value.
Action W means to move west by the given value.
Action L means to turn left the given number of degrees.
Action R means to turn right the given number of degrees.
Action F means to move forward by the given value in the direction the ship is currently facing.

The ship starts by facing east. Only the L and R actions change the direction the ship is facing. (That is, if the
ship is facing east and the next instruction is N10, the ship would move north 10 units, but would still move east if
the following action were F.)

For example:

F10
N3
F7
R90
F11

These instructions would be handled as follows:

F10 would move the ship 10 units east (because the ship starts by facing east) to east 10, north 0.
N3 would move the ship 3 units north to east 10, north 3.
F7 would move the ship another 7 units east (because the ship is still facing east) to east 17, north 3.
R90 would cause the ship to turn right by 90 degrees and face south; it remains at east 17, north 3.
F11 would move the ship 11 units south to east 17, south 8.

At the end of these instructions, the ship's Manhattan distance (sum of the absolute values of its east/west position
	and its north/south position) from its starting position is 17 + 8 = 25.

Figure out where the navigation instructions lead. What is the Manhattan distance between that location and the ship's
starting position?

--- Part Two ---
Before you can give the destination to the captain, you realize that the actual action meanings were printed on
the back of the instructions the whole time.

Almost all of the actions indicate how to move a waypoint which is relative to the ship's position:

Action N means to move the waypoint north by the given value.
Action S means to move the waypoint south by the given value.
Action E means to move the waypoint east by the given value.
Action W means to move the waypoint west by the given value.
Action L means to rotate the waypoint around the ship left (counter-clockwise) the given number of degrees.
Action R means to rotate the waypoint around the ship right (clockwise) the given number of degrees.
Action F means to move forward to the waypoint a number of times equal to the given value.

The waypoint starts 10 units east and 1 unit north relative to the ship. The waypoint is relative to the ship;
that is, if the ship moves, the waypoint moves with it.

For example, using the same instructions as above:

F10 moves the ship to the waypoint 10 times (a total of 100 units east and 10 units north), leaving the ship at
east 100, north 10. The waypoint stays 10 units east and 1 unit north of the ship.
N3 moves the waypoint 3 units north to 10 units east and 4 units north of the ship. The ship remains at east 100,
north 10.
F7 moves the ship to the waypoint 7 times (a total of 70 units east and 28 units north), leaving the ship at east
170, north 38. The waypoint stays 10 units east and 4 units north of the ship.
R90 rotates the waypoint around the ship clockwise 90 degrees, moving it to 4 units east and 10 units south of the
ship. The ship remains at east 170, north 38.
F11 moves the ship to the waypoint 11 times (a total of 44 units east and 110 units south), leaving the ship at east
214, south 72. The waypoint stays 4 units east and 10 units south of the ship.
After these operations, the ship's Manhattan distance from its starting position is 214 + 72 = 286.

Figure out where the navigation instructions actually lead. What is the Manhattan distance between that location and
the ship's starting position?

*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// It's very important, especially for the cardinal directions, that they keep this order
// It seems the absence of enums is not made up for by a greater, more progressive concept in golang..
const (
	NORTH int = iota
	EAST
	SOUTH
	WEST
	LEFT
	RIGHT
	FORWARD
)

type NavInstruction struct {
	op   int
	dist int
}

type Ship struct {
	facing int
	x      int
	y      int
	wpx    int
	wpy    int
}

func main() {

	input, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	shipNaieve := Ship{EAST, 0, 0, 10, 1}
	ship := Ship{EAST, 0, 0, 10, 1}
	var instructions []NavInstruction
	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	for i := 0; scanner.Scan(); i++ {
		inst := parseInstruction(scanner.Text())
		instructions = append(instructions, inst)

		shipNaieve.runOperation(inst)
		ship.runCorrectOperation(inst)

	}

	fmt.Printf("Naive Manhattan distance: %d\n", shipNaieve.manhattanDistance())
	fmt.Printf("Correct Manhattan distance: %d\n", ship.manhattanDistance())

}

func (ship *Ship) runCorrectOperation(inst NavInstruction) {
	switch inst.op {
	case NORTH, EAST, SOUTH, WEST:
		ship.moveWPNESW(inst.op, inst.dist)
	case FORWARD:
		ship.moveToWP(inst.dist)
	case LEFT, RIGHT:
		ship.rotateWP(inst.op, inst.dist)

	}
}

func (ship *Ship) runOperation(inst NavInstruction) {
	switch inst.op {
	case NORTH, EAST, SOUTH, WEST:
		ship.moveNESW(inst.op, inst.dist)
	case FORWARD:
		ship.moveNESW(ship.facing, inst.dist)
	case LEFT, RIGHT:
		ship.rotateShip(inst.op, inst.dist)

	}

}

// take this function...
// multiply it by all the other little fundamental utilities out there (floor, ceil, etc)
// multiply that by all the developers and orgs that will roll their own fundamental
// 		utilties when basic tools are missing (all devs & orgs?)
// apply a coefficient for bugs in all those (millions?) lines of code
// apply a separate coefficient for performance inefficienies
// apply yet another coefficient for security gaps
// take a moment to ruminate on the cost implications (for language users at large) of
// language owners not creating/supporting a robust standard lib
func GoShouldIncludeThisInItsStandardLibs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func (ship *Ship) manhattanDistance() int {
	return GoShouldIncludeThisInItsStandardLibs(ship.x) + GoShouldIncludeThisInItsStandardLibs(ship.y)
}

func (ship *Ship) moveWPNESW(dir int, dist int) {
	switch dir {
	case NORTH:
		ship.wpy += dist
	case SOUTH:
		ship.wpy -= dist
	case EAST:
		ship.wpx += dist
	case WEST:
		ship.wpx -= dist
	default:
		panic("Bad Direction")
	}
}

func (ship *Ship) moveToWP(dist int) {
	ship.x += ship.wpx * dist
	ship.y += ship.wpy * dist
}

func (ship *Ship) moveNESW(dir int, dist int) {
	switch dir {
	case NORTH:
		ship.y += dist
	case SOUTH:
		ship.y -= dist
	case EAST:
		ship.x += dist
	case WEST:
		ship.x -= dist
	default:
		panic("Bad Direction")
	}
}

// This could be much more gnarly to do it right, but given everything
// in the input data are multiples of 90 (or 90, 180, 270 to be explicit)
// we can cheat by just rotating values
func (ship *Ship) rotateWP(dir int, deg int) {
	for deg > 0 {
		// sanity check...
		// rotating (5,1) CW:
		// (5, 1) (1, -5) (-5, -1) (-1, 5)
		// -> (y, -x)  (y, -x)  (y, -x)
		// rotating (5,1) CCW:
		// (5, 1) (-1, 5) (-5, -1) (1, -5)
		// -> (-y, x)  (-y, x)  (-y, x)
		if dir == RIGHT { // CW
			tmp := ship.wpy
			ship.wpy = -ship.wpx
			ship.wpx = tmp
		} else if dir == LEFT { // CCW
			tmp := ship.wpx
			ship.wpx = -ship.wpy
			ship.wpy = tmp
		} else {
			panic("Bad Rotation")
		}

		deg -= 90
	}

}

func (ship *Ship) rotateShip(dir int, deg int) {
	compassCW := [4]int{NORTH, EAST, SOUTH, WEST}
	steps := int(deg / 90)

	if dir == RIGHT { // Clockwise
		ship.facing = compassCW[(ship.facing+steps)%4]
	} else if dir == LEFT { // CCW
		ship.facing = compassCW[(ship.facing+(3*steps))%4]
	} else {
		panic("Bad Rotation")
	}
}

func parseInstruction(s string) NavInstruction {
	var op int
	switch s[:1] {
	case "N":
		op = NORTH
	case "E":
		op = EAST
	case "S":
		op = SOUTH
	case "W":
		op = WEST
	case "L":
		op = LEFT
	case "R":
		op = RIGHT
	case "F":
		op = FORWARD
	default:
		panic("Bad insruction parse")
	}
	arg, err := strconv.Atoi(s[1:])
	if err != nil {
		panic(fmt.Errorf("Bad instruction parse: %s", s))
	}
	return NavInstruction{op, arg}
}
