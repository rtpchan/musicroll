package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/ebiten/v2"
	"iatearock.com/musicroll/dft"
)

// is music file supported
func IsMusicFile(path string) bool {
	suffix := []string{".mp3", ".ogg", ".wav"}
	for _, s := range suffix {
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}

// change path from mp3/wav/ogg to png
func ToPngPath(path string) string {
	i := strings.LastIndex(path, ".")
	return path[:i] + ".png"
}

func IsPngExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func LoadPng(path string) *ebiten.Image {
	pianoRollImgMu.Lock()
	defer pianoRollImgMu.Unlock()
	imgF, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	img, err := png.Decode(bytes.NewReader(imgF))
	if err != nil {
		log.Println(err)
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

// Analyse sound, to be run in a go routine
// use msg to pass message
func AnalyseSound(path string, spacing time.Duration, msg chan string, done chan bool) {
	var err error
	f, err := os.Open(path)
	// f, err := os.Open("/home/raymond/Music/Ed Sheeran - The Joker And The Queen [Official Lyric Video]-m9f4XtNj1Vg.mp3")
	if err != nil {
		return
	}
	defer f.Close()

	ext := filepath.Ext(path)
	// spacing:= time.Millisecond * 500
	var format beep.Format
	var streamer beep.StreamSeekCloser
	var k *dft.Keys
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	default:
		return
	}
	if err != nil {
		return
	}
	k = dft.NewKeys(format, streamer, path)
	k.SetSpacing(spacing)
	pianoRollImgHeight = k.GetImage().Bounds().Dy()
	// pianoRollImgY = keyboardImgY - float64(pianoRollImgHeight)

	pngPath := path[:strings.LastIndex(path, ".")] + ".png"
	current := time.Millisecond * 0
	count := 0
	for current < k.Len() {

		sp := k.Analyse(current)
		k.DrawStripe(sp, count)
		// update image
		pianoRollImg = ebiten.NewImageFromImage(k.GetImage())
		current += k.Spacing()
		count += 1
		msg <- fmt.Sprintf("%0.2f", current.Seconds()/k.Len().Seconds())
		if count%10 == 0 {
			k.SaveImage(pngPath)
		}
	}
	k.SaveImage(pngPath)
	done <- true
}

func UpdateAnalysisProgress(msg chan string, done chan bool) {
	for {
		select {
		case m := <-msg:
			pianoRollImgProgress = m
		case <-done:
			analysing = false
			pianoRollImgProgress = "Completed"
			return
		}
	}
}
