# timeago - A time formatting package

## Install

	go get github.com/xeonx/timeago

## Docs

<http://godoc.org/github.com/xeonx/timeago>

## Use

	package main

	import (
		"time"
		"github.com/xeonx/timeago"
	)
		
	func main() {
		t := time.Now().Add(42 * time.Second)
		
		s := timeago.English.Format(t)
		//s will contains "less than a minute ago"
		
		//...
	}
	
## Tests

`go test` is used for testing.