package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type BugReport struct {
	OS           string
	Architecture string
	GoVersion    string
	StackTrace   string
}

func GenerateBugReport() BugReport {
	return BugReport{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		StackTrace:   getStackTrace(),
	}
}

func getStackTrace() string {
	buf := new(bytes.Buffer)
	stack := make([]byte, 1024)
	for {
		n := runtime.Stack(stack, true)
		if n < len(stack) {
			stack = stack[:n]
			break
		}
		stack = make([]byte, 2*len(stack))
	}
	buf.Write(stack)
	return buf.String()
}

func (br BugReport) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("OS: %s\n", br.OS))
	sb.WriteString(fmt.Sprintf("Architecture: %s\n", br.Architecture))
	sb.WriteString(fmt.Sprintf("Go Version: %s\n", br.GoVersion))
	sb.WriteString(fmt.Sprintf("Stack Trace:\n%s\n", br.StackTrace))
	return sb.String()
}

func SendBugReport(report BugReport) error {
	cmd := exec.Command("curl", "-X", "POST", "https://example.com/bugreport", "-d", report.String())
	return cmd.Run()
}
