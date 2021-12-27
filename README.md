# Introduction

This golang code is meant to demonstrate an implementation of Brehensam's line algorithm in Golang. 
The article can be [found here](https://johntgz.github.io/2021/12/27/line_rasters_part1_brehensam/)

# Quick Start

To run brehensam's line algorithm on an input, do the following (If no argument flag is provided, the file path will default to "./input.txt")

```
cd src

go run . -f=<INPUT FILE PATH>

#Example1: draw a house
go run . -f=inputs/house.txt
```

line_render/line_render.go is a library kept for helper functions.

# Coordinate Input format

The code works by drawing a line from (x1,y1) to (x2,y2), and must be ended by a semicolon.
```
x1,y1 -> x2, y2;
x1,y1 -> x2, y2;
...

```