/*
	2048 game in Golang

	2021-02-11 Author: Miigon
*/

package main

import (
	"fmt"
	"math"
	"math/rand"
)

const boardSize = 4
const chanceToGet4Divider = 8 // the chance to get a 4 instead of a 2 is 1/8

// for Windows (or any environment that doesn't support color), set this to false
const enableColor = true

type squareNumber int

type gameState struct {
	board [boardSize][boardSize]squareNumber
	score int
}

var colorPalette = [...]string{"0", "1", "2", "4", "5", "3", "6", "1;1", "2;1", "4;1", "5;1", "3;1"}

func (state *gameState) drawBoard() {
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			num := state.board[i][j]
			if num == 0 {
				fmt.Print("    ")
			} else {
				if enableColor {
					fmt.Printf("\033[3%sm%4d\033[0m", colorPalette[int(math.Log2(float64(num)))], num)
				} else {
					fmt.Printf("%4d", num)
				}
			}
			if j != boardSize-1 {
				fmt.Print("|")
			}
		}
		fmt.Print("\n")
		if i != boardSize-1 {
			for j := 0; j < boardSize; j++ {
				if j != boardSize-1 {
					fmt.Print("-----")
				} else {
					fmt.Print("----")
				}
			}
			fmt.Print("\n")
		}
	}
}

type direction int

const (
	up = iota
	right
	down
	left
)

/*
	for every turn:
	(assuming right arrow key was pressed)
	for every row, from right to left, search for all the "mergable" pairs (adjacent numbers that are the same)
	for every pair found, add their numbers together, then place them at the right side of the board, one by one.
*/

func (state *gameState) move(dir direction) {
	for i := 0; i < boardSize; i++ {
		var lastNonZeroRef *squareNumber = nil
		var emptyJRef *squareNumber
		emptyJ := 0
		for j := 0; j < boardSize; j++ {
			var numRef *squareNumber
			switch dir {
			case up:
				numRef = &state.board[j][i]
				emptyJRef = &state.board[emptyJ][i]
			case down:
				numRef = &state.board[boardSize-j-1][i]
				emptyJRef = &state.board[boardSize-emptyJ-1][i]
			case left:
				numRef = &state.board[i][j]
				emptyJRef = &state.board[i][emptyJ]
			case right:
				numRef = &state.board[i][boardSize-j-1]
				emptyJRef = &state.board[i][boardSize-emptyJ-1]
			}
			num := *numRef
			if num == 0 {
				continue
			}
			if lastNonZeroRef != nil && *lastNonZeroRef == num {
				// merge adjacent equal numbers
				*lastNonZeroRef = num * 2
				state.score += int(num * 2)
				num, *numRef = 0, 0
				lastNonZeroRef = nil
			} else {
				// move number to the side
				if num != 0 {
					*numRef = 0
					*emptyJRef = num
					lastNonZeroRef = emptyJRef
					emptyJ++
				} else {
					lastNonZeroRef = numRef
				}
			}
		}
	}
}

func (state *gameState) randomlyPlaceNewNum() {
	// randomly place a new number anywhere on the board
	for true {
		i, j := rand.Intn(boardSize), rand.Intn(boardSize)
		if state.board[i][j] == 0 {
			if rand.Intn(chanceToGet4Divider) == 0 { // chance to get a 4
				state.board[i][j] = 4
			} else {
				state.board[i][j] = 2
			}
			break
		}
	}
}

func main() {
	var state gameState
	state.randomlyPlaceNewNum()
	for true {
		state.randomlyPlaceNewNum()
		fmt.Printf("Score: %d\n", state.score)
		state.drawBoard()

		var operation byte
		fmt.Scanf("%c\n", &operation)
		switch operation {
		case 'w':
			state.move(up)
		case 's':
			state.move(down)
		case 'a':
			state.move(left)
		case 'd':
			state.move(right)
		}
	}

}
