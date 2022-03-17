package game

import (
	"regame/internal/support"
	"time"
)

const (
	FrameStep    int64   = 25
	MaximumUnits int     = 100
	BulletSpeed  float32 = 800
	PlayerSpeed  float32 = 100
	EnemySpeed   float32 = 80
)

type Game struct {
	World *World
	Step  int64
	Turn  int64
}

func NewGame() *Game {
	return &Game{NewWorld(), FrameStep, 0}
}

func (g *Game) Start() {
	//for i := 0; i < 100; i++ {
	//	enemy := NewRandomEnemyFromTop(50)
	//	g.World.AddMovable(Movable(enemy))
	//}
	for {
		g.turn()
		g.Turn++
	}
}

func (g *Game) AddPlayer() *PlayerShip {
	player := NewPlayerShip(900, 900, 0, 0)
	g.World.AddMovable(Movable(player))
	return player
}

func (g *Game) turn() {
	start := support.MakeTimestamp()

	g.makeTurn()

	end := support.MakeTimestamp()
	if diff := g.Step - (end - start); diff > 0 {
		time.Sleep(time.Duration(diff) * time.Millisecond)
	}
}

func (g *Game) makeTurn() {
	if g.Turn%20 == 0 {
		enemy := NewRandomEnemyFromTop(50)
		g.World.AddMovable(Movable(enemy))
	}

	g.World.makeFire(g.Step * g.Turn)
	g.World.makeCollisions(g.Step)
	g.World.removeDeadUnits()
	g.World.makeMove(g.Step)
	g.World.removeOutBoundUnits(200)

	g.World.UnitsArray = g.World.GetUnitsArrayView()
}
