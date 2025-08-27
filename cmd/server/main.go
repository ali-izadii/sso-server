package main

import "github.com/common-nighthawk/go-figure"

func main() {
	myFigure := figure.NewColorFigure("SSO Server", "", "green", true)
	myFigure.Print()
}
