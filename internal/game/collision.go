package game

import (
	"math"
	"regame/pkg/kdtree"
	"sort"
)

type CollisionDetector struct {
	RootNode     *kdtree.T
	MaxTimeToHit float32
}

func NewCollisionDetector(MaxTimeToHit float32) *CollisionDetector {
	return &CollisionDetector{kdtree.New(nil), MaxTimeToHit}
}

func insertUnitToKdTree(tree *kdtree.T, unit Movable) *kdtree.T {
	unitT := new(kdtree.T)
	unitT.Point = kdtree.Point{float64(unit.GetX()), float64(unit.GetY())}
	unitT.Data = unit
	return tree.Insert(unitT)
}

func (c *CollisionDetector) AddItem(m Movable) {
	c.RootNode = insertUnitToKdTree(c.RootNode, m)
}

func TimeToHit(a, b Movable) (bool, float32) {
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
	sigma := a.GetRadius() + b.GetRadius()
	d := dvdr*dvdr - dvdv*(drdr-sigma*sigma)
	if d < 0 {
		return false, 0
	}
	return true, -(dvdr + float32(math.Sqrt(float64(d)))) / dvdv
}

type UnitCollision struct {
	Unit Movable
	d    float32
}
type UnitsCollisions []UnitCollision

func (a UnitsCollisions) Len() int           { return len(a) }
func (a UnitsCollisions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UnitsCollisions) Less(i, j int) bool { return a[i].d < a[j].d }

func (c *CollisionDetector) GetSortedCollisionsFor(m Movable) []Movable {
	nearestNodes := c.RootNode.InRange(kdtree.Point{float64(m.GetX()), float64(m.GetY())}, 100, nil)
	if len(nearestNodes) > 1 {
		unitsCollisions := make([]UnitCollision, 0)
		for _, node := range nearestNodes {
			nearest := node.Data.(Movable)
			isCollision, d := TimeToHit(m, nearest)
			if isCollision && d <= c.MaxTimeToHit {
				unitsCollisions = append(unitsCollisions, UnitCollision{nearest, d})
			}
		}
		sort.Sort(UnitsCollisions(unitsCollisions))
		res := make([]Movable, len(unitsCollisions))
		for i := range unitsCollisions {
			res[i] = unitsCollisions[i].Unit
		}
		return res
	}
	return nil
}
