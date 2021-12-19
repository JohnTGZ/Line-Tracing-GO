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

		//Get max width and height: this code assumes non-negative coordinate values
		max_width = line_render.Max(x1, x2)
		max_height = line_render.Max(y1, y2)

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
	} else if dx == 0 { // m == INF
		gradient_case = 9
	} else if dy == 0 { // m == 0
		gradient_case = 10
	} else if dy > 0 { // m > 0
		if dx > 0 { //1st Quadrant, Top right
			if line_render.Abs(dy) > line_render.Abs(dx) { //1st Octant ( 1 < m < INF )
				gradient_case = 1
			} else { //2nd Octant ( 0 < m <= 1 )
				gradient_case = 2
			}
		} else { //3rd Quadrant, Btm Left
			if line_render.Abs(dy) > line_render.Abs(dx) { //5th Octant ( 1 < m < INF )
				gradient_case = 5
			} else { //6th Octant ( 0 < m <= 1 )
				gradient_case = 6
			}
		}
	} else if dy < 0 { // Negative gradient
		if dx > 0 { //2nd Quadrant, Btm right
			if line_render.Abs(dy) > line_render.Abs(dx) { //4th Octant ( -INF < m < -1 )
				gradient_case = 4
			} else { //3rd Octant ( -1 <= m < 0 )
				gradient_case = 3
			}
		} else { //4th Quadrant, Top Left
			if line_render.Abs(dy) > line_render.Abs(dx) { //8th Octant ( -INF < m < -1 )
				gradient_case = 8
			} else { //7th Octant ( -1 <= m < 0 )
				gradient_case = 7
			}
		}
	}

	var m float64
	if gradient_case != 9 { //Only calculate gradient if it is not infinite
		m = (float64(dy)) / (float64(dx)) //gradient
	}

	// if absolute value of gradient > 1, then invert it
	if line_render.Abs(dy) > line_render.Abs(dx) {
		m = 1 / m
	}
	x_inc := line_render.CopySignInt(1, dx)
	y_inc := line_render.CopySignInt(1, dy)

	switch gradient_case {
	case 1: //Quad 1: 1 < m < INF
		fmt.Printf("1st Octant: 1 < m < INF \n")
		for x, y, err := x1, y1, 0.0; y != y2+1; y += y_inc {
			img.Set(x, y, color)
			if (err + m) < 0.5 {
				err += m
			} else {
				err += m - 1
				x += x_inc
			}
		}
	case 2: //Quad 2: 0 < m <= 1
		fmt.Printf("2nd Octant: 0 < m <= 1 \n")
		for x, y, err := x1, y1, 0.0; x != x2+1; x += x_inc {
			fmt.Printf("2nd: %d, %d \n", x, y)
			img.Set(x, y, color)
			if (err + m) < 0.5 {
				err += m
			} else {
				err += m - 1
				y += y_inc
			}
		}
	case 3: //Quad 3: -1 <= m < 0
		fmt.Printf("3rd Octant: -1 <= m < 0 \n")
		for x, y, err := x1, y1, 0.0; x != x2+1; x += x_inc {
			img.Set(x, y, color)
			if (err + m) > -0.5 {
				err += m
			} else {
				err += m + 1
				y += y_inc
			}
		}
	case 4: //Quad 4: -INF < m < -1
		fmt.Printf("4th Octant: -INF < m < -1 \n")
		for x, y, err := x1, y1, 0.0; y != y2-1; y += y_inc {
			img.Set(x, y, color)
			if (err + m) > -0.5 {
				err += m
			} else {
				err += m + 1
				x += x_inc
			}
		}
	case 5: //Quad 5:  1 < m < INF
		fmt.Printf("5th Octant: 1 < m < INF \n")
		for x, y, err := x1, y1, 0.0; y != y2-1; y += y_inc {
			img.Set(x, y, color)
			if (err + m) < 0.5 {
				err += m
			} else {
				err += m - 1
				x += x_inc
			}
		}
	case 6: //Quad 6:  0 < m <= 1
		fmt.Printf("6th Octant: 0 < m <= 1 \n")
		for x, y, err := x1, y1, 0.0; x != x2-1; x += x_inc {
			img.Set(x, y, color)
			if (err + m) < 0.5 {
				err += m
			} else {
				err += m - 1
				y += y_inc
			}
		}
	case 7: //Quad 7:  -1 <= m < 0
		fmt.Printf("7th Octant: -1 <= m < 0 \n")
		for x, y, err := x1, y1, 0.0; x != x2-1; x += x_inc {
			img.Set(x, y, color)
			if (err + m) > -0.5 {
				err += m
			} else {
				err += m + 1
				y += y_inc
			}
		}
	case 8: //Quad 8:  -INF < m < -1
		fmt.Printf("8th Octant: -INF < m < -1 \n")
		for x, y, err := x1, y1, 0.0; y != y2+1; y += y_inc {
			img.Set(x, y, color)
			if (err + m) > -0.5 {
				err += m
			} else {
				err += m + 1
				x += x_inc
			}
		}
	case 9: // m == INF
		fmt.Printf("m == INF \n")
		if dy > 0 { //Case 9a: Increasing y
			for y := y1; y <= y2; y++ {
				img.Set(x1, y, color)
			}
		} else { //Case 9b: Decreasing y
			for y := y1; y >= y2; y-- {
				img.Set(x1, y, color)
			}
		}
	case 10: // m == 0
		fmt.Printf(" m == 0 \n")
		if dx > 0 { //Case 10a: Increasing x
			for x := x1; x <= x2; x++ {
				img.Set(x, y1, color)
			}
		} else { //Case 10b: Decreasing x

			for x := x1; x >= x2; x-- {
				img.Set(x, y1, color)
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

	//reader for taking in input
	// consoleReader := bufio.NewReader(os.Stdin)

	input_file := procArg()

	coordinate_arr, max_width, max_height := getInput(input_file)

	fmt.Printf("Width,Height: %d, %d \n", max_width, max_height)

	topLeft := image.Point{0, 0}
	btmRight := image.Point{9, 9}

	img := image.NewRGBA(image.Rectangle{topLeft, btmRight})
	cyan := color.RGBA{100, 200, 200, 0xff}

	//iterate through each coord (x1, y1, x2, y2, dx, dy)
	for _, coord := range coordinate_arr {
		brehensamFloat64(img, coord, cyan)

		// char, _, err := consoleReader.ReadRune()
		// line_render.CheckErr(err)
		// fmt.Printf("Current input: %s", char)

		// switch char {
		// case 'n':
		// 	fmt.Println("Stepping through program")
		// 	continue
		// case 'q':
		// 	fmt.Println("Exiting line-tracing program")
		// 	break
		// }
	}

	// Encode as PNG.
	f, err := os.Create("image.png")
	line_render.CheckErr(err)
	png.Encode(f, img)

}
