package main

import (
	"log"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("unable to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatalf("unable to create renderer: %v", err)
	}
	defer renderer.Destroy()

	image, err := img.Load("sprite.png")
	if err != nil {
		log.Fatalf("unable to load PNG: %v", err)
	}
	defer image.Free()

	texture, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		log.Fatalf("unable to create texture: %v", err)
	}
	defer texture.Destroy()

	xPos := int32(0)
	xVel := int32(1)

	frameCount := 0
	lastFrameTime := time.Now()

	running := true
	for running {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		for i := int32(0); i < 20; i++ {
			for j := int32(0); j < 20; j++ {
				dst := sdl.Rect{xPos + j*25, i * 20, 19, 19}
				renderer.Copy(texture, nil, &dst)
			}
		}

		renderer.Present()
		frameCount++

		if frameCount%100 == 0 {
			currentFrameTime := time.Now()
			fps := 100.0 / currentFrameTime.Sub(lastFrameTime).Seconds()
			lastFrameTime = currentFrameTime
			log.Println("frameCount:", frameCount)
			log.Println("fps:", fps)
		}

		xPos += xVel

		if xPos < 0 {
			xPos = 0
			xVel = -xVel
		}
		if xPos > 600 {
			xPos = 600
			xVel = -xVel
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			log.Printf("%#v\n", event)

			switch event := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyDownEvent:
				switch event.Keysym.Scancode {
				case sdl.SCANCODE_Q:
					running = false
				}
			}
		}
	}
}
