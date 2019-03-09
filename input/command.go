package input

type KeyCommandExecutor interface {
	Execute()
}

type KeyCommandExecutorFunc func()

func (c KeyCommandExecutorFunc) Execute() {
	c()
}
