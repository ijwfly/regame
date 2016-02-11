package game

import "math"

type Structure interface {
	GetHealth() int
	SetHealth(int)
}

type Gun interface {
	Fire(int64) *Bullet
	IsFire() bool
	SetStateFire()
	SetStateLazy()
}

type Unit struct {
	Id        int
	X         float32
	Y         float32
	SpeedX    float32
	SpeedY    float32
	Radius    float32
	Type      int
	Structure Structure
	Gun       Gun
}

var currentId = 0

func GetNextUnitId() int {
	currentId++
	return currentId
}

func NewUnit(X, Y, SpeedX, SpeedY float32) *Unit {
	return &Unit{GetNextUnitId(), X, Y, SpeedX, SpeedY, 25, EnemyType, nil, nil}
}

type Movable interface {
	GetX() float32
	GetY() float32
	GetSpeedX() float32
	GetSpeedY() float32
	GetType() int
	GetId() int
	GetRadius() float32
	Move(gameStep int64)
	GetStructure() Structure
	SetStructure(Structure)
	CanCollide() bool
	Collide(Movable)
	GetGun() Gun
	TimeToHit(Movable) (bool, float32)
}

func (u *Unit) GetX() float32 {
	return u.X
}

func (u *Unit) GetY() float32 {
	return u.Y
}

func (u *Unit) GetSpeedX() float32 {
	return u.SpeedX
}

func (u *Unit) GetSpeedY() float32 {
	return u.SpeedY
}

func (u *Unit) GetType() int {
	return u.Type
}

func (u *Unit) GetId() int {
	return u.Id
}

func (u *Unit) GetRadius() float32 {
	return u.Radius
}

func (u *Unit) Move(gameStep int64) {
	u.X = u.X + u.SpeedX*float32(gameStep)/1000
	u.Y = u.Y + u.SpeedY*float32(gameStep)/1000
}

func (u *Unit) GetStructure() Structure {
	return u.Structure
}

func (u *Unit) SetStructure(s Structure) {
	u.Structure = s
}

func (u *Unit) GetGun() Gun {
	return u.Gun
}

func (u *Unit) CanCollide() bool {
	return true
}

func (a *Unit) Collide(b Movable) {
	if a.GetStructure() != nil && b.GetStructure() != nil {
		aHealth := a.GetStructure().GetHealth()
		bHealth := b.GetStructure().GetHealth()
		if aHealth > 0 && bHealth > 0 {
			a.GetStructure().SetHealth(aHealth - bHealth)
			b.GetStructure().SetHealth(bHealth - aHealth)
		}
	}
}

func (a *Unit) TimeToHit(b Movable) (bool, float32) {
	if a.GetId() == b.GetId() {
		return false, 0
	}
	dx, dy := b.GetX()-a.GetX(), b.GetY()-a.GetY()
	dvx, dvy := b.GetSpeedX()-a.GetSpeedX(), b.GetSpeedY()-a.GetSpeedY()
	dvdr := dx*dvx + dy*dvy
	if dvdr > 0 {
		return false, 0
	}
	dvdv := dvx*dvx + dvy*dvy
	drdr := dx*dx + dy*dy
	sigma := a.Radius + b.GetRadius()
	d := dvdr*dvdr - dvdv*(drdr-sigma*sigma)
	if d < 0 {
		return false, 0
	}
	return true, -(dvdr + float32(math.Sqrt(float64(d)))) / dvdv
}
