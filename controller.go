package main

import (
	"fmt"
)

type Controller struct {
	fswCh <-chan FswState
}

func NewController(fswCh <-chan FswState) *Controller {
	return &Controller{
		fswCh: fswCh,
	}
}

func (c *Controller) Loop() (err error) {
	for {
		select {
		case state := <-c.fswCh:
			fmt.Printf("state = %v\n", state)
			break
		}
	}
	return nil
}
