// midi_darwin
package main

type midiImpl struct {
}

func NewMidi() (midi Midi, err error) {
	midi = &midiImpl{}
	return
}

func (m *midiImpl) Close() error {
	return nil
}
func (m *midiImpl) CC(channel uint8, controller uint8, value uint8) error {
	return nil
}
func (m *midiImpl) PC(channel uint8, program uint8) error {
	return nil
}
func (m *midiImpl) Sysex(data []byte) error {
	return nil
}
