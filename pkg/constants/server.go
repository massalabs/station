package constants

//go:generate stringer -type=Status
type Status uint8

const (
	Running Status = iota
	Stopping
	Stopped
)
