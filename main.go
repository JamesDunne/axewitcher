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
	defer midi.Close()

	// Initialize controller:
	controller := NewController(midi)

	// Run an idle loop awaiting events:
	for {
		select {
		case ev := <-fswCh:
			controller.HandleFswEvent(ev)
			break
		}
	}
}
