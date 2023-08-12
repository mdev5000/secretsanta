package svelte_check_wrapper

import (
	"bufio"
	"fmt"
	lg "github.com/charmbracelet/lipgloss"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Type string

const (
	Error     Type = "ERROR"
	Start     Type = "START"
	Completed Type = "COMPLETED"
)

type CheckMessage struct {
	Num            int
	Type           Type
	ErrorMessage   string
	Filename       string
	FileLineAndCol string
}

var red = lg.Color("#f00")
var lightGreen = lg.Color("#4cb848")
var lightOrange = lg.Color("#f9a51b")
var lightBlue = lg.Color("#b9e5fa")
var styleError = lg.NewStyle().Foreground(red)
var stylePath = lg.NewStyle().Foreground(lightGreen)
var stylePathLineAndCol = lg.NewStyle().Foreground(lightOrange)
var styleErrorMessage = lg.NewStyle().Foreground(lightBlue)

func (m CheckMessage) Print(w io.Writer) {
	var styledType string
	switch m.Type {
	case Error:
		styledType = styleError.Render(string(m.Type))
	default:
		styledType = string(m.Type)
	}
	fmt.Fprintf(w, "\n%d %s \"%s:%s\"\n%s\n\n", m.Num,
		styledType,
		stylePath.Render(m.Filename),
		stylePathLineAndCol.Render(m.FileLineAndCol),
		styleErrorMessage.Render(m.ErrorMessage))
}

func Passed(msgs []CheckMessage) bool {
	for _, msg := range msgs {
		if msg.Type != Start && msg.Type != Completed {
			return false
		}
	}
	return true
}

type Filter = func(message CheckMessage) bool

func And(filter ...Filter) Filter {
	return func(message CheckMessage) bool {
		for _, f := range filter {
			if !f(message) {
				return false
			}
		}
		return true
	}
}

// IgnoreNodeModules ignores any errors from the node_modules directory.
func IgnoreNodeModules(m CheckMessage) bool {
	return !strings.HasPrefix(m.Filename, "node_modules")
}

// IgnoreDataTestIdMessage ignores missing data-testid parameters in svelte-check results.
func IgnoreDataTestIdMessage(m CheckMessage) bool {
	testIdMatcher := regexp.MustCompile(`^Type '{.*"data-testid": string;.*is not assignable to type`)
	return !testIdMatcher.MatchString(m.ErrorMessage)
}

func ParseLines(r io.Reader, filter Filter) []CheckMessage {
	var out []CheckMessage
	s := bufio.NewScanner(r)
	for s.Scan() {
		msg, ok := ParseLine(s.Text())
		if !ok || !filter(msg) {
			continue
		}
		out = append(out, msg)
	}
	return out
}

func ParseLine(line string) (CheckMessage, bool) {
	numAndRest := strings.SplitN(line, " ", 2)
	if len(numAndRest) != 2 {
		return CheckMessage{}, false
	}
	num, err := strconv.Atoi(numAndRest[0])
	if err != nil {
		return CheckMessage{}, false
	}
	re := regexp.MustCompile("^(ERROR|START|COMPLETED) (.+)$")
	parts := re.FindStringSubmatch(numAndRest[1])
	if len(parts) == 0 {
		panic(fmt.Errorf("bad line '%s'", line))
	}
	var errorMessage string
	var filename string
	var fileLineAndCol string
	lineType := Type(parts[1])

	switch lineType {

	case Start:
		filename = parts[2][1 : len(parts[2])-1] // trim quotes

	case Error:
		errorMatcher := regexp.MustCompile(`^"([^"]+)" ([0-9:]+) "(.+)$`)
		errParts := errorMatcher.FindStringSubmatch(parts[2])
		if len(errParts) == 0 {
			panic(fmt.Errorf("bad error line '%s'", line))
		}
		filename = errParts[1]
		fileLineAndCol = errParts[2]
		rawMessage := errParts[3][:len(errParts[3])-1] // trim left over "
		errorMessage = strings.ReplaceAll(rawMessage, `\"`, `"`)
		errorMessage = strings.ReplaceAll(errorMessage, "\\n", "\n")

	}

	return CheckMessage{
		Num:            num,
		Filename:       filename,
		FileLineAndCol: fileLineAndCol,
		Type:           lineType,
		ErrorMessage:   errorMessage,
	}, true
}
