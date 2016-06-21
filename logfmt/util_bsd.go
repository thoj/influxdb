// +build darwin dragonfly freebsd netbsd openbsd

package logfmt

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
