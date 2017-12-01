// Copyright (c) 2013 Aaron Torres. All rights reserved.
package main

import (
	"fmt"
	"github.com/agtorre/gocolorize"
)

func main() {
	// one way to define a stateful colorizer
	var green gocolorize.Colorize
	green.SetColor(gocolorize.Green)
	g := green.Paint

	// Another way to do it
	red := gocolorize.Colorize{Fg: gocolorize.Red}
	r := red.Paint

	// now with string construction
	green_black := gocolorize.Colorize{Fg: gocolorize.Blue}

	// toggle attributes
	green_black.ToggleUnderline()
	b := green_black.Paint

	//all in 1 line
	c := gocolorize.NewColor("yellow:black").Paint

	fmt.Println(b("On the twelfth day of Christmas"))
	fmt.Println(b("my true love sent to me:"))
	fmt.Println(g("Twelve"), c("Drummers"), r("Drumming"))
	fmt.Println(g("Eleven"), c("Pipers"), r("Piping"))
	fmt.Println(g("Ten"), c("Lords"), r("a Leaping"))
	fmt.Println(g("Nine"), c("Ladies"), r("Dancing"))
	fmt.Println(g("Eight"), c("Maids"), r("a Milking"))
	fmt.Println(g("Seven"), c("Swans"), r("a Swimming"))
	fmt.Println(g("Six"), c("Geese"), r("a Laying"))
	fmt.Println(g("Five"), c("Golden"), r("Rings"))
	fmt.Println(g("Four"), c("Calling"), r("Birds"))
	fmt.Println(g("Three"), c("French"), r("Hens"))
	fmt.Println(g("Two"), c("Turtle"), r("Doves"))
	fmt.Println(b("and a Partridge in a Pear Tree"))

}
