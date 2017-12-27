// fsw
package main

type FswState uint8

const (
	FswPrev = FswState(1 << iota)
	FswReset
	FswNext
)
