package input

type KeyCommand interface {
	Execute()
}

type KeyCommandFunc func()

func (c KeyCommandFunc) Execute() {
	c()
}

type MouseMoveCommand interface {
	Execute(xOffset, yOffset float64)
}

type MouseMoveCommandFunc func(xOffset, yOffset float64)

func (c MouseMoveCommandFunc) Execute(xOffset, yOffset float64) {
	c(xOffset, yOffset)
}

type MouseScrollCommand interface {
	Execute(offsetX, yOffset float64)
}

type MouseScrollCommandFunc func(xOffset, yOffset float64)

func (c MouseScrollCommandFunc) Execute(xOffset, yOffset float64) {
	c(xOffset, yOffset)
}

type MouseButtonCommand interface {
	Execute()
}

type MouseButtonCommandFunc func()

func (c MouseButtonCommandFunc) Execute() {
	c()
}
