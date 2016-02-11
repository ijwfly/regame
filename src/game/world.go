package game

import (
	//"kdtree"
	"kdtree"
	"sort"
	"support"
)

type World struct {
	Height  float32
	Width   float32
	Movable map[int]Movable
}

func NewWorld() *World {
	return &World{Height: 1200, Width: 2000, Movable: make(map[int]Movable)}
}

func (w *World) AddMovable(m Movable) {
	w.Movable[m.GetId()] = m
}

func (w *World) removeOutBoundUnits(boundary float32) {
	for _, current := range w.Movable {
		x := current.GetX()
		y := current.GetY()
		if !(-boundary < x && x < float32(w.Width)+boundary &&
			-boundary < y && y < float32(w.Height)+boundary) {
			delete(w.Movable, current.GetId())
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

func insertUnitToKdTree(tree *kdtree.T, unit Movable) *kdtree.T {
	unitT := new(kdtree.T)
	unitT.Point = kdtree.Point{float64(unit.GetX()), float64(unit.GetY())}
	unitT.Data = unit
	return tree.Insert(unitT)
}

type UnitCollision struct {
	Unit Movable
	d    float32
}
type UnitsCollisions []UnitCollision

func (a UnitsCollisions) Len() int           { return len(a) }
func (a UnitsCollisions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UnitsCollisions) Less(i, j int) bool { return a[i].d < a[j].d }

func (w *World) makeCollisions(gameStep int64) {
	unitsTree := kdtree.New(nil)
	for i := range w.Movable {
		obj := w.Movable[i]
		obj_structure := obj.GetStructure()
		if obj_structure != nil && obj_structure.GetHealth() > 0 {
			unitsTree = insertUnitToKdTree(unitsTree, obj)
		}

	}

	for i := range w.Movable {
		obj := w.Movable[i]
		obj_structure := obj.GetStructure()
		if obj.CanCollide() && obj_structure != nil && obj_structure.GetHealth() > 0 {
			nearestNodes := unitsTree.InRange(kdtree.Point{float64(obj.GetX()), float64(obj.GetY())}, 100, nil)
			if len(nearestNodes) > 1 {
				unitsCollisions := make([]UnitCollision, 0)
				mainUnit := obj.(Movable)
				for _, node := range nearestNodes {
					nodeUnit := node.Data.(Movable)
					isCollision, d := mainUnit.TimeToHit(nodeUnit)
					if isCollision && d < float32(gameStep)/1000 {
						unitsCollisions = append(unitsCollisions, UnitCollision{nodeUnit, d})
					}
				}
				sort.Sort(UnitsCollisions(unitsCollisions))
				for _, collision := range unitsCollisions {
					mainUnit.Collide(collision.Unit)
				}
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
				w.AddMovable(Movable(NewEffect(obj.GetX(), obj.GetY())))
			}
		}

	}
}

func (w *World) GetUnitsArrayView() [][]interface{} {
	res := make([][]interface{}, len(w.Movable))
	currentIndex := 0
	for i := range w.Movable {
		unit := w.Movable[i]
		unitView := make([]interface{}, 7)
		unitView[0] = unit.GetId()
		unitView[1] = unit.GetType()
		unitView[2] = support.Round2(unit.GetX())
		unitView[3] = support.Round2(unit.GetY())
		unitView[4] = support.Round2(unit.GetSpeedX())
		unitView[5] = support.Round2(unit.GetSpeedY())
		unitView[6] = unit.GetRadius()
		if currentIndex < len(res) {
			res[currentIndex] = unitView
		} else {
			res = append(res, unitView)
		}
		currentIndex++
	}
	return res
}
