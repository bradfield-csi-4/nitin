package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("performing initialization...")
		done <- struct{}{}
	}()

	<-done
	fmt.Println("initialization done, continuing with rest of program")
}

/*
Issue: The send of the 1st message on a buffered channel won't synchronize with a receive, therefore
instead of blocking until initialization completes, there is a race condition between the initialization
goroutine and the main goroutine. WHen the main goroutine wins, the original program will simply print
"initialization done" and exit (even though initialization wasn't performed).

Fix: We can swap the send and receive operations. Now the receive blocks until it detects a message
on the buffered channel, and by the time a message is sent on the channel, initialization has completed.
*/
