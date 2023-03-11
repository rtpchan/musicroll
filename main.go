package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iatearock/dango"
	"github.com/iatearock/dango/ui"
	"github.com/sqweek/dialog"
	"golang.org/x/image/font"
)

//go:embed assets/*
var dataFS embed.FS

var (
	screenWidth   int = 800
	screenHeight  int = 600
	vfs           *dango.FS
	game          *Game
	font18        font.Face
	buttonFile    *ui.Button
	buttonPlay    *ui.Button
	buttonPause   *ui.Button
	buttonBack    *ui.Button
	buttonForward *ui.Button
	buttonRewind  *ui.Button
	buttonAnalyse *ui.Button
	analysing     bool

	infoMsg              string
	musicPath            string
	pianoRollPath        string
	pianoRollExist       bool
	pianoRollImg         *ebiten.Image
	pianoRollImgOp       *ebiten.DrawImageOptions
	pianoRollImgMu       sync.Mutex
	pianoRollImgHeight   int
	pianoRollImgProgress string
	pianoRollImgY        float64

	keyboardImg   *ebiten.Image
	keyboardImgOp *ebiten.DrawImageOptions
	keyboardImgY  float64 = float64(screenHeight - 98 - 30)

	ac *AudioControl
)

type Game struct {
}

func (g *Game) Draw(screen *ebiten.Image) {

	text.Draw(screen, musicPath, font18, 10, 30, color.White)
	text.Draw(screen, fmt.Sprintf("Debug - %v", pianoRollExist),
		font18, 10, 110, color.White)

	// No Music file
	if musicPath == "" {
		buttonFile.Draw(screen)
	}

	// Music file loaded
	if ac != nil {
		text.Draw(screen,
			fmt.Sprintf("Time: %.2f / -%.2f",
				ac.Current().Seconds(),
				ac.length-ac.Current().Seconds()),
			font18, screenWidth-180, 100, color.White)
		if ac.IsPlaying() {
			buttonPause.Draw(screen)
		} else if ac.IsEnded() {
			buttonRewind.Draw(screen)
		} else {
			buttonPlay.Draw(screen)
		}
		buttonBack.Draw(screen)
		buttonForward.Draw(screen)
	}

	// Has piano roll image
	if pianoRollImg != nil && ac != nil {
		pianoRollImgOp.GeoM.Reset()
		ratio := ac.Current().Seconds() / ac.length
		deltaImgY := float64(pianoRollImgHeight) * (1.0 - ratio)
		pianoRollImgY = keyboardImgY - deltaImgY
		pianoRollImgOp.GeoM.Translate(0, pianoRollImgY)
		screen.DrawImage(pianoRollImg, pianoRollImgOp)
	} else if ac != nil && !analysing && pianoRollImg == nil {
		// No image, Has file, and image file not found
		buttonAnalyse.Draw(screen)
	}

	screen.DrawImage(keyboardImg, keyboardImgOp)

	// ebitenutil.DebugPrintAt(screen, infoMsg, 10, 460)
	text.Draw(screen, infoMsg, font18, 10, screenHeight-10, color.White)
}

