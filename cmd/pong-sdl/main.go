package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/jackc/pong"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

var options struct {
	width  int
	height int
	seed   int
}

func init() {
	runtime.LockOSThread()
}

var font *ttf.Font

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&options.width, "width", 800, "width of world")
	flag.IntVar(&options.height, "height", 600, "height of world")
	flag.IntVar(&options.seed, "seed", -1, "seed")
	flag.Parse()

	if options.seed < 0 {
		options.seed = time.Now().Nanosecond()
	}

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		options.width, options.height, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("unable to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatalf("unable to create renderer: %v", err)
	}
	defer renderer.Destroy()

	err = ttf.Init()
	if err != nil {
		log.Fatalf("unable to initialize ttf subsystem: %v", err)
	}
	defer ttf.Quit()

	font, err = ttf.OpenFont("ArchivoBlack-for-Print/ArchivoBlack.ttf", 20)
	if err != nil {
		log.Fatalf("unable to open font: %v", err)
	}
	defer font.Close()

	var kpc keyboardPaddleController

	game, err := pong.NewGame(
		pong.Vec2D{X: float32(options.width), Y: float32(options.height)},
		[]pong.PaddleController{&kpc},
		int64(options.seed),
	)
	if err != nil {
		log.Fatalln(err)
	}

	lastFrameTime := time.Now()

	running := true
	for running {
		err = Render(game, renderer)
		if err != nil {
			log.Fatalln(err)
		}

		currentFrameTime := time.Now()
		frameDuration := currentFrameTime.Sub(lastFrameTime)
		lastFrameTime = currentFrameTime

		game.Tick(frameDuration)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch event := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyDownEvent:
				switch event.Keysym.Scancode {
				case sdl.SCANCODE_Q:
					running = false
				case sdl.SCANCODE_UP:
					kpc.Up = true
				case sdl.SCANCODE_DOWN:
					kpc.Down = true
				}
			}
		}
	}
}

func Render(game *pong.Game, renderer *sdl.Renderer) error {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	renderer.SetDrawColor(40, 40, 255, 255)

	err := renderText(0, 0, strconv.Itoa(game.Players[0].Score), renderer)
	if err != nil {
		return err
	}
	err = renderText(int32(options.width-150), 0, strconv.Itoa(game.Players[1].Score), renderer)
	if err != nil {
		return err
	}

	for _, b := range game.Balls {
		r := sdl.Rect{X: int32(b.Pos.X), Y: int32(b.Pos.Y), W: 16, H: 16}
		renderer.FillRect(&r)
	}

	for _, p := range game.Players {
		r := sdl.Rect{X: int32(p.Paddle.Pos.X), Y: int32(p.Paddle.Pos.Y), W: int32(p.Paddle.Size.X), H: int32(p.Paddle.Size.Y)}
		renderer.FillRect(&r)
	}

	renderer.Present()

	return nil
}

func renderText(x, y int32, text string, renderer *sdl.Renderer) error {
	// nop while render is failing
	return nil

	textSurf, err := font.RenderUTF8_Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return err
	}
	defer textSurf.Free()

	textText, err := renderer.CreateTextureFromSurface(textSurf)
	if err != nil {
		return err
	}
	defer textText.Destroy()

	renderer.Copy(textText, nil, &sdl.Rect{X: x, Y: y, W: textSurf.W, H: textSurf.H})

	return nil
}

type keyboardPaddleController struct {
	Up   bool
	Down bool
}

func (kpc *keyboardPaddleController) Act() *pong.PaddleInput {
	pi := &pong.PaddleInput{Up: kpc.Up, Down: kpc.Down}
	*kpc = keyboardPaddleController{}
	return pi
}
