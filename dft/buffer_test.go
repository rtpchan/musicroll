package dft

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

func TestXxx(t *testing.T) {
	// InitNote()
	initDraw()
	// k, err := loadFile("/home/raymond/Music/Ed Sheeran - The Joker And The Queen [Official Lyric Video]-m9f4XtNj1Vg.mp3")
	k, err := loadFile("/home/raymond/Music/bach.mp3")

	if err != nil {
		t.Error(err)
	}

	tm := k.Len()
	if tm != time.Millisecond*8208+time.Minute*3 {
		t.Errorf("buf length = %v", tm)
	}
	spacing := time.Millisecond * 100
	// spacing := time.Millisecond * 100
	k.SetSpacing(spacing)

	// keymap := k.Analyse(time.Second * 5)
	// log.Println(keymap["C8"])
	// log.Println(k.combine[240000:240030])
	// if len(keymap) != 88 {
	// 	t.Errorf("want 88 keys output, got %v", len(keymap))
	// }
	// // log.Println(keymap["A0"])
	// if (keymap["A0"] - 0.042532) > 0.0001 {
	// 	t.Errorf("A0 should be 0.042532, got %v", keymap["A0"])
	// }
	// // log.Println(keymap["C4"])
	// if (keymap["C4"] - 0.042532) > 0.0001 {
	// 	t.Errorf("A0 should be 0.042532, got %v", keymap["C4"])
	// }
	// for k, v := range keymap {
	// 	log.Printf("%s, %0.4f", k, v)
	// }
	now := time.Now()
	sp := k.AnalyseAll()
	elapsed := time.Since(now)
	log.Println(time.Since(now).String())

	numStrip := int(tm/spacing) + 1
	if len(k.spectrum) != numStrip {
		t.Errorf("want 19 spectrums, got %d", numStrip)
	}

	log.Printf("done analysing: %d instant", numStrip)
	log.Printf(" %0.2f per instant", (elapsed.Seconds() / float64(numStrip)))

	imageHeight := numStrip * 10
	bg := DrawEmpty(imageHeight) // 10 pixel * num of strips

	// bg := DrawSpectrum(sp, 10)

	for i, spec := range sp {
		log.Println(i, spec["C4"])
		// 	img := DrawNewStripe(spec, 0.001, 10)
		// 	bg = DrawOnImage(bg, img, image.Point{0, imageHeight - ((i + 1) * 10)})
		// 	// err = Export(fmt.Sprintf("test%02d.png", i), img)
		// 	// if err != nil {
		// 	// 	t.Errorf("export png : %v", err)
		// 	// }
	}
	err = Export("bach.png", bg)
	if err != nil {
		t.Errorf("export png : %v", err)
	} else {
		log.Println("export image successfully")
	}

	// img := DrawNewStripe(keymap, 1.0, 10)
	// err = Export("test.png", img)
	// if err != nil {
	// 	t.Errorf("export png : %v", err)
	// }

}

func TestGonumFFTData(t *testing.T) {
	sample := []float64{1, 0, 2, 0, 4, 0, 2, 0}
	n := NoteDFT(sample, "test", 1, 8)
	if (n - 3./8.) > 0.000001 {
		t.Errorf("wanted 3/8. got %f", n)
	}
}

func loadFile(path string) (*Keys, error) {

	f, err := os.Open(path)
	// f, err := os.Open("/home/raymond/Music/Ed Sheeran - The Joker And The Queen [Official Lyric Video]-m9f4XtNj1Vg.mp3")
	if err != nil {
		return &Keys{}, errors.New("Error reading this file. (OpenFile) %s")
	}
	defer f.Close()
	// defer f.Close()   // do not close f, streamer will cose it
	ext := filepath.Ext(path)
	var format beep.Format
	var streamer beep.StreamSeekCloser
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
		return &Keys{}, errors.New("Error, file format not recognised. (OpenFile)")
	}
	if err != nil {
		return &Keys{}, errors.New("Error decoding this file. (OpenFile)")
	}
	keys := NewKeys(format, streamer, "test")
	keys.SetSpacing(time.Second)
	return keys, nil
}
