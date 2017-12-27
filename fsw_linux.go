package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gvalkov/golang-evdev"
)

func ListenFootswitch() (fswCh chan FswState, err error) {
	fsw := (*evdev.InputDevice)(nil)

	// List all input devices:
	devs, err := evdev.ListInputDevices()
	if err != nil {
		return
	}
	for _, dev := range devs {
		// Find foot switch device:
		if strings.Contains(dev.Name, "PCsensor FootSwitch3") {
			fsw = dev
			break
		}
	}
	if fsw == nil {
		err = errors.New("No footswitch device found!")
		return
	}
	fmt.Printf("%v\n", fsw)

	fswCh = make(chan FswState)
	go func() {
		defer close(fswCh)

		fswState := FswState(0)
		for {
			ev, err := fsw.ReadOne()
			if err != nil {
				break
			}
			if ev.Type != evdev.EV_KEY {
				continue
			}

			key := evdev.NewKeyEvent(ev)
			if key.State == evdev.KeyHold {
				continue
			}

			// Determine which footswitch was pressed/released:
			mask := FswState(0)
			if key.Scancode == evdev.KEY_A {
				mask = FswReset
			} else if key.Scancode == evdev.KEY_B {
				mask = FswPrev
			} else if key.Scancode == evdev.KEY_C {
				mask = FswNext
			}

			// Apply mask based on down/up state:
			if key.State == evdev.KeyDown {
				fswState |= mask
			} else if key.State == evdev.KeyUp {
				fswState &= ^mask
			}

			fswCh <- fswState
		}
	}()

	return
}
