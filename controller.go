package axewitcher

import (
	"log"
)

type AmpMode int

const axeMidiChannel = 2

const (
	AmpDirty AmpMode = iota
	AmpClean
	AmpAcoustic
)

type FXState struct {
	Name    string
	MidiCC  uint8
	Enabled bool
}

type AmpState struct {
	Mode      AmpMode
	DirtyGain uint8
	CleanGain uint8
	Volume    uint8
	Fx        [5]FXState
}

type AmpConfig struct {
}

type ControllerState struct {
	sceneIdx int
	Scene    *Scene
	prIdx    int
	pr       *Program
	Amp      [2]AmpState
}

type Scene struct {
	Amp [2]AmpState
}

type Program struct {
	Scenes    []*Scene
	AmpConfig [2]AmpConfig
}

type Controller struct {
	midi Midi

	Programs []*Program
	curr     ControllerState
	prev     ControllerState
}

func NewController(midi Midi) *Controller {
	return &Controller{
		midi: midi,
	}
}

func (c *Controller) HandleFswEvent(ev FswEvent) (err error) {
	// Handle footswitch event:
	if ev.State {
		// Handle footswitch press:
		switch ev.Fsw {
		case FswNext:
			log.Println("next")
			c.curr.sceneIdx++
			if c.curr.sceneIdx >= len(c.curr.pr.Scenes) {
				c.curr.sceneIdx = 0
				c.curr.prIdx++
				if c.curr.prIdx >= len(c.Programs) {
					c.curr.prIdx = 0
				}

				// Update pointers:
				c.curr.pr = c.Programs[c.curr.prIdx]
				c.curr.Scene = c.curr.pr.Scenes[c.curr.sceneIdx]
			}
			break
		case FswPrev:
			break
		case FswReset:
			break
		}
	} else {
		// Handle footswitch release:
	}

	// Send MIDI diff:
	for a := 0; a < 2; a++ {
		// Change amp mode:
		if c.curr.Amp[a].Mode != c.prev.Amp[a].Mode {
			// TODO
			c.midi.CC(axeMidiChannel, 0, 0)
		}
	}

	// Copy to prev state:
	c.prev = c.curr
	return nil
}
