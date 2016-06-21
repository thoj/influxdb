package logfmt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/apex/log"
)

var Default = New(os.Stderr)

const timeFormat = "2006/01/02 15:04:05 -0700"

const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

var Colors = [...]int{
	log.DebugLevel: gray,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

type field struct {
	Name  string
	Value interface{}
}

type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type Handler struct {
	UseColors bool
	mu        sync.Mutex
	io.Writer
}

func New(w io.Writer) *Handler {
	useColors := false
	if f, ok := w.(*os.File); ok {
		useColors = IsTerminal(int(f.Fd()))
	}
	return &Handler{
		UseColors: useColors,
		Writer:    w,
	}
}

func (h *Handler) HandleLog(e *log.Entry) error {
	now := time.Now()
	var color int
	if h.UseColors {
		color = Colors[e.Level]
	}
	level := Strings[e.Level]

	fields := make([]field, 0, len(e.Fields))
	for k, v := range e.Fields {
		if k == "service" {
			continue
		}
		fields = append(fields, field{k, v})
	}
	sort.Sort(byName(fields))

	var buf bytes.Buffer
	buf.WriteString(now.Format(timeFormat))
	if h.UseColors {
		fmt.Fprintf(&buf, " |\033[%dm%5s\033[0m| ", color, level)
	} else {
		fmt.Fprintf(&buf, " |%5s| ", level)
	}
	if service, ok := e.Fields["service"]; ok {
		fmt.Fprintf(&buf, "[%s] ", service)
	}
	buf.WriteString(e.Message)

	for _, f := range fields {
		if h.UseColors {
			fmt.Fprintf(&buf, " \033[%dm%6s\033[0m=%v", green, f.Name, f.Value)
		} else {
			fmt.Fprintf(&buf, " %6s=%v", f.Name, f.Value)
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintln(h.Writer, buf.String())
	return nil
}
