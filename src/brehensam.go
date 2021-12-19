package main

import (
	"bufio" //For reading line by line
	"flag"  //For command line parsing
	"fmt"
	"image"
	"image/color"
	"image/png"
	"line_render/line_render"
	"os" //for opening filess
	"strconv"
	"strings"
)

/* Process command line arguemnts to read from either
test input or normal input file
*/
func procArg() string {
	testingPtr := flag.Bool("t", false, "Enable testing")

	flag.Parse()

	input_file := "input.txt"
	if *testingPtr {
		input_file = "test_input.txt"
	}
	return input_file
}

/* Read the input into an iterable array
 */
func getInput(filepath string) ([][]int, int, int) {

	f, err := os.Open(filepath)
	line_render.CheckErr(err)

	//close file at end of program
	defer f.Close()

	//create scanner object to read line by line
	scanner := bufio.NewScanner(f)

	var coordinate_arr [][]int

	max_width := 0
	max_height := 0

	//read file line by line
	for scanner.Scan() {
		text := scanner.Text()

		coordinates := strings.Split(text, " -> ")
		xy1 := strings.Split(coordinates[0], ",")
		xy2 := strings.Split(coordinates[1], ",")

		x1, _ := strconv.Atoi(xy1[0])
		y1, _ := strconv.Atoi(xy1[1])
		x2, _ := strconv.Atoi(xy2[0])
		y2, _ := strconv.Atoi(xy2[1])

		//Get max width and height
		if x1 > max_width {
			max_width = x1
		} else if x2 > max_width {
			max_width = x2
		}
		if y1 > max_height {
			max_height = y1
		} else if y2 > max_height {
			max_height = y2
		}

		xy_pos := []int{x1, y1, x2, y2}

		coordinate_arr = append(coordinate_arr, xy_pos)
	}

	return coordinate_arr, max_width + 1, max_height + 1
}

func brehensamFloat64(img *image.RGBA, coords []int, color color.RGBA) {
	x1, y1, x2, y2 := coords[0], coords[1], coords[2], coords[3]
	dx, dy := x2-x1, y2-y1
	fmt.Printf("dx: %d, dy: %d \n", dx, dy)

	gradient_case := 0

	if dx == 0 && dy == 0 {
		fmt.Printf("Start and end coordinates are the same \n")
		return
	} else if dx == 0 { //infinite gradient
		gradient_case = 1
	} else if dy == 0 { //zero gradient
		gradient_case = 2
	} else if dy > 0 {
		if dx > 0 { //+ve gradient
			gradient_case = 3
		} else { //-ve gradient
			gradient_case = 4
		}
	} else if dy < 0 {
		if dx < 0 { //+ve gradient
			gradient_case = 3
		} else { //-ve gradient
			gradient_case = 4
		}
	}

	var m float64
	if gradient_case != 1 { //Only calculate gradient if it is not infinite
		m = (float64(dy)) / (float64(dx)) //gradient
	}

	switch gradient_case {
	case 1: //Infinite Gradient
		fmt.Printf("Infinite gradient \n")
		if dy > 0 { //Case 1a: Increasing y
			for y := y1; y <= y2; y++ {
				img.Set(x1, y, color)
			}
		} else { //Case 1b: Decreasing y
			for y := y1; y >= y2; y-- {
				img.Set(x1, y, color)
			}
		}
	case 2: //Zero Gradient
		fmt.Printf("Zero gradient \n")
		if dx > 0 { //Case 2a: Increasing x
			for x := x1; x <= x2; x++ {
				img.Set(x, y1, color)
			}
		} else { //Case 2b: Decreasing x

			for x := x1; x >= x2; x-- {
				img.Set(x, y1, color)
			}
		}
	case 3: //Positive Gradient
		fmt.Printf("Positive gradient \n")
		if dx > 0 { //Case 3a: (dx > 0 && dy > 0)
			for x, y, err := x1, y1, 0.0; x <= x2; x++ {
				img.Set(x, y, color)
				if (err + m) < 0.5 {
					err += m
				} else {
					err += m + 1
					y++
				}
			}
		} else { //Case 3b: (dx < 0 && dy < 0)
			for x, y, err := x1, y1, 0.0; x >= x2; x-- {
				img.Set(x, y, color)
				if (err + m) < 0.5 {
					err += m
				} else {
					err += m + 1
					y--
				}
			}
		}
	case 4: //Negative Gradient
		fmt.Printf("Negative gradient \n")
		if dx > 0 { //Case 4a: (dx > 0 && dy < 0)
			for x, y, err := x1, y1, 0.0; x <= x2; x++ {
				img.Set(x, y, color)
				if (err + m) > -0.5 {
					err += m
				} else {
					err += m + 1
					y--
				}
			}
		} else { //Case 4b: (dx < 0 && dy > 0)
			for x, y, err := x1, y1, 0.0; x >= x2; x-- {
				img.Set(x, y, color)
				if (err + m) > -0.5 {
					err += m
				} else {
					err += m + 1
					y++
				}
			}
		}
	}

}

func brehensamInt(img *image.RGBA, coords []int, color color.RGBA) {
	x1, y1, x2, y2 := coords[0], coords[1], coords[2], coords[3]
	dx, dy := x2-x1, y2-y1

	//error (fluctutates between -0.5 and 0.5)
	err := 0

	if dx > 0 {
		x, y := x1, y1
		for x <= x2 {
			img.Set(x, y, color)
			if 2*(err+dy) < dx {
				err += dy
			} else {
				err += dy - dx
				y++
			}
			x++
		}
	}

}

func main() {

	input_file := procArg()

	coordinate_arr, max_width, max_height := getInput(input_file)

	topLeft := image.Point{0, 0}
	btmRight := image.Point{max_width, max_height}

	img := image.NewRGBA(image.Rectangle{topLeft, btmRight})
	cyan := color.RGBA{100, 200, 200, 0xff}

	//iterate through each coord (x1, y1, x2, y2, dx, dy)
	for _, coord := range coordinate_arr {
		brehensamFloat64(img, coord, cyan)
	}

	// Encode as PNG.
	f, err := os.Create("image.png")
	line_render.CheckErr(err)
	png.Encode(f, img)

}
