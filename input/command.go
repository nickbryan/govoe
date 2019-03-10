package input

// KeyCommandExecutor wraps the Execute command for key press callback.
//
// Execute will be called once per simulation if the criteria is met within the Manager.
type KeyCommandExecutor interface {
	Execute(dt float64)
}

// KeyCommandExecutorFunc allows a callback to satisfy the KeyCommandExecutor interface. You would use this
// over satisfying the KeyCommandExecutor within a struct when the KeyCommandExecutor has no state.
type KeyCommandExecutorFunc func(dt float64)

// Execute will be called once per simulation if the criteria is met within the Manager.
func (c KeyCommandExecutorFunc) Execute(dt float64) {
	c(dt)
}
