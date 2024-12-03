package logrus

type level int

const (
	Panic level = iota
	Fatal
	Error
	Warn
	Info
	Debug
	Trace
)

type Configs struct {
	Level         level
	CustomEntries map[string]any
}
