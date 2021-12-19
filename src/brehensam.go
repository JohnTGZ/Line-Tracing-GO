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
		//all relevant data end at ";"
		test_processed := strings.Split(text, ";")

		coordinates := strings.Split(test_processed[0], " -> ")
		xy1 := strings.Split(coordinates[0], ",")
		xy2 := strings.Split(coordinates[1], ",")

		x1, _ := strconv.Atoi(xy1[0])
		y1, _ := strconv.Atoi(xy1[1])
		x2, _ := strconv.Atoi(xy2[0])
		y2, _ := strconv.Atoi(xy2[1])

		//Get max width and height: this code assumes non-negative coordinate values
		max_x := line_render.Max(x1, x2)
		max_y := line_render.Max(y1, y2)
		if max_width < max_x {
			max_width = max_x
		}
		if max_height < max_y {
			max_height = max_y
		}

		xy_pos := []int{x1, y1, x2, y2}

		coordinate_arr = append(coordinate_arr, xy_pos)
	}

	return coordinate_arr, max_width + 1, max_height + 1
}

func brehensamFloat(img *image.RGBA, coords []int, color color.RGBA) {
	x1, y1, x2, y2 := coords[0], coords[1], coords[2], coords[3]
	dx, dy := x2-x1, y2-y1

	gradient_case := 0

	if dx == 0 && dy == 0 {
		fmt.Printf("Start and end coordinates are the same \n")
		return
	} else if dx == 0 { // m == INF
		gradient_case = 5
	} else if dy == 0 { // m == 0
		gradient_case = 6
	} else if dy > 0 { // m > 0
		if dx > 0 { //1st Quadrant, Top right
			if line_render.Abs(dy) > line_render.Abs(dx) { //1st Octant ( 1 < m < INF )
				gradient_case = 1
			} else { //2nd Octant ( 0 < m <= 1 )
				gradient_case = 2
			}
		} else { //4th Quadrant, Top left
			if line_render.Abs(dy) > line_render.Abs(dx) { //8th Octant ( -INF < m < -1 )
				gradient_case = 4
			} else { //7th Octant ( -1 <= m < 0 )
				gradient_case = 3
			}
		}
	} else if dy < 0 { // m < 0
		if dx > 0 { //2nd Quadrant, Btm right
			if line_render.Abs(dy) > line_render.Abs(dx) { //4th Octant ( -INF < m < -1 )
				gradient_case = 4
			} else { //3rd Octant ( -1 <= m < 0 )
				gradient_case = 3
			}
		} else { //3rd Quadrant, Btm Left
			if line_render.Abs(dy) > line_render.Abs(dx) { //5th Octant ( 1 < m < INF )
				gradient_case = 1
			} else { //6th Octant ( 0 < m <= 1 )
				gradient_case = 2
			}
		}
	}

	var m float64
	if gradient_case != 5 { //Only calculate gradient if it is not infinite
		m = (float64(dy)) / (float64(dx)) //gradient
	}

	// if absolute value of gradient > 1, then invert it
	if line_render.Abs(dy) > line_render.Abs(dx) {
		m = 1 / m
	}
	x_inc := line_render.CopySignInt(1, dx)
	y_inc := line_render.CopySignInt(1, dy)

	switch gradient_case {
	case 1: //Octant 1 and 5: 1 < m < INF
		fmt.Printf("1st Octant: 1 < m < INF \n")
		for x, y, err := x1, y1, 0.0; y != y2+y_inc; y += y_inc {
			img.Set(x, y, color)
			if (err + m) < 0.5 {
				err += m
			} else {
				err += m - 1
				x += x_inc
			}
		}
	case 2: //Octant 2 and 6: 0 < m <= 1
		fmt.Printf("2nd Octant: 0 < m <= 1 \n")
		for x, y, err := x1, y1, 0.0; x != x2+x_inc; x += x_inc {
			img.Set(x, y, color)
			if (err + m) < 0.5 {
				err += m
			} else {
				err += m - 1
				y += y_inc
			}
		}
	case 3: //Octant 3 and 7: -1 <= m < 0
		fmt.Printf("3rd Octant: -1 <= m < 0 \n")
		for x, y, err := x1, y1, 0.0; x != x2+x_inc; x += x_inc {
			img.Set(x, y, color)
			if (err + m) > -0.5 {
				err += m
			} else {
				err += m + 1
				y += y_inc
			}
		}
	case 4: //Octant 4 and 8: -INF < m < -1
		fmt.Printf("4th Octant: -INF < m < -1 \n")
		for x, y, err := x1, y1, 0.0; y != y2+y_inc; y += y_inc {
			img.Set(x, y, color)
			if (err + m) > -0.5 {
				err += m
			} else {
				err += m + 1
				x += x_inc
			}
		}
	case 5: // m == INF
		fmt.Printf("Vertical Line: m == INF \n")
		for y := y1; y != y2+y_inc; y += y_inc {
			img.Set(x1, y, color)
		}
	case 6: // m == 0
		fmt.Printf("Horizontal Line: m == 0 \n")
		for x := x1; x != x2+x_inc; x += x_inc {
			img.Set(x, y1, color)
		}

	}

}

