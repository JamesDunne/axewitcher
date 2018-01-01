package axewitcher

import (
	"log"
	"math"
	"strings"
)

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

func findFxMidiCC(name string) uint8 {
	nameLower := strings.ToLower(name)
	for i, cmp := range fxNames {
		if strings.ToLower(cmp) == nameLower {
			return uint8(41 + i)
		}
	}
	return 255
}

//  p = 10 ^ (dB / 20)
// dB = log10(p) * 20
// Log20A means 20% percent at half-way point of knob, i.e. dB = 20 * ln(0.20) / ln(10) = -13.98dB
func dB(percent float64) float64 {
	db := math.Log10(percent) * 20.0
	return db
}

func MIDItoDB(n uint8) float64 {
	p := float64(n) / 127.0
	// log20a taper (50% -> 20%)
	p = (math.Pow(15.5, p) - 1.0) / 14.5
	//fmt.Printf("%3.f\n", p * 127.0)
	db := dB(p) + 6.0
	return db
}

func round(n float64) float64 {
	if (n - math.Floor(n)) >= 0.5 {
		return math.Ceil(n)
	} else {
		return math.Floor(n)
	}
}

func DBtoMIDI(db float64) uint8 {
	db = db - 6.0
	p := math.Pow(10.0, (db / 20.0))
	plog := math.Log10(p*14.5+1.0) / math.Log10(15.5)
	plog *= 127.0
	return uint8(round(plog))
}

func logTaper(b int) int {
	// 127 * (ln(x+1)^2) / (ln(127+1)^2)
	return int(127.0 * math.Pow(math.Log2(float64(b)+1.0), 2) / math.Pow(math.Log2(127.0+1.0), 2))
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
	DirtyGain uint8
	CleanGain uint8
	Fx        [5]FXConfig
}

type ControllerState struct {
	SceneIdx  int
	Scene     *Scene
	PrIdx     int
	Pr        *Program
	Amp       [2]AmpState
	AmpConfig [2]AmpConfig
}

type Scene struct {
	Name string
	Amp  [2]AmpState
}

type Program struct {
	Name      string
	Tempo     int
	Scenes    []*Scene
	AmpConfig [2]AmpConfig
}

type Controller struct {
	midi Midi

	DefaultAmpConfig [2]AmpConfig
	Programs         []*Program

	Curr ControllerState
	Prev ControllerState
}

func NewController(midi Midi) *Controller {
	return &Controller{
		midi: midi,
	}
}

func (c *Controller) Init() {
	c.Curr.PrIdx = 0
	c.ActivateProgram()
}

func (c *Controller) ActivateProgram() {
	curr := &c.Curr

	if curr.PrIdx >= len(c.Programs) {
		curr.PrIdx = 0
	}
	curr.SceneIdx = 0

	log.Println("activate program", curr.PrIdx+1, len(c.Programs))

	curr.Pr = c.Programs[curr.PrIdx]
	if curr.Pr != nil {
		curr.AmpConfig = curr.Pr.AmpConfig
	}

	// Activate scene:
	c.ActivateScene()
}

func (c *Controller) ActivateScene() {
	curr := &c.Curr

	log.Println("activate scene", curr.SceneIdx+1, len(curr.Pr.Scenes))

	curr.Scene = curr.Pr.Scenes[curr.SceneIdx]
	if curr.Scene != nil {
		curr.Amp = curr.Scene.Amp
	}
}

func (c *Controller) HandleFswEvent(ev FswEvent) (err error) {
	curr := &c.Curr

	if ev.State {
		// Handle footswitch press:
		switch ev.Fsw {
		case FswNext:
			curr.SceneIdx++
			if curr.SceneIdx >= len(curr.Pr.Scenes) {
				curr.PrIdx++
				if curr.PrIdx >= len(c.Programs) {
					curr.PrIdx = 0
				}
				curr.SceneIdx = 0

				// Activate program:
				c.ActivateProgram()
			} else {
				c.ActivateScene()
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

	// Send new MIDI state:
	return c.SendMidi()
}

func (c *Controller) SendMidi() error {
	curr := &c.Curr

	// Send MIDI diff:
	for a := 0; a < 2; a++ {
		// Change amp mode:
		if curr.Amp[a].Mode != c.Prev.Amp[a].Mode {
			// TODO
			c.midi.CC(axeMidiChannel, 0, 0)
		}
	}

	// Copy to prev state:
	c.Prev = c.Curr
	return nil
}
