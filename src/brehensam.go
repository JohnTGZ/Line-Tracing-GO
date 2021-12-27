package main

import (
	"bufio" //For reading line by line
	"flag"

	//For command line parsing
	"fmt"
	"image"
	"image/color"
	"image/png"
	"line_render/line_render"
	"log"
	"os" //for opening filess
	"strconv"
	"strings"
	"time"
)

/* Process command line arguemnts to read from either
test input or normal input file
*/
func processArgFlags() string {
	file_name_ptr := flag.String("f", "input.txt", "Input file name")

	flag.Parse()

	return *file_name_ptr
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

	grad_sign := 1.0                  //Sign of gradient
	if bool(dx < 0) != bool(dy < 0) { // (dx is negative) XOR (dy is negative)
		grad_sign = -1.0
	}

	//We use x_inc and y_inc to either increment or decrement x or y depending on the sign of the gradient
	x_inc := line_render.CopySignInt(1, dx)
	y_inc := line_render.CopySignInt(1, dy)

	var m float64
	if dx != 0 { //Only calculate gradient if it is not infinite
		m = (float64(dy)) / (float64(dx)) //gradient
	}
	if line_render.Abs(dy) > line_render.Abs(dx) { // if absolute value of gradient > 1, then invert it
		m = 1 / m
	}

	if dx == 0 && dy == 0 {
		fmt.Printf("Start and end coordinates are the same \n")
		return
	} else if dx == 0 { // m == INF
		for y := y1; y != y2+y_inc; y += y_inc {
			img.Set(x1, y, color)
		}
	} else if dy == 0 { // m == 0
		for x := x1; x != x2+x_inc; x += x_inc {
			img.Set(x, y1, color)
		}
	} else if line_render.Abs(dy) > line_render.Abs(dx) { // 1 < abs(m) < INF
		// fmt.Printf("1st, 5th, 4th and 8th Octant: 1 < abs(m) < INF \n")
		for x, y, err := x1, y1, 0.0; y != y2+y_inc; y += y_inc {
			img.Set(x, y, color)
			if grad_sign*(err+m) < 0.5 {
				err += m
			} else {
				err += m - grad_sign*1
				x += x_inc
			}
		}
	} else { // 0 < abs(m) <= 1
		// fmt.Printf("2nd, 6th, 3rd and 7th Octant: 0 < abs(m) <= 1 \n")
		for x, y, err := x1, y1, 0.0; x != x2+x_inc; x += x_inc {
			img.Set(x, y, color)
			if grad_sign*(err+m) < 0.5 {
				err += m
			} else {
				err += m - grad_sign*1
				y += y_inc
			}
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
		for y := y1; y != y2+y_inc; y += y_inc {
			img.Set(x1, y, color)
		}
	} else if dy == 0 { // m == 0
		for x := x1; x != x2+x_inc; x += x_inc {
			img.Set(x, y1, color)
		}
	} else if line_render.Abs(dy) > line_render.Abs(dx) { // 1 < abs(m) < INF
		// fmt.Printf("1st, 5th, 4th and 8th Octant: 1 < abs(m) < INF \n")
		for x, y, err := x1, y1, 0; y != y2+y_inc; y += y_inc {
			img.Set(x, y, color)
			if grad_sign*2*(err+dx) < dy {
				err += y_inc * (dx)
			} else {
				err += y_inc * (dx - (grad_sign * dy))
				x += x_inc
			}
		}
	} else { // 0 < m <= 1
		// fmt.Printf("2nd, 6th, 3rd and 7th Octant: 0 < abs(m) <= 1 \n")
		for x, y, err := x1, y1, 0; x != x2+x_inc; x += x_inc {
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
	input_file := processArgFlags()

	coordinate_arr, max_width, max_height := getInput(input_file)
	// fmt.Printf("Image Width,Height: %d, %d \n", max_width, max_height)

	img_topLeft := image.Point{0, 0}
	img_btmRight := image.Point{max_width, max_height}

	img_int := image.NewRGBA(image.Rectangle{img_topLeft, img_btmRight})
	img_float := image.NewRGBA(image.Rectangle{img_topLeft, img_btmRight})
	cyan := color.RGBA{100, 200, 200, 0xff}

	start := time.Now()
	for _, coord := range coordinate_arr { //iterate through each coord (x1, y1, x2, y2)
		brehensamFloat(img_float, coord, cyan)
	}
	elapsed_float := time.Since(start)
	// Encode as PNG and save it to a file
	f_float, err := os.Create("house_float.png")
	line_render.CheckErr(err)
	png.Encode(f_float, img_float)

	start = time.Now()
	for _, coord := range coordinate_arr {
		brehensamInt(img_int, coord, cyan)
	}
	elapsed_int := time.Since(start)

	// Encode as PNG and save it to a file
	f_int, err := os.Create("house_int.png")
	line_render.CheckErr(err)
	png.Encode(f_int, img_int)

	log.Printf("Brehensam Float took %s", elapsed_float)
	log.Printf("Brehensam Integer took %s", elapsed_int)
}
