package devrunner

type (
	Source string

	OutputData struct {
		Source Source
		Err    error
		Output string
	}
)

const (
	Backend  Source = "backend"
	Frontend Source = "frontend"
)

const ScreenClear = "\033c"
