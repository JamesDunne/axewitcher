// fsw
package main

type FswState uint8

const (
	FswReset = FswState(1 << iota)
	FswPrev
	FswNext
)
