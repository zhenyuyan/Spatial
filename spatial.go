// Zhenyu Yang 2017/10/11

package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"
)

// The data stored in a single cell of a field
type Cell struct {
	strategy string  //represents "C" or "D" corresponding to the type of prisoner in the cell
	score    float64 //represents the score of the cell based on the prisoner's relationship with neighboring cells
}

// The game board is a 2D slice of Cell objects
type GameBoard [][]Cell

func ReadFile(filename string) []string {
	// open the file and make sure all went well
	in, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: couldn’t open the file")
		os.Exit(1)
	}
	// create the variable to hold the lines
	var lines []string = make([]string, 0)
	// for every line in the file
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		// append it to the lines slice
		lines = append(lines, scanner.Text())
	}
	// check that all went ok
	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}
	// close the file and return the lines
	in.Close()
	return lines
}

func IsInTheField(i, j, width, length int) bool { // to check whether a node is in field
	if i < 0 || i >= width || j < 0 || j >= length {
		// Return false when the a square is out of range
		return false
	}
	return true
}

// Calculate each prisoner's score
func CalculateScore(r GameBoard, i int, j int, b float64, width int, length int) float64 {
	count := 0.0
	for k := i - 1; k <= i+1; k++ {
		for l := j - 1; l <= j+1; l++ {
			if IsInTheField(k, l, width, length) {
				// to check each circumstances correct. Use the criteria in the provided write up
				if r[i][j].strategy == "C" && r[k][l].strategy == "C" {
					count += 1.0
				}
				if r[i][j].strategy == "C" && r[k][l].strategy == "D" {
					count += 0.0
				}
				if r[i][j].strategy == "D" && r[k][l].strategy == "C" {
					count += b
				}
				if r[i][j].strategy == "D" && r[k][l].strategy == "D" {
					count += 0.0
				}
			}
		}
	}
	if r[i][j].strategy == "C" {
		count = count - 1.0
	}
	return count
}

func FindMaxCordinate(r GameBoard, b float64, i, j, width, length int) string {
	//find the make r.score among r[i][j].score
	var e1, f1 int
	max := 0.0
	for p := i - 1; p <= i+1; p++ {
		for q := j - 1; q <= j+1; q++ {
			if IsInTheField(p, q, width, length) == true {
				if max <= r[p][q].score {
					max = r[p][q].score
					e1 = p
					f1 = q
				}
			}
		}
	}
	return r[e1][f1].strategy
}

func EvolveStep(r GameBoard, b float64, width int, length int) GameBoard {
	// To update strategy that is better for a prisoner to higher scores
	// initiate r.score. Copy the score into an empty slice r1. Put the 9 scores
	//into an array
	for p := 0; p < width; p++ {
		for q := 0; q < length; q++ {
			r[p][q].score = CalculateScore(r, p, q, b, width, length)
		}
	}
	r1 := make(GameBoard, length)
	for k := range r1 {
		r1[k] = make([]Cell, width)
	}

	for i := 0; i < width; i++ {
		for j := 0; j < length; j++ {
			//Find the max coordinate for every prisoners around
			r1[i][j].strategy = FindMaxCordinate(r, b, i, j, width, length)
		}
	}

	return r1
}

func DrawBoard(board GameBoard, width int, length int) Canvas {
	// set all the parameters
	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)
	pic := CreateNewCanvas(1000, 1000)
	// Set the orignal coordinate
	x1 := 0
	x2 := 1000 / width
	y1 := 0
	y2 := 1000 / length
	// xcell and ycell stands for a single unit of a square
	xcell := 1000 / width
	ycell := 1000 / length
	// for every square in canvas, fill the related color
	for i := 0; i < width; i++ {
		for j := 0; j < length; j++ {
			if board[i][j].strategy == "D" {
				pic.SetFillColor(red)
				pic.ClearRect(x1, y1, x2, y2)
				y1 = y1 + ycell
				y2 = y2 + ycell
			} else if board[i][j].strategy == "C" {
				pic.SetFillColor(blue)
				pic.ClearRect(x1, y1, x2, y2)
				y1 = y1 + ycell
				y2 = y2 + ycell
			}
		}
		// Move to another square that is waiting to update
		y1 = 0
		x1 = x1 + xcell
		y2 = ycell
		x2 = x2 + xcell
	}
	return pic
}

func Animation(board GameBoard, b float64, steps int, width int, length int) []image.Image {
	// Set empty gifImages slices
	gifImages := make([]image.Image, steps)
	pic := DrawBoard(board, width, length)
	// Update the board each time
	for i := 0; i < steps; i++ {
		board1 := EvolveStep(board, b, width, length)
		board = board1
		pic = DrawBoard(board1, width, length)
		gifImages[i] = pic.img

	}
	// Save the final picture
	pic.SaveToPNG("Prisoners.png")
	return gifImages
}

func main() {
	// Set all the parameters in commandline
	lines := ReadFile(os.Args[1])
	b, _ := strconv.ParseFloat(os.Args[2], 64)
	steps, _ := strconv.Atoi(os.Args[3])
	// Read the first line containing two numbers indicating width and length
	var items1 []string = strings.Split(lines[0], " ")
	width, _ := strconv.Atoi(items1[1])
	length, _ := strconv.Atoi(items1[0])
	// make a empty 2-D board
	board := make(GameBoard, length)
	for k := range board {
		board[k] = make([]Cell, width)
	}
	// insert all the elements from the file into board
	for i := 1; i <= length; i++ {
		var items2 []string = strings.Split(lines[i], "")
		for j := 0; j < width; j++ {
			board[i-1][j].strategy = items2[j]
		}
	}
	// Calculate RelatedScores after inserting
	for p := 0; p < length; p++ {
		for q := 0; q < width; q++ {
			board[p][q].score = CalculateScore(board, p, q, b, width, length)
		}
	}
	//Animation of gif
	gifImages := Animation(board, b, steps, width, length)
	Process(gifImages, "Prisoners")
}
