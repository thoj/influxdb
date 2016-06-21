package logfmt

// IsTerminal determines if this file descriptor references a terminal.
// This function is not implemented on Windows and will always return false.
func IsTerminal(fd int) bool {
	return false
}
