package game

import (
	"math"
	"math/rand"
)

const (
	Effect     = 30
	BulletType = 20
	EnemyType  = 10
	PlayerType = 1

	LazyState = 0
	FireState = 1
)

//=========================================================================
type SimpleGun struct {
	State      int
	Unit       Movable
	ShootTime  int64
	ReloadTime int64
}

func (o *SimpleGun) IsFire() bool {
	return o.State == FireState
}

func (o *SimpleGun) SetStateFire() {
	o.State = FireState
}

func (o *SimpleGun) SetStateLazy() {
	o.State = LazyState
}

func (o *SimpleGun) Fire(currentTime int64) *Bullet {
	if o.ShootTime+o.ReloadTime < currentTime {
		bullet := NewBullet(o.Unit.GetX(), o.Unit.GetY()-o.Unit.GetRadius(), 0, -BulletSpeed)
		o.ShootTime = currentTime
		return bullet
	} else {
		return nil
	}
}

//=========================================================================
type EffectGun struct {
	*SimpleGun
}

func (o *EffectGun) Fire(currentTime int64) *Bullet {
	if o.ShootTime == 0 {
		o.ShootTime = currentTime
		return nil
	} else if o.ShootTime+o.ReloadTime < currentTime {
		o.Unit.SetStructure(Structure(&SimpleStructure{-1}))
		return nil
	} else {
		return nil
	}
}

//=========================================================================
type SimpleStructure struct {
	Health int
}

func (o *SimpleStructure) GetHealth() int {
	return o.Health
}

func (o *SimpleStructure) SetHealth(h int) {
	o.Health = h
}

//=========================================================================
type EnemyShip struct {
	*Unit
}

func (a *EnemyShip) Collide(b Movable) {
	if b.GetType() != a.GetType() {
		a.Unit.Collide(b)
	}
}

func (e *EnemyShip) GetExplosion() Movable {
	return Movable(NewEffect(e.GetX(), e.GetY(), 40))
}

func NewEnemyShip(X, Y, SpeedX, SpeedY float32) *EnemyShip {
	unit := NewUnit(X, Y, SpeedX, SpeedY)
	unit.Type = EnemyType
	unit.Structure = Structure(&SimpleStructure{1})
	return &EnemyShip{unit}
}

func NewRandomEnemyUnit(speed float32) *EnemyShip {
	unit := NewEnemyShip(900, 400, 0, 0)

	swap := float32(0.0)
	t := 2 * math.Pi * rand.Float64()
	u := rand.Float32() + rand.Float32()
	if u > 1 {
		swap = 2 - u
	} else {
		swap = u
	}
	unit.SpeedX = speed * swap * float32(math.Cos(t))
	unit.SpeedY = speed * swap * float32(math.Sin(t))

	return unit
}

func NewRandomEnemyFromTop(speed float32) *EnemyShip {
	unit := NewRandomEnemyUnit(speed)
	unit.Y = -100
	unit.X = rand.Float32()*900 + 400
	unit.SpeedY = float32(math.Abs(float64(unit.SpeedY)))
	return unit
}

//=========================================================================
type PlayerShip struct {
	*Unit
}

func NewPlayerShip(X, Y, SpeedX, SpeedY float32) *PlayerShip {
	unit := NewUnit(X, Y, SpeedX, SpeedY)
	unit.Type = PlayerType
	unit.Structure = Structure(&SimpleStructure{1})
	unit.Gun = Gun(&SimpleGun{LazyState, Movable(unit), 0, 100})
	return &PlayerShip{unit}
}

func (p *PlayerShip) GetExplosion() Movable {
	return Movable(NewEffect(p.GetX(), p.GetY(), 80))
}

//=========================================================================
type Bullet struct {
	*Unit
}

func NewBullet(X, Y, SpeedX, SpeedY float32) *Bullet {
	unit := NewUnit(X, Y, SpeedX, SpeedY)
	unit.Type = BulletType
	unit.Radius = 5
	unit.Structure = Structure(&SimpleStructure{1})
	return &Bullet{unit}
}

func (b *Bullet) CanCollide() bool {
	return false
}

//=========================================================================
type EffectBullet struct {
	*Unit
}

func NewEffect(X, Y, R float32) *EffectBullet {
	unit := NewUnit(X, Y, 0, 0)
	unit.Radius = R
	unit.Type = Effect
	unit.Gun = Gun(&EffectGun{&SimpleGun{FireState, Movable(unit), 0, 800}})
	return &EffectBullet{unit}
}

func (b *EffectBullet) CanCollide() bool {
	return false
}
