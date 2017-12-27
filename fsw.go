// fsw
package main

type FswButton uint8

const (
	FswNone FswButton = iota
	FswReset
	FswPrev
	FswNext
)

type FswEvent struct {
	Fsw   FswButton
	State bool
}
