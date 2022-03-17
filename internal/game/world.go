package game

type World struct {
	Height     float32
	Width      float32
	Movable    map[int]Movable
	UnitsArray [][]interface{}
}

func NewWorld() *World {
	return &World{Height: 1200, Width: 2000, Movable: make(map[int]Movable)}
}

func (w *World) AddMovable(m Movable) {
	w.Movable[m.GetId()] = m
}

func (w *World) RemoveUnit(Id int) {
	delete(w.Movable, Id)
}

func (w *World) removeOutBoundUnits(boundary float32) {
	for i := range w.Movable {
		current := w.Movable[i]
		x := current.GetX()
		y := current.GetY()
		if !(-boundary < x && x < float32(w.Width)+boundary &&
			-boundary < y && y < float32(w.Height)+boundary) {
			w.RemoveUnit(current.GetId())
		}
	}
}

func (w *World) makeFire(time int64) {
	for i := range w.Movable {
		moveUnit := w.Movable[i]
		gun := moveUnit.GetGun()
		if gun != nil && gun.IsFire() {
			fire := gun.Fire(time)
			if fire != nil {
				w.AddMovable(Movable(fire))
			}
		}
	}
}

func (w *World) makeMove(step int64) {
	for i := range w.Movable {
		moveUnit := w.Movable[i]
		moveUnit.Move(step)
	}
}

func (w *World) makeCollisions(gameStep int64) {
	MaxTimeToHit := float32(gameStep) / 1000
	detector := NewCollisionDetector(MaxTimeToHit)
	for i := range w.Movable {
		obj := w.Movable[i]
		if obj.CanBeDestroyed() {
			detector.AddItem(obj)
		}

	}
	for i := range w.Movable {
		obj := w.Movable[i]
		if obj.CanCollide() && obj.CanBeDestroyed() {
			collideUnits := detector.GetSortedCollisionsFor(obj)
			for i := range collideUnits {
				obj.Collide(collideUnits[i])
			}
		}
	}
}

func (w *World) removeDeadUnits() {
	for i := range w.Movable {
		obj := w.Movable[i]
		obj_structure := obj.GetStructure()
		if obj_structure != nil && obj_structure.GetHealth() <= 0 {
			delete(w.Movable, obj.GetId())
			if obj.CanCollide() {
				w.AddMovable(obj.GetExplosion())
			}
		}

	}
}

func (w *World) GetUnitsArrayView() [][]interface{} {
	res := make([][]interface{}, len(w.Movable))
	currentIndex := 0
	for i := range w.Movable {
		unit := w.Movable[i]
		if currentIndex < len(res) {
			res[currentIndex] = unit.ToArray()
		} else {
			res = append(res, unit.ToArray())
		}
		currentIndex++
	}
	return res
}
