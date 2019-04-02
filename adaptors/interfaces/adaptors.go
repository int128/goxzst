package adaptors

type Cmd interface {
	Run(args []string) int
}

type Logger interface {
	Logf(format string, v ...interface{})
}
