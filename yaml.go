package axewitcher

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func (c *Controller) Load() error {
	// Read all bytes in from yaml file:
	bytes, err := ioutil.ReadFile("all_programs-v5.yml")
	if err != nil {
		return err
	}

	// Parse yaml:
	dict := make(map[string]interface{})
	err = yaml.Unmarshal(bytes, dict)
	if err != nil {
		return err
	}

	// Extract default amp configs:
	extractAmpConfig := func(ampsYaml []interface{}) [2]AmpConfig {
		ac := [2]AmpConfig{}
		for a := 0; a < 2; a++ {
			layoutYaml := ampsYaml[a].(map[interface{}]interface{})
			layout := layoutYaml["fx_layout"]
			// Parse fx names into MIDI CC:
			for n := 0; n < 5; n++ {
				name := layout.([]interface{})[n].(string)
				ac[a].Fx[n].Name = name
				ac[a].Fx[n].MidiCC = findFxMidiCC(name)
			}
		}
		return ac
	}
	c.DefaultAmpConfig = extractAmpConfig(dict["amp"].([]interface{}))

	// Extract programs:
	programsYaml := dict["programs"].([]interface{})
	c.Programs = make([]*Program, 0, len(programsYaml))
	for _, y := range programsYaml {
		yp := y.(map[interface{}]interface{})
		log.Println(yp)

		ac := c.DefaultAmpConfig
		if ampconfig, ok := yp["amp"]; ok {
			ac = extractAmpConfig(ampconfig.([]interface{}))
		}

		scenesYaml := yp["scenes"].([]interface{})
		scenes := make([]*Scene, 0, len(scenesYaml))
		for _, sy := range scenesYaml {
			syp := sy.(map[interface{}]interface{})
			//log.Println(syp)

			extractAmpState := func(ampConfig *AmpConfig, sypa map[interface{}]interface{}) AmpState {
				m := AmpMode(0)
				switch sypa["channel"].(string) {
				case "clean":
					m = AmpClean
					break
				case "dirty":
					m = AmpDirty
					break
				case "acoustic":
					m = AmpAcoustic
					break
				}

				g := 0
				if v, ok := sypa["gain_log"]; ok {
					g = logTaper(v.(int))
				}
				if v, ok := sypa["gain"]; ok {
					g = v.(int)
				}

				// Extract volume as dB:
				volumeDB := float64(0)
				if v, ok := sypa["level"]; ok {
					if fv, ok := v.(float64); ok {
						volumeDB = fv
					} else if iv, ok := v.(int); ok {
						volumeDB = float64(iv)
					}
				}
				volume := DBtoMIDI(volumeDB)

				ampState := AmpState{
					Mode:      m,
					DirtyGain: uint8(g),
					CleanGain: 0x20, // TODO
					Volume:    volume,
				}

				// Enable named fx:
				if v, ok := sypa["fx"]; ok {
					fx := v.([]interface{})
					for _, f := range fx {
						name := f.(string)
						// Enable the FX by name:
						for i, x := range ampConfig.Fx {
							if x.Name == name {
								ampState.Fx[i].Enabled = true
								break
							}
						}
					}
				}

				return ampState
			}

			scene := &Scene{
				Name: syp["name"].(string),
				Amp: [2]AmpState{
					extractAmpState(&ac[0], syp["MG"].(map[interface{}]interface{})),
					extractAmpState(&ac[1], syp["JD"].(map[interface{}]interface{})),
				},
			}
			log.Println(scene)

			scenes = append(scenes, scene)
		}

		p := &Program{
			Name:      yp["name"].(string),
			Tempo:     yp["tempo"].(int),
			AmpConfig: ac,
			Scenes:    scenes,
		}
		c.Programs = append(c.Programs, p)
	}

	return nil
}
