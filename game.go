package pong

import (
	"errors"
	// "errors"
	"math/rand"
	"time"
)

const ballRadius = 8
const paddleSpeed = 300

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
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
}

func (a Rect) Intersect(b Rect) bool {
	return !(a.Right < b.Left || b.Right < a.Left || a.Bottom < b.Top || b.Bottom < a.Top)
}

type Ball struct {
	Pos Vec2D
	Vel Vec2D
}

func (b Ball) BoundingRect() Rect {
	return Rect{Left: b.Pos.X, Top: b.Pos.Y, Right: b.Pos.X + ballRadius*2, Bottom: b.Pos.Y + ballRadius*2}
}

type Paddle struct {
	Pos  Vec2D
	Size Vec2D
}

func (p Paddle) BoundingRect() Rect {
	return Rect{Left: p.Pos.X, Top: p.Pos.Y, Right: p.Pos.X + p.Size.X, Bottom: p.Pos.Y + p.Size.Y}
}

type Player struct {
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
	game.Balls[0] = game.newBall()

	if len(controllers) > 2 {
		return nil, errors.New("must receive at most 2 controllers")
	}

	for len(controllers) < 2 {
		controllers = append(controllers, nil)
	}

	game.Players = make([]Player, len(controllers))
	game.Players[0] = Player{
		Paddle:     Paddle{Pos: Vec2D{X: 40, Y: game.Dimensions.Y * 0.5}, Size: Vec2D{X: 20, Y: 60}},
		Controller: controllers[0],
	}
	game.Players[1] = Player{
		Paddle:     Paddle{Pos: Vec2D{X: game.Dimensions.X - 40, Y: game.Dimensions.Y * 0.5}, Size: Vec2D{X: 20, Y: 60}},
		Controller: controllers[1],
	}

	for i := range game.Players {
		player := &game.Players[i]
		if player.Controller == nil {
			player.Controller = &AIPaddleController{game: game, player: player}
		}
	}

	return game, nil
}

func (game *Game) newBall() Ball {
	vel := Vec2D{X: 100 + game.prng.Float32()*200, Y: 100 + game.prng.Float32()*200}
	if game.prng.Float32() >= 0.5 {
		vel.X = -vel.X
	}
	if game.prng.Float32() >= 0.5 {
		vel.Y = -vel.Y
	}

	return Ball{
		Pos: Vec2D{X: game.Dimensions.X * 0.5, Y: game.Dimensions.Y * 0.5},
		Vel: vel,
	}
}

func (game *Game) Tick(duration time.Duration) bool {
	t := float32(duration.Seconds())

	for i := range game.Players {
		p := &game.Players[i]
		input := p.Controller.Act()

		if input.Up {
			p.Paddle.Pos.Y -= paddleSpeed * t
		}
		if input.Down {
			p.Paddle.Pos.Y += paddleSpeed * t
		}
	}

	for i := range game.Balls {
		b := &game.Balls[i]
		deltaPos := b.Vel.MultScalar(t)
		b.Pos = b.Pos.Add(deltaPos)

		if b.Pos.X < 0 {
			game.Players[1].Score++
			*b = game.newBall()
		}
		if b.Pos.X > game.Dimensions.X-1 {
			game.Players[0].Score++
			*b = game.newBall()
		}
		if b.Pos.Y < 0 {
			b.Pos.Y = 0
			b.Vel.Y = -b.Vel.Y
		}
		if b.Pos.Y > game.Dimensions.Y-1 {
			b.Pos.Y = game.Dimensions.Y - 1
			b.Vel.Y = -b.Vel.Y
		}

		for _, p := range game.Players {
			if b.BoundingRect().Intersect(p.Paddle.BoundingRect()) {
				b.Vel.X = -b.Vel.X
				b.Pos.X = b.Pos.X - deltaPos.X
			}
		}
	}

	return true
}

type NilPaddleController int

func (npc NilPaddleController) Act() *PaddleInput {
	return &PaddleInput{}
}

type AIPaddleController struct {
	game   *Game
	player *Player
}

func (ai *AIPaddleController) Act() *PaddleInput {
	ball := ai.game.Balls[0]
	paddle := ai.player.Paddle
	paddleCenter := paddle.Pos.Y + (paddle.Size.Y * 0.5)
	ballCenter := ball.Pos.Y + ballRadius
	if ballCenter < paddleCenter {
		return &PaddleInput{Up: true}
	}
	return &PaddleInput{Down: true}
}
