package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	mp3decoder "github.com/hajimehoshi/go-mp3"
)

type AudioControl struct {
	path   string // file path to music file
	ctx    *audio.Context
	player *audio.Player
	length float64 // length of music in second
}

// Create a new audio control with path to a music file
func NewAudioControl(path string) *AudioControl {
	// ac := &AudioControl{ctx: audio.NewContext(44100), path: path}
	var sampleRate int
	var length float64 // length of music in second
	var ctx *audio.Context
	var s *mp3.Stream
	var err error
	if strings.HasSuffix(path, ".mp3") {
		b, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
			return nil
		}
		reader := bytes.NewReader(b)
		d, err := mp3decoder.NewDecoder(reader)
		if err != nil {
			log.Println(err)
			return nil
		}

		// length in second
		// go-mp3 doc say sample always convert to 16bit, 2 channels i.e. 4 bytes
		// TODO but is the length before the convertion?
		length = float64(d.Length()) / float64(d.SampleRate()) / 4.0
		sampleRate = d.SampleRate()

		ctx = audio.NewContext(sampleRate)
		s, err = mp3.DecodeWithSampleRate(ctx.SampleRate(), reader)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	ac := &AudioControl{ctx: ctx}
	ac.player, err = ac.ctx.NewPlayer(s)
	ac.path = path
	ac.length = length
	if err != nil {
		log.Println(err)
		return nil
	}
	return ac
}

func (ac *AudioControl) Play() {
	ac.player.Play()
}

func (ac *AudioControl) Pause() {
	ac.player.Pause()
}

func (ac *AudioControl) Current() time.Duration {
	return ac.player.Current()
}

func (ac *AudioControl) IsPlaying() bool {
	return ac.player.IsPlaying()
}

func (ac *AudioControl) IsEnded() bool {
	return (ac.length - ac.Current().Seconds()) < 0.1
}

func (ac *AudioControl) Length() time.Duration {
	return time.Duration(ac.length * 1e9)
}

// Foward by offset amount of time
func (ac *AudioControl) Forward(offset time.Duration) {

	c := ac.Current() + offset
	// log.Println(offset, ac.Current(), c, time.Duration(ac.length*1e9))
	if c > time.Duration(ac.length*1e9) {
		c = time.Duration(ac.length * 1e9)
	}
	ac.player.Seek(c)
}

// Back by offset amount of time
func (ac *AudioControl) Back(offset time.Duration) {
	c := ac.Current() + offset
	if c < time.Duration(0) {
		c = 0
	}
	ac.player.Seek(c)
}
