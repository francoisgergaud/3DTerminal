package runner

import "log"

//Runnable can be run by a runner
type Runnable interface {
	Run() error
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
		err := runnable.Run()
		if err != nil {
			log.Println(err)
		}
	}()
}
