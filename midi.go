// midi
package axewitcher

type Midi interface {
	Close() error

	// Controller change:
	CC(channel uint8, controller uint8, value uint8) error
	// Program change:
	PC(channel uint8, program uint8) error
	// send sysex data; F0 and F7 are automatically included
	Sysex(data []byte) error
}
