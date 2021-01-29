/*
	Tic-tac-toe game in Golang

	2021-01-29 Author: Miigon
	https://blog.miigon.ml/posts/golang-project-tic-tac-toe/
*/

package main

import (
	"fmt"
)

// what piece is on a certain square
type squareState int
type player int

const (
	none   = iota
	cross  = iota
	circle = iota
)

func (e player) String() string {
	switch e {
	case none:
		return "none"
	case cross:
		return "cross"
	case circle:
		return "circle"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

// current state of the game
type gameState struct {
	board      [3][3]squareState
	turnPlayer player
}

// define a method for struct type `gameState`
func (state *gameState) drawBoard() {
	// (# challenge: use Stringer to simplify this function)
	for i, row := range state.board {
		for j, square := range row {
			fmt.Print(" ")
			switch square {
			case none:
				fmt.Print(" ")
			case cross:
				fmt.Print("X")
			case circle:
				fmt.Print("O")
			}
			if j != len(row)-1 {
				fmt.Print(" |")
			}
		}
		if i != len(state.board)-1 {
			fmt.Print("\n------------")
		}
		fmt.Print("\n")
	}
}

type markAlreadyExistError struct {
	row    int
	column int
}

type positionOutOfBoundError struct {
	row    int
	column int
}

func (e *markAlreadyExistError) Error() string {
	return fmt.Sprintf("position (%d,%d) already has a mark on it.", e.row, e.column)
}

func (e *positionOutOfBoundError) Error() string {
	return fmt.Sprintf("position (%d,%d) is out of bound.", e.row, e.column)
}

func (state *gameState) placeMark(row int, column int) error {
	if row < 0 || column < 0 || row >= len(state.board) || column >= len(state.board[row]) {
		return &positionOutOfBoundError{row, column}
	}
	if state.board[row][column] != none {
		return &markAlreadyExistError{row, column}
	}

	state.board[row][column] = squareState(state.turnPlayer)
	return nil // no error
}

type gameResult int

const (
	noWinnerYet = iota
	crossWon
	circleWon
	draw
)

func (state *gameState) whosNext() player {
	return state.turnPlayer
}

func (state *gameState) nextTurn() {
	if state.turnPlayer == cross {
		state.turnPlayer = circle
	} else {
		state.turnPlayer = cross
	}
}

func (state *gameState) checkForWinner() gameResult {
	boardSize := len(state.board) // assuming the board is always square-shaped.

	// define a lambda function for checking one line
	checkLine := func(startRow int, startColumn int, deltaRow int, deltaColumn int) gameResult {
		var lastSquare squareState = state.board[startRow][startColumn]
		row, column := startRow+deltaRow, startColumn+deltaColumn

		// loop starts from the second square(startRow + deltaRow, startColumn + deltaColumn)
		for row >= 0 && column >= 0 && row < boardSize && column < boardSize {

			// there can't be a winner if a empty square is present within the line
			if state.board[row][column] == none {
				return noWinnerYet
			}

			if lastSquare != state.board[row][column] {
				return noWinnerYet
			}

			lastSquare = state.board[row][column]
			row, column = row+deltaRow, column+deltaColumn
		}

		// someone has won the game
		if lastSquare == cross {
			return crossWon
		} else if lastSquare == circle {
			return circleWon
		}

		return noWinnerYet
	}

	// check horizontal rows
	for row := 0; row < boardSize; row++ {
		if result := checkLine(row, 0, 0, 1); result != noWinnerYet {
			return result
		}
	}
	// check vertical columns
	for column := 0; column < boardSize; column++ {
		if result := checkLine(column, 0, 0, 1); result != noWinnerYet {
			return result
		}
	}
	// check top-left to bottom-right diagonal
	if result := checkLine(0, 0, 1, 1); result != noWinnerYet {
		return result
	}
	// check top-right to bottom-left diagonal
	if result := checkLine(0, boardSize-1, 1, -1); result != noWinnerYet {
		return result
	}
	// check for draw
	for _, row := range state.board {
		for _, square := range row {
			if square == none {
				return noWinnerYet
			}
		}
	}
	// if no one wins yet, but none of the squares are empty
	return draw
}

func main() {
	state := gameState{}
	state.turnPlayer = cross // cross goes first

	var result gameResult = noWinnerYet

	// the main game loop
	for {
		fmt.Printf("next player to place a mark is: %v\n", state.whosNext())

		// 1. draw the board onto the screen
		state.drawBoard()

		fmt.Printf("where to place a %v? (input row then column, separated by space)\n> ", state.whosNext())

		// 2. use a loop to take input
		for {
			var row, column int
			fmt.Scan(&row, &column)

			e := state.placeMark(row-1, column-1) // -1 so coordinate starts at (1,1) instead of (0,0)

			// if a valid position was entered, break out from the input loop
			if e == nil {
				break
			}

			// if an invalid position was entered, prompt the player to re-enter another position
			fmt.Println(e)
			fmt.Printf("please re-enter a position:\n> ")
		}

		// 3. check if anyone has won the game
		result = state.checkForWinner()
		if result != noWinnerYet {
			break
		}

		// 4. if no one has won in this turn, go on for next turn and continue the game loop
		state.nextTurn()

		fmt.Println()
	}

	state.drawBoard()

	switch result {
	case crossWon:
		fmt.Printf("cross won the game!\n")
	case circleWon:
		fmt.Printf("circle won the game!\n")
	case draw:
		fmt.Printf("the game has ended with a draw!\n")
	}
}
