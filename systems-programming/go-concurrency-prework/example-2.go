package main

import (
	"fmt"
)

const numTasks = 3

func main() {
	done := make(chan struct{})
	for i := 0; i < numTasks; i++ {
		go func() {
			fmt.Println("running task...")

			// Signal that task is done
			done <- struct{}{}
		}()
	}
	// Wait for tasks to complete
	for i := 0; i < numTasks; i++ {
		<-done
	}
	fmt.Printf("all %d tasks done!\n", numTasks)
}

/*
Issue: The `done` channel was declared but not initialized and is therefore set to its zero value of `nil`. You cannot send or receive
on a `nil` channel.

Fix: To properly create a channel, we can use the `make` keyword.
*/
