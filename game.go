package pong

import (
	"errors"
	// "errors"
	"math/rand"
	"time"
)

type Vec2D struct {
	X float32
	Y float32
}

func (a Vec2D) Add(b Vec2D) Vec2D {
	return Vec2D{X: a.X + b.X, Y: a.Y + b.Y}
}

func (v Vec2D) MultScalar(s float32) Vec2D {
	return Vec2D{X: v.X * s, Y: v.Y * s}
}

type Rect struct {
	Pos  Vec2D
	Size Vec2D
}

type Ball struct {
	Pos Vec2D
	Vel Vec2D
}

type Paddle struct {
	Pos  Vec2D
	Size Vec2D
}

type Goal struct {
	Rect Rect
}

type Player struct {
	Goal       Goal
	Paddle     Paddle
	Controller PaddleController
	Score      int
}

type PaddleInput struct {
	Up   bool
	Down bool
}

type PaddleController interface {
	Act() *PaddleInput
}

type Game struct {
	Balls        []Ball
	Players      []Player
	Dimensions   Vec2D
	WinningScore int
	prng         *rand.Rand
}

func NewGame(dim Vec2D, controllers []PaddleController, rngSeed int64) (*Game, error) {
	game := &Game{Dimensions: dim, prng: rand.New(rand.NewSource(rngSeed))}
	game.Balls = make([]Ball, 1)
	game.Balls[0] = Ball{
		Pos: Vec2D{X: game.Dimensions.X * 0.5, Y: game.Dimensions.Y * 0.5},
		Vel: Vec2D{X: 100 + game.prng.Float32()*200, Y: 100 + game.prng.Float32()*200},
	}

	if len(controllers) > 2 {
		return nil, errors.New("must receive at most 2 controllers")
	}

	for len(controllers) < 2 {
		controllers = append(controllers, NilPaddleController(0))
	}

	game.Players = make([]Player, len(controllers))
	game.Players[0] = Player{
		Goal:       Goal{Rect: Rect{Pos: Vec2D{X: 0, Y: 0}, Size: Vec2D{X: 20, Y: game.Dimensions.Y}}},
		Paddle:     Paddle{Pos: Vec2D{X: 40, Y: game.Dimensions.Y * 0.5}, Size: Vec2D{X: 20, Y: 60}},
		Controller: controllers[0],
	}
	game.Players[1] = Player{
		Goal:       Goal{Rect: Rect{Pos: Vec2D{X: game.Dimensions.X - 20, Y: 0}, Size: Vec2D{X: 20, Y: game.Dimensions.Y}}},
		Paddle:     Paddle{Pos: Vec2D{X: game.Dimensions.X - 40, Y: game.Dimensions.Y * 0.5}, Size: Vec2D{X: 20, Y: 60}},
		Controller: controllers[1],
	}

	return game, nil
}

func (game *Game) Tick(duration time.Duration) bool {
	t := float32(duration.Seconds())

	for i := range game.Balls {
		b := &game.Balls[i]
		deltaPos := b.Vel.MultScalar(t)
		b.Pos = b.Pos.Add(deltaPos)

		if b.Pos.X < 0 {
			b.Pos.X = 0
			b.Vel.X = -b.Vel.X
		}
		if b.Pos.X > game.Dimensions.X-1 {
			b.Pos.X = game.Dimensions.X - 1
			b.Vel.X = -b.Vel.X
		}
		if b.Pos.Y < 0 {
			b.Pos.Y = 0
			b.Vel.Y = -b.Vel.Y
		}
		if b.Pos.Y > game.Dimensions.Y-1 {
			b.Pos.Y = game.Dimensions.Y - 1
			b.Vel.Y = -b.Vel.Y
		}
	}

	for i := range game.Players {
		p := &game.Players[i]
		input := p.Controller.Act()

		if input.Up {
			p.Paddle.Pos.Y -= 200 * t
		}
		if input.Down {
			p.Paddle.Pos.Y += 200 * t
		}
	}

	return true
}

type NilPaddleController int

func (npc NilPaddleController) Act() *PaddleInput {
	return &PaddleInput{}
}
