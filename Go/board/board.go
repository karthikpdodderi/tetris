package board

type Mover interface {
	Down() State
	Drop() State
	Left()
	Right()
	RotateLeft()
	RotateRight()
}

type Printer interface {
	Start()
	Refresh()
	Stop()
}
