/*
Write two goroutines which have a race condition when executed concurrently.
Explain what the race condition is and how it can occur.
*/

// run: go run -race init.go to check that there is a race condition

package main

var shared int = 0
var result int

func read() {
	result = shared
}

func write() {
	shared += 1
}

func main() {
	// there is 2 read/write goroutines that use common variable
	// as the order is undefined, race condition occurs
	go read()
	go write()
}
