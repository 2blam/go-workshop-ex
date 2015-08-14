package main

import "fmt"
import "math"

// Shape defines the methods that a shape should have
type Shape interface {
	// Name returns the name of the Shape
	Name() string

	// Area returns the area of the Shape
	Area() float64
}

func main() {
	shapes := []Shape{
		Rectangle{Width: 10, Height: 2},
		Triangle{Base: 3, Height: 4},
		Circle{Radius: 3},
	}

	for _, shape := range shapes {
		fmt.Printf("%s.Area() = %v\n", shape.Name(), shape.Area())
	}
}

type Rectangle struct {
	Width, Height float64
}

// implement Rectangle methods here

func (rect Rectangle) Name() string {
	return "Rectangle"
}

func (rect Rectangle) Area() float64 {
	return (rect.Width + rect.Height) * 2
}

type Triangle struct {
	Base, Height float64
}

// implement Triangle methods here
func (tri Triangle) Name() string {
	return "Triangle"
}

func (tri Triangle) Area() float64 {
	return (tri.Base + tri.Height) / 2
}

type Circle struct {
	Radius float64
}

// implement Circle methods here
func (cir Circle) Name() string {
	return "Circle"
}

func (cir Circle) Area() float64 {
	return math.Pow(cir.Radius, 2) * 3.14
}
