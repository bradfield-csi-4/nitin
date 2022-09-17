package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go func(val int) {
			fmt.Printf("launched goroutine %d\n", val)
		}(i)
	}
	// Wait for goroutines to finish
	time.Sleep(time.Second)
}

/*
Issue: The variable `i` is initialized by the loop and incremented at the end of each iteration. When you pass `i` to the
goroutine, a data race starts between the loop's writing (i.e. incrementing) and the Printf's reading of `i`.

Fix: We can fix this my initializing a new variable and assigning it to `i`, and then using that new variable in the fmt.Printf
call. Alternatively, we can add a parameter to the goroutine and pass `i` as an argument which effectively does the same
thing by copying the value into a new variable within the function.
*/
