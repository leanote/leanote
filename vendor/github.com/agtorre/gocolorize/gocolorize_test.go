package gocolorize

import (
	"fmt"
	"log"
	"testing"
)

var (
	DEBUG *log.Logger
	WARN  *log.Logger
)

func TestPaint(t *testing.T) {
	var blue Colorize

	//set some state
	blue.SetColor(Blue)

	outString := fmt.Sprintf("Testing %s", blue.Paint("paint"))
	basisString := "Testing \033[0;34mpaint\033[0m"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}

func TestPaintString(t *testing.T) {
	var red Colorize

	//set some state
	red.SetColor(Red)

	outString := red.Paint("Returning a string")
	basisString := "\033[0;31mReturning a string\033[0m"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}

func TestSetColorSetBgColor(t *testing.T) {
	var whiteRedBg Colorize

	//set color and background
	whiteRedBg.SetColor(White)
	whiteRedBg.SetBgColor(Red)

	outString := whiteRedBg.Paint("Setting a foreground and background color!")
	basisString := "\033[0;37m\033[41mSetting a foreground and background color!\033[0m"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}

func TestPaintMultipleInterface(t *testing.T) {
	blue := Colorize{Fg: Blue}
	outString := blue.Paint("Multiple types of args:", 1, 1.24)
	basisString := "\033[0;34mMultiple types of args: 1 1.24\033[0m"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}

func TestPaintComplexType(t *testing.T) {
	green := Colorize{Bg: Green}
	outString := green.Paint("Complex types:",
		struct {
			int
			string
		}{})
	basisString := fmt.Sprintf("\033[42mComplex types: %v\033[0m", struct {
		int
		string
	}{})
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}

func TestInitialize(t *testing.T) {
	blackOnWhite := Colorize{Fg: Black, Bg: White}
	f := blackOnWhite.Paint

	outString := f("Now this is cool")
	basisString := "\033[0;30m\033[47mNow this is cool\033[0m"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}

func TestToggle(t *testing.T) {
	craziness := Colorize{Fg: Yellow, Bg: Black}
	craziness.ToggleFgIntensity()
	craziness.ToggleBgIntensity()
	craziness.ToggleBold()
	craziness.ToggleBlink()
	craziness.ToggleUnderline()
	craziness.ToggleInverse()

	outString := craziness.Paint("craziness")
	basisString := "\033[0;1;5;4;7;93m\033[100mcraziness\033[0m"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}
func TestNewAllToggle(t *testing.T) {
	n := NewColor("yellow+bBuih:black+h")
	outString := n.Paint("all toggles in 1 line!")
	basisString := "\033[0;1;5;4;7;93m\033[100mall toggles in 1 line!\033[0m"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
}

func TestPlain(t *testing.T) {
	plain := Colorize{Fg: Magenta}
	SetPlain(true)
	outString := plain.Paint("plain", "text")
	basisString := "plain text"
	if outString != basisString {
		t.Errorf("Error: string '%s' does not match '%s'\n", outString, basisString)
	} else {
		fmt.Printf("Success: string: '%s' matches '%s'\n", outString, basisString)
	}
	SetPlain(false)
}
