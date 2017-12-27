// midi
package main

type Midi interface {
	Close() error

	CC(channel uint8, controller uint8, value uint8) error
	PC(channel uint8, program uint8) error
	Sysex(data []byte) error
}
