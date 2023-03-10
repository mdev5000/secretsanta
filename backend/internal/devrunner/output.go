package devrunner

type (
	Source string

	Status string

	OutputData struct {
		Source Source
		Status Status
		Err    error
		Output string
	}
)

const (
	Backend  Source = "backend"
	Frontend Source = "frontend"

	ErrorS    Status = "error"
	Compiling Status = "compiling"
	Loading   Status = "loading"
	Running   Status = "running"
	Data      Status = "data"
)

const ScreenClear = "\033c"
