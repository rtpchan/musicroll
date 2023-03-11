package dft

// Generate a wave, by frequency, by note
// l := generateNote("A5", 44100, 8, 1)
// c := generateNote("C6", 44100, 8, 1)
// a0 := generateNote("A0", 44100, 8, 1)
// c8 := generateNote("C8", 44100, 8, 1)

// l = combineNotes(l, c)
// l = combineNotes(l, a0)
// l = combineNotes(l, c8)

import "math"

// generate sound wave of a note
func generateFreq(freq float64, sampleRate int, bitDepth int, second float64) []float64 {
	amplitude := math.Pow(2, float64(bitDepth)-1) - 1

	numSample := int(math.Round(float64(sampleRate) * second))
	note := make([]float64, numSample)
	for i := 0; i < numSample; i++ {
		note[i] = math.Sin(freq*2*math.Pi*float64(i)/float64(numSample)) * amplitude
	}
	return note
}

func generateNote(note string, sampleRate int, bitDepth int, second float64) []float64 {
	freq := NoteFreq[note]
	return generateFreq(freq, sampleRate, bitDepth, second)
}

// combine two notes into a single sound wave
func combineNotes(a, b []float64) []float64 {
	numSample := math.Min(float64(len(a)), float64(len(b)))
	note := make([]float64, int(numSample))
	for i := 0; i < int(numSample); i++ {
		note[i] = a[i] + b[i]
	}
	return note
}

// generate sin wave , with sampleRate/sec, and number of sample points
func generateSin(freq float64, sampleRate, sampleSize int) []float64 {
	note := make([]float64, sampleSize)
	for i := 0; i < sampleSize; i++ {
		note[i] = math.Sin(freq * 2 * math.Pi * float64(i) / float64(sampleRate))
	}
	return note
}

// generate cos wave , with sampleRate/sec, and number of sample points
func generateCos(freq float64, sampleRate, sampleSize int) []float64 {
	note := make([]float64, sampleSize)
	for i := 0; i < sampleSize; i++ {
		note[i] = math.Cos(freq * 2 * math.Pi * float64(i) / float64(sampleRate))
	}
	return note
}

// return sin and cos wave for each key
// use Euler formula for Fourier transform, don't need these sine waves
func WaveForm(sampleRate, sampleSize int) (map[string][]float64, map[string][]float64) {
	sins := map[string][]float64{}
	coss := map[string][]float64{}

	for _, n := range noteName {
		sins[n] = generateSin(NoteFreq[n], sampleRate, sampleSize)
		coss[n] = generateCos(NoteFreq[n], sampleRate, sampleSize)
	}
	return sins, coss
}
