package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gvalkov/golang-evdev"
)

func ListenFootswitch() (fswCh chan FswEvent, err error) {
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

	go func() {
		defer close(fswCh)

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
			// NOTE: unfortunately the footswitch driver does not allow multiple switches to be depressed simultaneously.
			button := FswNone
			if key.Scancode == evdev.KEY_A {
				button = FswReset
			} else if key.Scancode == evdev.KEY_B {
				button = FswPrev
			} else if key.Scancode == evdev.KEY_C {
				button = FswNext
			}

			fswCh <- FswEvent{
				Fsw:   button,
				State: key.State == evdev.KeyDown,
			}
		}
	}()

	return
}
