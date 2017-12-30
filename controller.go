package axewitcher

//"log"

const axeMidiChannel = 2

var fxNames = []string{
	"Cho1", // 41
	"Cho2", // 42
	"Cmp1", // 43
	"Cmp2", // 44
	"Crs1", // 45
	"Crs2", // 46
	"Dly1", // 47
	"Dly2", // 48
	"Drv1", // 49
	"Drv2", // 50
	"Enh1", // 51
	"Flt1", // 52
	"Flt2", // 53
	"Flt3", // 54
	"Flt4", // 55
	"Flg1", // 56
	"Flg2", // 57
	"Fmnt", // 58
	"Fxlp", // 59
	"Gte1", // 60
	"Gte2", // 61
	"Geq1", // 62
	"Geq2", // 63
	"Geq3", // 64
	"Geq4", // 65
	"Mega", // 66
	"Mcm1", // 67
	"Mcm2", // 68
	"Mdy1", // 69
	"Mdy2", // 70
	"Peq2", // 72
	"Peq3", // 73
	"Peq4", // 74
	"Peq1", // 71
	"Phr1", // 75
	"Phr2", // 76
	"Pit1", // 77
	"Pit2", // 78
	"Qch1", // 79
	"Qch2", // 80
	"Rsn1", // 81
	"Rsn2", // 82
	"Rvb1", // 83
	"Rvb2", // 84
	"Ring", // 85
	"Rtr1", // 86
	"Rtr2", // 87
	"Syn1", // 88
	"Syn2", // 89
	"Trm1", // 90
	"Trm2", // 91
	"Voco", // 92
	"Vol1", // 93
	"Vol2", // 94
	"Vol3", // 95
	"Vol4", // 96
	"Wah1", // 97
	"Wah2", // 98
}

type AmpMode int

const (
	AmpDirty AmpMode = iota
	AmpClean
	AmpAcoustic
)

type FXState struct {
	Enabled bool
}
type FXConfig struct {
	Name   string
	MidiCC uint8
}

type AmpState struct {
	Mode      AmpMode
	DirtyGain uint8
	CleanGain uint8
	Volume    uint8
	Fx        [5]FXState
}

type AmpConfig struct {
	Fx [5]FXConfig
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

	Curr ControllerState
	Prev ControllerState
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
			c.Curr.sceneIdx++
			if c.Curr.sceneIdx >= len(c.Curr.pr.Scenes) {
				c.Curr.sceneIdx = 0
				c.Curr.prIdx++
				if c.Curr.prIdx >= len(c.Programs) {
					c.Curr.prIdx = 0
				}

				// Update pointers:
				c.Curr.pr = c.Programs[c.Curr.prIdx]
				c.Curr.Scene = c.Curr.pr.Scenes[c.Curr.sceneIdx]
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
		if c.Curr.Amp[a].Mode != c.Prev.Amp[a].Mode {
			// TODO
			c.midi.CC(axeMidiChannel, 0, 0)
		}
	}

	// Copy to prev state:
	c.Prev = c.Curr
	return nil
}
