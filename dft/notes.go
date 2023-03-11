package dft

import (
	"math"
	"math/cmplx"
)

var noteName []string
var NoteFreq map[string]float64 // node frequency
var NoteCycle map[string]int    // number of wave cycle to be sample

func init() {
	noteName = []string{"A0", "As0", "B0",
		"C1", "Cs1", "D1", "Ds1", "E1", "F1", "Fs1", "G1", "Gs1", "A1", "As1", "B1",
		"C2", "Cs2", "D2", "Ds2", "E2", "F2", "Fs2", "G2", "Gs2", "A2", "As2", "B2",
		"C3", "Cs3", "D3", "Ds3", "E3", "F3", "Fs3", "G3", "Gs3", "A3", "As3", "B3",
		"C4", "Cs4", "D4", "Ds4", "E4", "F4", "Fs4", "G4", "Gs4", "A4", "As4", "B4",
		"C5", "Cs5", "D5", "Ds5", "E5", "F5", "Fs5", "G5", "Gs5", "A5", "As5", "B5",
		"C6", "Cs6", "D6", "Ds6", "E6", "F6", "Fs6", "G6", "Gs6", "A6", "As6", "B6",
		"C7", "Cs7", "D7", "Ds7", "E7", "F7", "Fs7", "G7", "Gs7", "A7", "As7", "B7",
		"C8",
	}
	NoteFreq = make(map[string]float64)
	for i, n := range noteName {
		num := float64(i + 1)
		NoteFreq[n] = math.Pow(2, (num-49)/12) * 440
	}
	NoteFreq["test"] = 1.0
	NoteCycle = make(map[string]int)
	for _, n := range noteName {
		NoteCycle[n] = 20
	}

}

// Sample size required given rate, freq and mode(no. of cycles)
func SampleSize(sampleRate int, freq float64, cycle int) int {
	wavelength := 1. / freq // second
	timePerPoint := 1. / float64(sampleRate)
	pointPerCycle := wavelength / timePerPoint // no. of Points, per cycle
	count := math.Ceil(pointPerCycle * float64(cycle))
	return int(count)
}

// DFT for a single `note` using `sample`, N number
// of sample, in the k mode. (actual sample length use adjusted by
// the k mode used, if k=3, then wavelength x 3 sample size used
func NoteDFT(sample []float64, note string, k float64, sampleRate int) float64 {
	freq := NoteFreq[note]
	// NFloat := wavelength / (1 / float64(sampleRate)) * k
	NSample := int(math.Floor(float64(sampleRate) * k / freq)) // number of samples in k number of wave
	// math.Min - so we can accept lack of sample at the end of file
	// N := int(math.Min(math.Floor(NFloat), float64(len(sample))))

	// log.Printf("freq %s %f", note, freq)
	// log.Printf("sample rate %d", sampleRate)
	// log.Printf("nSample %d, k %f", NSample, k)
	if NSample > len(sample) {
		NSample = len(sample)
		k = numWave(NSample, note, sampleRate)
	}

	// Euler formula
	sum := complex(0, 0)
	for n := 0; n < NSample; n++ {
		X := complex(sample[n], 0) *
			cmplx.Exp(complex(0, (-2*math.Pi/float64(NSample))*float64(k)*float64(n)))
		sum += X
	}
	return cmplx.Abs(sum) / float64(NSample)
}

func PianoDFT(sample []float64, sampleRate int) map[string]float64 {
	noteValue := map[string]float64{}
	// nCycle := float64(len(sample)) / (float64(sampleRate) / NoteFreq["A0"])
	for _, key := range noteName {
		v := NoteDFT(sample, key, 25, sampleRate)
		// v := NoteDFT(sample, key, NoteCycle[key], sampleRate)
		// v := NoteDFT(sample, key, nCycle, sampleRate)
		noteValue[key] = v
	}
	return noteValue
}

// Number of wave that can fit in the sample for a note
func numWave(sampleSize int, note string, sampleRate int) float64 {
	freq := NoteFreq[note]
	samplePerWave := float64(sampleRate) / freq
	return float64(sampleSize) / samplePerWave
}

// ============== FFT =================
// log.Println("----- FFT -----")
// numFrames := 32000
// SampleRate := 44100
// fft := fourier.NewFFT(numFrames)
// coeff := fft.Coefficients(nil, l[0:numFrames]) // l is the waveform

// // mag := []float64{}
// // freq := []float64{}
// var maxFreq, magnitude, mean float64
// for i, c := range coeff {
// 	if i == 0 {
// 		// Zero cycle per period (over the sample)
// 		continue
// 	}
// 	// log.Println(i, c)
// 	m := cmplx.Abs(c)
// 	mean += m
// 	freqency := fft.Freq(i) * float64(SampleRate)
// 	// mag = append(mag, m)
// 	// freq = append(freq, freqency)

// 	log.Printf("%d - magnitude:%.02f, freq:%.02f ", i, m, freqency)
// 	// if m > magnitude {
// 	// 	magnitude = m
// 	// 	maxFreq = fft.Freq(i) * float64(numFrames)
// 	// 	log.Printf("--- magnitude:%.02f, freq:%.02f ", m, fft.Freq(i)*float64(numFrames))
// 	// }
// }
// log.Println(maxFreq, magnitude, mean/float64(numFrames))

// Plot(freq, mag, "Spectrum", "Freqency(Hz)", "Magnitude", "spectrum.png")
// seq := fft.Sequence(nil, coeff)
// for s := range seq {
// 	log.Println(s)

// }
