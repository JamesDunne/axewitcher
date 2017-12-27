// midi_linux
package main

import (
	"os"
)

type midiImpl struct {
	f *os.File
}

func NewMidi() (midi Midi, err error) {
	// Open midi device:
	var f *os.File
	f, err = os.OpenFile("/dev/midi1", os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	midi = &midiImpl{
		f: f,
	}
	return
}

func (m *midiImpl) Close() error {
	return m.f.Close()
}

func (m *midiImpl) CC(channel uint8, controller uint8, value uint8) error {
	n, err := m.f.Write([]byte{0xB0 | channel, controller, value})
	return err
}

func (m *midiImpl) PC(channel uint8, program uint8) error {
	n, err := m.f.Write([]byte{0xC0 | channel, program})
	return err
}

func (m *midiImpl) Sysex(data []byte) error {
	n, err := m.f.Write([]byte{0xF0, ...data, 0xF7})
	return err
}
