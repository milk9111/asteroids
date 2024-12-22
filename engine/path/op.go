package path

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
	dmath "github.com/yohamta/donburi/features/math"
)

type Op interface {
	Target() dmath.Vec2
	Do(path *vector.Path, targetModifier func(target dmath.Vec2) dmath.Vec2)
	Copy(target dmath.Vec2) Op
}

type op struct {
	target dmath.Vec2
}

func (op op) Target() dmath.Vec2 {
	return op.target
}

type moveToOp struct {
	op
}

func NewMoveToOp(target dmath.Vec2) Op {
	return moveToOp{
		op{
			target: target,
		},
	}
}

func (op moveToOp) Do(path *vector.Path, targetModifier func(target dmath.Vec2) dmath.Vec2) {
	t := targetModifier(op.target)
	path.MoveTo(float32(t.X), float32(t.Y))
}

func (op moveToOp) Copy(target dmath.Vec2) Op {
	return NewMoveToOp(target)
}

type lineToOp struct {
	op
}

func NewLineToOp(target dmath.Vec2) Op {
	return lineToOp{
		op{
			target: target,
		},
	}
}

func (op lineToOp) Do(path *vector.Path, targetModifier func(target dmath.Vec2) dmath.Vec2) {
	t := targetModifier(op.target)
	path.LineTo(float32(t.X), float32(t.Y))
}

func (op lineToOp) Copy(target dmath.Vec2) Op {
	return NewLineToOp(target)
}

type closeOp struct {
	op
}

func NewCloseOp() Op {
	return closeOp{}
}

func (op closeOp) Do(path *vector.Path, targetModifier func(target dmath.Vec2) dmath.Vec2) {
	path.Close()
}

func (op closeOp) Copy(target dmath.Vec2) Op {
	return NewCloseOp()
}
