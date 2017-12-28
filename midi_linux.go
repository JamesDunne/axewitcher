// midi_linux
package axewitcher

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
	_, err := m.f.Write([]byte{0xB0 | channel, controller, value})
	return err
}

func (m *midiImpl) PC(channel uint8, program uint8) error {
	_, err := m.f.Write([]byte{0xC0 | channel, program})
	return err
}

func (m *midiImpl) Sysex(data []byte) error {
	sysex := make([]byte, len(data)+2)
	sysex = append(sysex, 0xF0)
	sysex = append(sysex, data...)
	sysex = append(sysex, 0xF7)
	_, err := m.f.Write(sysex)
	return err
}
