package runner

//Runnable can be run by a runner
type Runnable interface {
	Run()
}

//Runner trigger a function
type Runner interface {
	Start(Runnable)
}

//AsyncRunner implements an asynchronous runner
type AsyncRunner struct {
}

//Start triggers a goroutine for the function given in parameter
func (runner *AsyncRunner) Start(runnable Runnable) {
	go func() {
		runnable.Run()
	}()
}
