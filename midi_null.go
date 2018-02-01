package axewitcher

type nullMidiImpl struct {
}

func NewNullMidi() (midi Midi, err error) {
	midi = &nullMidiImpl{}
	return
}

func (m *nullMidiImpl) Close() error {
	return nil
}
func (m *nullMidiImpl) CC(channel uint8, controller uint8, value uint8) error {
	return nil
}
func (m *nullMidiImpl) PC(channel uint8, program uint8) error {
	return nil
}
func (m *nullMidiImpl) Sysex(data []byte) error {
	return nil
}
