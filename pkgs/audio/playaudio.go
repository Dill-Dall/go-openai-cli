package audio

import (
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

func PlaySound(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, ready, err := oto.NewContext(d.SampleRate(), 2, 2)
	if err != nil {
		return err
	}
	<-ready

	p := c.NewPlayer(d)
	defer p.Close()

	p.Play()

	stopPlayback := false

	// Listen for space key
	if err := keyboard.Open(); err != nil {
		return err
	}
	defer keyboard.Close()

	for {
		if _, key, err := keyboard.GetKey(); err == nil {
			if key == keyboard.KeySpace || key == keyboard.KeyEnter {
				if stopPlayback {
					break
				}
				stopPlayback = true
			}
		}

		time.Sleep(time.Millisecond * 10)

		if !p.IsPlaying() || stopPlayback {
			break
		}
	}

	return nil
}
