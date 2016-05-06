package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jackc/pong"
	"github.com/nsf/termbox-go"
)

var options struct {
	width  float64
	height float64
	seed   int
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Float64Var(&options.width, "width", 800, "width of world")
	flag.Float64Var(&options.height, "height", 600, "height of world")
	flag.IntVar(&options.seed, "seed", -1, "seed")
	flag.Parse()

	if options.seed < 0 {
		options.seed = time.Now().Nanosecond()
	}

	var kpc keyboardPaddleController

	game, err := pong.NewGame(
		pong.Vec2D{X: float32(options.width), Y: float32(options.height)},
		[]pong.PaddleController{&kpc},
		int64(options.seed),
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = termbox.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(30 * time.Millisecond)
	lastTime := time.Now()
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC:
					return
				case ev.Key == termbox.KeyArrowUp:
					kpc.Up = true
				case ev.Key == termbox.KeyArrowDown:
					kpc.Down = true
				}
			}
		case currentTime := <-ticker.C:
			Render(game, os.Stdout)
			frameDuration := currentTime.Sub(lastTime)
			lastTime = currentTime
			game.Tick(frameDuration)
		}
	}
}

func Render(game *pong.Game, wr io.Writer) {
	ball := termbox.Cell{Ch: '*', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack}
	paddle := termbox.Cell{Ch: ']', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for _, b := range game.Balls {
		x := int(b.Pos.X / 10)
		y := int(b.Pos.Y / 20)
		termbox.SetCell(x, y, ball.Ch, ball.Fg, ball.Bg)
	}

	for _, p := range game.Players {
		left := int(p.Paddle.Pos.X / 10)
		top := int(p.Paddle.Pos.Y / 20)
		right := left + int(p.Paddle.Size.X/10)
		bottom := top + int(p.Paddle.Size.Y/20)

		for y := top; y < bottom; y++ {
			for x := left; x < right; x++ {
				termbox.SetCell(x, y, paddle.Ch, paddle.Fg, paddle.Bg)
			}
		}
	}

	termbox.Flush()
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
