/*
There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.
Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).
In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
The host allows no more than 2 philosophers to eat concurrently.
Each philosopher is numbered, 1 through 5.
When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>” on a line by itself, where <number> is the number of the philosopher.
When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>” on a line by itself, where <number> is the number of the philosopher.
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	wg   sync.WaitGroup
	host = make(chan bool, 2)
)

// Chopstick
type ChopS struct{ sync.Mutex }

// Philosopher
type Philo struct {
	number  int
	leftCS  *ChopS
	rightCS *ChopS
}

// Returns a random boolean value based on the current time
func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func (p Philo) eat() {

	// eat 3 times
	for i := 0; i < 3; i++ {

		// goroutine/eating is blocked if channel is full
		host <- true

		// pick chopsticks up randomly
		switch RandBool() {
		case true:
			p.leftCS.Lock()
			p.rightCS.Lock()
		default:
			p.rightCS.Lock()
			p.leftCS.Lock()
		}

		fmt.Println("starting to eat ", p.number)
		fmt.Println("finishing eating ", p.number)

		p.leftCS.Unlock()
		p.rightCS.Unlock()

		<-host

	}
	wg.Done()
}

func main() {

	ChopSticks := make([]*ChopS, 5)
	for i := 0; i < 5; i++ {
		ChopSticks[i] = new(ChopS)
	}

	Philosophers := make([]*Philo, 5)
	for i := 0; i < 5; i++ {
		Philosophers[i] = &Philo{i + 1, ChopSticks[i], ChopSticks[(i+1)%5]}
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go Philosophers[i].eat()
	}

	// main gouroutine always waits other ones
	wg.Wait()

}
