/*
--- Day 11: Seating System ---
Your plane lands with plenty of time to spare. The final leg of your journey is a ferry that goes directly to the
tropical island where you can finally start your vacation. As you reach the waiting area to board the ferry, you
realize you're so early, nobody else has even arrived yet!

By modeling the process people use to choose (or abandon) their seat in the waiting area, you're pretty sure you can
predict the best place to sit. You make a quick map of the seat layout (your puzzle input).

The seat layout fits neatly on a grid. Each position is either floor (.), an empty seat (L), or an occupied seat (#).
For example, the initial seat layout might look like this:

L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL

Now, you just need to model the people who will be arriving shortly. Fortunately, people are entirely predictable
and always follow a simple set of rules. All decisions are based on the number of occupied seats adjacent to a given
seat (one of the eight positions immediately up, down, left, right, or diagonal from the seat). The following rules are
applied to every seat simultaneously:

If a seat is empty (L) and there are no occupied seats adjacent to it, the seat becomes occupied.
If a seat is occupied (#) and four or more seats adjacent to it are also occupied, the seat becomes empty.
Otherwise, the seat's state does not change.
Floor (.) never changes; seats don't move, and nobody sits on the floor.

After one round of these rules, every seat in the example layout becomes occupied:

#.##.##.##
#######.##
#.#.#..#..
####.##.##
#.##.##.##
#.#####.##
..#.#.....
##########
#.######.#
#.#####.##

After a second round, the seats with four or more occupied adjacent seats become empty again:

#.LL.L#.##
#LLLLLL.L#
L.L.L..L..
#LLL.LL.L#
#.LL.LL.LL
#.LLLL#.##
..L.L.....
#LLLLLLLL#
#.LLLLLL.L
#.#LLLL.##

This process continues for three more rounds:

#.##.L#.##
#L###LL.L#
L.#.#..#..
#L##.##.L#
#.##.LL.LL
#.###L#.##
..#.#.....
#L######L#
#.LL###L.L
#.#L###.##
#.#L.L#.##
#LLL#LL.L#
L.L.L..#..
#LLL.##.L#
#.LL.LL.LL
#.LL#L#.##
..L.L.....
#L#LLLL#L#
#.LLLLLL.L
#.#L#L#.##
#.#L.L#.##
#LLL#LL.L#
L.#.L..#..
#L##.##.L#
#.#L.LL.LL
#.#L#L#.##
..L.L.....
#L#L##L#L#
#.LLLLLL.L
#.#L#L#.##

At this point, something interesting happens: the chaos stabilizes and further applications of these rules cause no
seats to change state! Once people stop moving around, you count 37 occupied seats.

Simulate your seating area by applying the seating rules repeatedly until no seats change state. How many seats end
up occupied?

--- Part Two ---
As soon as people start to arrive, you realize your mistake. People don't just care about adjacent seats - they care
about the first seat they can see in each of those eight directions!

Now, instead of considering just the eight immediately adjacent seats, consider the first seat in each of those
eight directions. For example, the empty seat below would see eight occupied seats:

.......#.
...#.....
.#.......
.........
..#L....#
....#....
.........
#........
...#.....

The leftmost empty seat below would only see one empty seat, but cannot see any of the occupied ones:

.............
.L.L.#.#.#.#.
.............

The empty seat below would see no occupied seats:

.##.##.
#.#.#.#
##...##
...L...
##...##
#.#.#.#
.##.##.

Also, people seem to be more tolerant than you expected: it now takes five or more visible occupied seats for an
occupied seat to become empty (rather than four or more from the previous rules). The other rules still apply:
empty seats that see no occupied seats become occupied, seats matching no rule don't change, and floor never changes.

Given the same starting layout as above, these new rules cause the seating area to shift around as follows:

L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL

#.##.##.##
#######.##
#.#.#..#..
####.##.##
#.##.##.##
#.#####.##
..#.#.....
##########
#.######.#
#.#####.##

#.LL.LL.L#
#LLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLL#
#.LLLLLL.L
#.LLLLL.L#

#.L#.##.L#
#L#####.LL
L.#.#..#..
##L#.##.##
#.##.#L.##
#.#####.#L
..#.#.....
LLL####LL#
#.L#####.L
#.L####.L#

#.L#.L#.L#
#LLLLLL.LL
L.L.L..#..
##LL.LL.L#
L.LL.LL.L#
#.LLLLL.LL
..L.L.....
LLLLLLLLL#
#.LLLLL#.L
#.L#LL#.L#

#.L#.L#.L#
#LLLLLL.LL
L.L.L..#..
##L#.#L.L#
L.L#.#L.L#
#.L####.LL
..#.#.....
LLL###LLL#
#.LLLLL#.L
#.L#LL#.L#

#.L#.L#.L#
#LLLLLL.LL
L.L.L..#..
##L#.#L.L#
L.L#.LL.L#
#.LLLL#.LL
..#.L.....
LLL###LLL#
#.LLLLL#.L
#.L#LL#.L#

Again, at this point, people stop shifting around and the seating area reaches equilibrium. Once this occurs,
you count 26 occupied seats.

Given the new visibility method and the rule change for occupied seats becoming empty, once equilibrium is
reached, how many seats end up occupied?
*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

// honestly just felt like i'm using hand-wavy array slices too much, so doing a bit of this for variation
const COLUMNS = 98 // test data is 10
const NAIEVE = false

type SeatStatus int

const (
	NOSEAT SeatStatus = iota
	VACANT
	OCCUPIED
)

type Seat struct {
	status     SeatStatus
	nextStatus SeatStatus
}

type coord struct {
	row int
	col int
}

type SeatArrangement [][COLUMNS]Seat

func main() {

	input, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var seats SeatArrangement
	var seatData []string
	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	for i := 0; scanner.Scan(); i++ {
		row := scanner.Text()
		seatData = append(seatData, row)
		var seatRow [COLUMNS]Seat
		for j, c := range row {
			status := getStatus(c)
			seatRow[j].status = status
			seatRow[j].nextStatus = status
		}
		seats = append(seats, seatRow)
	}

	for runLife(&seats) {
		for i := range seats {
			for j := 0; j < COLUMNS; j++ {
				seats[i][j].status = seats[i][j].nextStatus
			}
		}
	}

	// this could be cleaner like keeping track as we go
	occupants := 0
	for i := range seats {
		for j := 0; j < COLUMNS; j++ {
			if seats[i][j].status == OCCUPIED {
				occupants++
			}
		}
	}

	fmt.Printf("occupants: %d\n", occupants)
}

func runLife(seats *SeatArrangement) bool {
	var swg sync.WaitGroup
	changeChan := make(chan bool, len(*seats)*COLUMNS)
	doneChan := make(chan bool)
	changed := false

	// top level wait group blocks separately and listens for a worker to say its value changed
	var tlwg sync.WaitGroup
	tlwg.Add(1)
	go func() {
		for {
			select {
			case <-changeChan:
				changed = true
			case <-doneChan:
				tlwg.Done()
				return
			}
		}
	}()

	// run through the seats and spawn a worker to figure out the next state per seat
	// workers will report via changeChan if the value changes
	for i := range *seats {
		for j := 0; j < COLUMNS; j++ {
			swg.Add(1)
			go updateNextStatus(seats, i, j, changeChan, &swg)

		}
	}
	// wait for everyone to finish and communicate when it's done
	swg.Wait()
	doneChan <- true

	tlwg.Wait()

	//fmt.Println(changed)

	return changed

}

// If a seat is empty (L) and there are no occupied seats adjacent to it, the seat becomes occupied.
// If a seat is occupied (#) and four or more seats adjacent to it are also occupied, the seat becomes empty.
// Otherwise, the seat's state does not change.
func updateNextStatus(seats *SeatArrangement, row int, col int, changed chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	tolerance := 4

	if !NAIEVE {
		tolerance = 5
	}

	neighbors := countSeatedNeighbors(seats, row, col)

	if (*seats)[row][col].status == VACANT && neighbors == 0 {
		(*seats)[row][col].nextStatus = OCCUPIED
		changed <- true
	} else if (*seats)[row][col].status == OCCUPIED && neighbors >= tolerance {
		(*seats)[row][col].nextStatus = VACANT
		changed <- true
	}

}

func updateStatus(seats *SeatArrangement, row int, col int) {
	(*seats)[row][col].status = (*seats)[row][col].nextStatus
}

func getCoordDir(orig coord, dir coord) coord {
	row := orig.row
	col := orig.col
	if dir.row > 0 {
		row++
	} else if dir.row < 0 {
		row--
	}

	if dir.col > 0 {
		col++
	} else if dir.col < 0 {
		col--
	}

	return coord{row, col}
}

func isDirSeatOccupied(seats *SeatArrangement, orig coord, dir coord, adjacent bool) bool {
	cur := orig
	for {
		cur = getCoordDir(cur, dir)
		if cur.row < 0 || cur.row >= len(*seats) || cur.col < 0 || cur.col >= COLUMNS {
			break
		}

		if (*seats)[cur.row][cur.col].status == OCCUPIED {
			return true
		} else if (*seats)[cur.row][cur.col].status == VACANT {
			return false
		}

		// bail after one go for adjacency checks
		if adjacent {
			break
		}
	}

	return false
}

func countSeatedNeighbors(seats *SeatArrangement, row int, col int) int {
	dirs := []coord{
		{1, -1}, {1, 0}, {1, 1},
		{0, -1}, {0, 1},
		{-1, -1}, {-1, 0}, {-1, 1},
	}

	count := 0
	countChan := make(chan bool, 8)

	go func() {
		var wg sync.WaitGroup
		wg.Add(8)

		for _, dir := range dirs {
			go func(d coord) {
				defer wg.Done()
				countChan <- isDirSeatOccupied(seats, coord{row, col}, d, NAIEVE)
			}(dir)
		}

		wg.Wait()
		close(countChan)
	}()

	for result := range countChan {
		if result {
			count++
		}
	}

	return count
}

func countSeatedNeighborsNaieve(seats *SeatArrangement, row int, col int) int {
	count := 0
	for i := row - 1; i <= (row + 1); i++ {
		if (i < 0) || (i >= len(*seats)) {
			continue
		}
		for j := col - 1; j <= (col + 1); j++ {
			if (j < 0) || (j >= COLUMNS) || (i == row && j == col) {
				continue
			}
			if (*seats)[i][j].status == OCCUPIED {
				count++
			}
		}
	}

	return count
}

func getStatus(seat rune) SeatStatus {

	switch seat {
	case '.':
		return NOSEAT
	case 'L':
		return VACANT
	case '#':
		return OCCUPIED
	default:
		panic(fmt.Errorf("no seat status for %s", seat))
	}
}

func printSeats(seats *SeatArrangement) {
	return
	for i := range *seats {
		for j := 0; j < COLUMNS; j++ {
			fmt.Printf("%c", getSeatRune((*seats)[i][j].status))
		}
		fmt.Println("")
	}
	fmt.Println("----------------------\n")
}

func getSeatRune(ss SeatStatus) rune {
	switch ss {
	case NOSEAT:
		return '.'
	case VACANT:
		return 'L'
	case OCCUPIED:
		return '#'
	default:
		return '?'
	}
}