func (g *Game) Layout(outsideWidth, ousideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	if musicPath == "" {
		// ======== NO FILE, NEED TO SELECT A FILE =========
		if buttonFile.IsJustReleased() {
			filename, err := dialog.File().Load()
			if err != nil {
				infoMsg = err.Error()
			} else {
				if IsMusicFile(filename) {
					musicPath = filename
					pianoRollPath = ToPngPath(filename)
					pianoRollExist = IsPngExist(pianoRollPath)
					if pianoRollExist {
						pianoRollImg = LoadPng(pianoRollPath)
						pianoRollImgHeight = pianoRollImg.Bounds().Dy()
						buttonAnalyse.SetActive(false)
					} else {
						buttonAnalyse.SetActive(true)
					}
				} else {
					infoMsg = "File type not supported."
					musicPath = ""
				}
			}
		}
	} else if ac == nil {
		// ===== has valid file, load audio control ======
		ac = NewAudioControl(musicPath)

	} else {
		// ================ HAS FILE ===========================
		// ===== got audio control ======
		// if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		// 	if ac.IsPlaying() {
		// 		ac.Pause()
		// 		isPlaying = false
		// 		buttonPlay.SetActive(true)
		// 		buttonPause.SetActive(false)
		// 	} else {
		// 		ac.Play()
		// 		isPlaying = true
		// 		buttonPlay.SetActive(false)
		// 		buttonPause.SetActive(true)
		// 	}
		// }
		if ac.IsPlaying() && buttonPause.IsJustPressed() {
			ac.Pause()
		} else if ac.IsEnded() && buttonRewind.IsJustReleased() {
			ac.player.Seek(time.Second * 0)
		} else if !ac.IsPlaying() && buttonPlay.IsJustPressed() {
			ac.Play()
		}

		if buttonBack.IsJustPressed() {
			ac.Back(time.Second * -5)
		}
		if buttonForward.IsJustPressed() {
			ac.Forward(time.Second * 5)
		}

		// ====== Analysis =======
		if buttonAnalyse.IsJustReleased() {
			// analysis sound file
			analysing = true
			msg := make(chan string)
			done := make(chan bool)
			go AnalyseSound(musicPath, time.Millisecond*100, msg, done)
			go UpdateAnalysisProgress(msg, done)
			buttonAnalyse.SetActive(false)
		}
		// On going analysis
		if analysing {
			infoMsg = fmt.Sprintf("Analysis Progress: %s", pianoRollImgProgress)
		}
	}

	return nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	vfs = dango.NewFS(dataFS)
	// var err error
	font18, _ = vfs.GetFontFace("assets/fonts/Roboto-Regular.ttf", 18, 72)
	bc := [4]color.Color{
		color.RGBA{170, 170, 170, 255}, // normal
		color.RGBA{205, 205, 205, 255}, // hover
		color.RGBA{120, 120, 120, 255}, // pressed
		color.RGBA{50, 50, 50, 255},    // disable
	}
	bi := ButtonImages(50, 30, bc)
	buttonFile = ui.NewButton(bi[0], bi[1], bi[2], bi[3], 30, 40)
	buttonFile.SetText("File", font18, color.Black)

	biA := ButtonImages(80, 30, bc)
	buttonAnalyse = ui.NewButton(biA[0], biA[1], biA[2], biA[3], 30, 40)
	buttonAnalyse.SetText("Analyse", font18, color.Black)
	buttonAnalyse.SetActive(false)

	imgPlay, _ := vfs.GetImage("assets/images/play_circle_FILL1_wght400_GRAD0_opsz48.png")
	imgPause, _ := vfs.GetImage("assets/images/pause_circle_FILL1_wght400_GRAD0_opsz48.png")
	imgBack, _ := vfs.GetImage("assets/images/replay_5_FILL1_wght400_GRAD0_opsz48.png")
	imgForward, _ := vfs.GetImage("assets/images/forward_5_FILL1_wght400_GRAD0_opsz48.png")
	imgRewind, _ := vfs.GetImage("assets/images/replay_FILL1_wght400_GRAD0_opsz48.png")
	buttonPlay = ui.NewButton(imgPlay, imgPlay, imgPlay, imgPlay, screenWidth-120, 10)
	buttonPause = ui.NewButton(imgPause, imgPause, imgPause, imgPause, screenWidth-120, 10)
	buttonBack = ui.NewButton(imgBack, imgBack, imgBack, imgBack, screenWidth-180, 10)
	buttonForward = ui.NewButton(imgForward, imgForward, imgForward, imgForward, screenWidth-60, 10)
	buttonRewind = ui.NewButton(imgRewind, imgRewind, imgRewind, imgRewind, screenWidth-120, 10)
	// buttonPause.SetActive(false)

	pianoRollImgOp = &ebiten.DrawImageOptions{}
	// pianoRollImgOp.GeoM.Translate(0, 0)
	// pianoRollImgOp.GeoM.Translate(0, pianoRollImgY)

	keyboardImg, _ = vfs.GetImage("assets/images/keyboard1.png")
	keyboardImgOp = &ebiten.DrawImageOptions{}
	keyboardImgOp.GeoM.Translate(0, keyboardImgY)
	//                       98 height of keyboard image
	//                       30 height of message bar

	analysing = false
	game = &Game{}
}

func main() {
	icon, err := vfs.GetImage("assets/images/logo-universal.png")
	if err != nil {
		log.Println(err)
	}
	ebiten.SetTPS(30)
	ebiten.SetWindowTitle("Musicroll")
	ebiten.SetWindowIcon([]image.Image{icon})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
