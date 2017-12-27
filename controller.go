package main

import (
	"fmt"
)

type AmpMode int

const (
	AmpDirty AmpMode = iota
	AmpClean
	AmpAcoustic
)

type AmpState struct {
	mode      AmpMode
	dirtyGain uint8
	cleanGain uint8
}

type ControllerState struct {
	fsw FswState
	amp [2]AmpState
}

type Controller struct {
	fswCh <-chan FswState
	curr  ControllerState
	prev  ControllerState
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
			c.curr.fsw = state
			fmt.Printf("state = %v\n", state)
			break
		}
	}
	return nil
}
