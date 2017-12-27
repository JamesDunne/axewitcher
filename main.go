package main

func main() {
	// Listen for footswitch activity:
	fswCh, err := ListenFootswitch()
	if err != nil {
		panic(err)
	}

	controller := NewController(fswCh)
	controller.Loop()

}