func brehensamInt(img *image.RGBA, coords []int, color color.RGBA) {
	x1, y1, x2, y2 := coords[0], coords[1], coords[2], coords[3]
	dx, dy := x2-x1, y2-y1

	grad_sign := 1                    //Sign of gradient
	if bool(dx < 0) != bool(dy < 0) { // (dx is negative) XOR (dy is negative)
		grad_sign = -1
	}
	x_inc := line_render.CopySignInt(1, dx)
	y_inc := line_render.CopySignInt(1, dy)

	if dx == 0 && dy == 0 {
		fmt.Printf("Start and end coordinates are the same \n")
		return
	} else if dx == 0 { // m == INF
		// fmt.Printf("Vertical Line: m == INF \n")
		for y := y1; y != y2+y_inc; y += y_inc {
			img.Set(x1, y, color)
		}
	} else if dy == 0 { // m == 0
		// fmt.Printf("Horizontal Line: m == 0 \n")
		for x := x1; x != x2+x_inc; x += x_inc {
			img.Set(x, y1, color)
		}
	} else if line_render.Abs(dy) > line_render.Abs(dx) { // 1 < m < INF
		// fmt.Printf("1st, 5th, 4th and 8th Octant: 1 < m < INF \n")
		for x, y, err := x1, y1, 0; y != y2+y_inc; y += y_inc {
			// fmt.Printf("(%d, %d), err: %d \n", x, y, err)
			img.Set(x, y, color)
			if grad_sign*2*(err+dx) < dy {
				err += y_inc * (dx)
			} else {
				err += y_inc * (dx - (grad_sign * dy))
				x += x_inc
			}
		}
	} else { // 0 < m <= 1
		// fmt.Printf("2nd, 6th, 3rd and 7th Octant: 0 < m <= 1 \n")
		for x, y, err := x1, y1, 0; x != x2+x_inc; x += x_inc {
			// fmt.Printf("(%d, %d), err: %d \n", x, y, err)
			img.Set(x, y, color)
			if grad_sign*2*(err+dy) < dx {
				err += x_inc * (dy)
			} else {
				err += x_inc * (dy - (grad_sign * dx))
				y += y_inc
			}
		}
	}

}

func main() {

	//reader for taking in input
	// consoleReader := bufio.NewReader(os.Stdin)

	input_file := procArg()

	coordinate_arr, max_width, max_height := getInput(input_file)

	fmt.Printf("Image Width,Height: %d, %d \n", max_width, max_height)

	topLeft := image.Point{0, 0}
	btmRight := image.Point{max_width, max_height}

	img := image.NewRGBA(image.Rectangle{topLeft, btmRight})
	cyan := color.RGBA{100, 200, 200, 0xff}

	//iterate through each coord (x1, y1, x2, y2, dx, dy)
	for _, coord := range coordinate_arr {
		// brehensamFloat(img, coord, cyan)
		brehensamInt(img, coord, cyan)
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
	f, err := os.Create("house.png")
	line_render.CheckErr(err)
	png.Encode(f, img)

}
