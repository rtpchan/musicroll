package dft

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"sync"
	"time"

	"github.com/faiface/beep"
)

type Keys struct {
	filepath string
	mu       sync.Mutex // combine buffer lock
	s        beep.StreamSeekCloser
	buffer   *beep.Buffer // [sample][channel]float64
	combine  []float64    // each value is a sample, sum or left and right

	combineStart int // combine indexes
	combineEnd   int
	windowStart  int // analysis window indexes
	windowEnd    int

	windowSize int

	fileLength time.Duration
	spacing    time.Duration // how frequent is a DFT is performed
	dataMu     sync.Mutex
	spectrum   []map[string]float64
	progress   float64 // num specturm analysed out of whole file

	image   image.Image
	imageMu sync.Mutex
}

func NewKeys(f beep.Format, s beep.StreamSeekCloser, filepath string) *Keys {
	initDraw()
	k := &Keys{
		filepath:     filepath,
		combine:      []float64{},
		combineStart: 0,
		combineEnd:   0,
		windowStart:  0,
		windowEnd:    0,
	}
	k.buffer = beep.NewBuffer(f)
	k.s = s
	// hardcode 0.1 second window size
	k.windowSize = f.SampleRate.N(time.Millisecond * 100)
	k.windowEnd = k.windowSize
	for k.windowEnd > k.combineEnd {
		k.nextBuffer()
	}
	k.fileLength = f.SampleRate.D(s.Len())
	// k.spacing = time.Millisecond * 100
	k.SetSpacing(time.Second * 1)
	k.progress = 0.0
	return k
}

// AnalyseAll perform PianoDFT for the entire file
func (k *Keys) AnalyseAll() []map[string]float64 {
	totalDataPoints := float64((k.Len() / k.spacing) + 1)
	// log.Printf("expecting %f data points", totalDataPoints)
	current := time.Millisecond * 0
	stripCount := 0
	imageHeight := k.image.Bounds().Dy()
	for current < k.fileLength {
		// log.Println(current)
		sp := k.Analyse(current)
		strip := DrawNewStripe(sp, 0.001, 10)
		k.imageMu.Lock()
		k.image = DrawOnImage(k.image, strip,
			image.Point{0, imageHeight - ((stripCount + 1) * 10)})
		k.imageMu.Unlock()
		k.dataMu.Lock()
		k.spectrum = append(k.spectrum, sp)
		k.progress = float64(len(k.spectrum)) / totalDataPoints
		k.dataMu.Unlock()
		current += k.spacing
		stripCount += 1
	}
	return k.spectrum
}

// Analyse spectrum at time t
func (k *Keys) Analyse(t time.Duration) map[string]float64 {
	k.windowStart = k.buffer.Format().SampleRate.N(t)
	k.windowEnd = k.windowStart + k.windowSize
	for k.windowEnd > k.combineEnd {
		k.nextBuffer()
	}
	// log.Println(k.windowStart, k.windowEnd, k.combine[k.windowStart:k.windowStart+10])
	return PianoDFT(k.combine[k.windowStart:k.windowEnd],
		k.buffer.Format().SampleRate.N(time.Second))
}

// GetSpectrum return spectrum data
// `got` is the number of specturm already recieved, this function return
// new specturm not send before
func (k *Keys) GetSpectrum(got int) []map[string]float64 {
	k.dataMu.Lock()
	defer k.dataMu.Unlock()
	log.Println(got, len(k.spectrum))
	if got > len(k.spectrum) {
		return []map[string]float64{}
	}
	return k.spectrum[got:]
}

// Draw 1 strip to the final image, with 1 spectrum, and index of the strip
func (k *Keys) DrawStripe(spectrum map[string]float64, stripCount int) {
	imageHeight := k.image.Bounds().Dy()
	strip := DrawNewStripe(spectrum, 0.001, 10)
	k.imageMu.Lock()
	k.image = DrawOnImage(k.image, strip,
		image.Point{0, imageHeight - ((stripCount + 1) * 10)})
	k.imageMu.Unlock()
}

func (k *Keys) GetImage() image.Image {
	k.imageMu.Lock()
	defer k.imageMu.Unlock()
	return k.image
}

func (k *Keys) SaveImage(name string) {
	err := Export(name, k.image)
	if err != nil {
		log.Printf("export png : %v", err)
	}
}

func (k *Keys) GetImageByte() []byte {
	buf := new(bytes.Buffer)
	k.imageMu.Lock()
	err := png.Encode(buf, k.image)
	k.imageMu.Unlock()

	if err != nil {
		log.Printf("exporting image []byte error, %s", err)
	}
	return buf.Bytes()
}

func (k *Keys) Progress() float64 {
	return k.progress
}

func (k *Keys) Spacing() time.Duration {
	return k.spacing
}

// Append read entire file, it is slow
func (k *Keys) Append(s beep.StreamSeekCloser) {
	k.buffer.Append(s)
}

// Set Spacing for each discrete fourier transform
func (k *Keys) SetSpacing(t time.Duration) {
	k.spacing = t
	k.initImage()
}

func (k *Keys) Len() time.Duration {
	return k.fileLength
}

func (k *Keys) TrimBuffer() {
	// TODO trim used buffer
	k.mu.Lock()
	k.combine = []float64{}
	defer k.mu.Unlock()
}

func (k *Keys) nextBuffer() int {
	var samples [1024][2]float64
	n, ok := k.s.Stream(samples[:])
	if !ok {
		return 0
	}
	for _, sample := range samples[:n] {
		sum := 0.0
		for _, channel := range sample {
			sum += channel
		}
		k.mu.Lock()
		k.combine = append(k.combine, sum)
		k.combineEnd += 1
		k.mu.Unlock()
	}
	return n
}

func (k *Keys) initImage() {
	tm := k.Len()
	numStrip := int(tm/k.spacing) + 1
	imageHeight := numStrip * 10 // 10 pixel per strip
	k.imageMu.Lock()
	k.image = DrawEmpty(imageHeight) // 10 pixel * num of strips
	k.imageMu.Unlock()
}
