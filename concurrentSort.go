/*
Write a program to sort an array of integers.
The program should partition the array into 4 parts, each of which is sorted by a different goroutine.
Each partition should be of approximately equal size.
Then the main goroutine should merge the 4 sorted subarrays into one large sorted array.
The program should prompt the user to input a series of integers.
Each goroutine which sorts ¼ of the array should print the subarray that it will sort.
When sorting is complete, the main goroutine should print the entire sorted list.
*/

// waitgroups and channels are used

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var slice *[]int

func ParseInts(s string) (*[]int, error) {
	var (
		fields = strings.Fields(s)
		ints   = make([]int, len(fields))
		err    error
	)
	for i, f := range fields {
		ints[i], err = strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
	}
	return &ints, nil
}

func Sort(slice *[]int, wg *sync.WaitGroup) {
	sort.Ints(*slice)
	wg.Done()
}

func Merge(left, right *[]int, c chan *[]int) {
	size, i, j := len(*left)+len(*right), 0, 0
	result := make([]int, size, size)
	for k := 0; k < size; k++ {
		if i > len(*left)-1 && j <= len(*right)-1 {
			result[k] = (*right)[j]
			j++
		} else if j > len(*right)-1 && i <= len(*left)-1 {
			result[k] = (*left)[i]
			i++
		} else if (*left)[i] < (*right)[j] {
			result[k] = (*left)[i]
			i++
		} else {
			result[k] = (*right)[j]
			j++
		}
	}
	c <- &result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter an array:")
	scanner.Scan()
	inputLine := scanner.Text()
	slice, _ = ParseInts(inputLine)
	if len(*slice) < 4 {
		sort.Ints(*slice)
		fmt.Println(*slice)
	} else {
		step := len(*slice) / 4
		slice0 := (*slice)[:step]
		slice1 := (*slice)[step : step*2]
		slice2 := (*slice)[step*2 : step*3]
		slice3 := (*slice)[step*3:]
		wg := new(sync.WaitGroup)
		wg.Add(4)
		go Sort(&slice0, wg)
		go Sort(&slice1, wg)
		go Sort(&slice2, wg)
		go Sort(&slice3, wg)
		wg.Wait()
		channel := make(chan *[]int, 2)
		go Merge(&slice0, &slice1, channel)
		go Merge(&slice2, &slice3, channel)
		lMerge := <-channel
		rMerge := <-channel
		Merge(lMerge, rMerge, channel)
		result := <-channel
		fmt.Println(*result)
	}
}
