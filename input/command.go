package input

type KeyCommandExecutor interface {
	Execute(dt float64)
}

type KeyCommandExecutorFunc func(dt float64)

func (c KeyCommandExecutorFunc) Execute(dt float64) {
	c(dt)
}
