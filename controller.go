package main

type AmpMode int

const axeMidiChannel = 2

const (
	AmpDirty AmpMode = iota
	AmpClean
	AmpAcoustic
)

type FXState struct {
	name    string
	midiCC  uint8
	enabled bool
}

type AmpState struct {
	mode      AmpMode
	dirtyGain uint8
	cleanGain uint8
	volume    uint8
	fx        [5]FXState
}

type ControllerState struct {
	sceneIdx int
	scene    *Scene
	prIdx    int
	pr       *Program
	amp      [2]AmpState
}

type Scene struct {
	amp [2]AmpState
}

type Program struct {
	scenes []Scene
}

type Controller struct {
	midi Midi

	programs []Program
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
			c.curr.sceneIdx++
			if c.curr.sceneIdx >= len(c.curr.pr.scenes) {
				c.curr.sceneIdx = 0
				c.curr.prIdx++
				if c.curr.prIdx >= len(c.programs) {
					c.curr.prIdx = 0
				}

				// Update pointers:
				c.curr.pr = &c.programs[c.curr.prIdx]
				c.curr.scene = &c.curr.pr.scenes[c.curr.sceneIdx]
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
		if c.curr.amp[a].mode != c.prev.amp[a].mode {
			// TODO
			c.midi.CC(axeMidiChannel, 0, 0)
		}
	}

	// Copy to prev state:
	c.prev = c.curr
	return nil
}
