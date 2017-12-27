package main

import (
	"fmt"
)

func main() {
	// Find footswitch input device by name:
	fswCh, err := ListenFootswitch()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case state := <-fswCh:
			fmt.Printf("state = %v\n", state)
			break
		}
	}
}
