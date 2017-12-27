package main

func main() {
	// Listen for footswitch activity:
	fswCh, err := ListenFootswitch()
	if err != nil {
		panic(err)
	}

	// Create MIDI interface:
	midi, err := NewMidi()
	if err != nil {
		panic(err)
	}

	controller := NewController(fswCh, midi)
	controller.Loop()

}
