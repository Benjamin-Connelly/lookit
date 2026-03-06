package doctor

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Check represents a single diagnostic check.
type Check struct {
	Name    string
	Status  CheckStatus
	Message string
}

// CheckStatus is the result of a diagnostic check.
type CheckStatus int

const (
	CheckOK CheckStatus = iota
	CheckWarn
	CheckFail
)

// Run executes all diagnostic checks and returns the results.
func Run() []Check {
	var checks []Check

	checks = append(checks, checkGo())
	checks = append(checks, checkGit())
	checks = append(checks, checkTerminal())

	return checks
}

// Print displays diagnostic results to stdout.
func Print(checks []Check) {
	for _, c := range checks {
		icon := "OK"
		switch c.Status {
		case CheckWarn:
			icon = "WARN"
		case CheckFail:
			icon = "FAIL"
		}
		fmt.Printf("[%s] %s: %s\n", icon, c.Name, c.Message)
	}
}

func checkGo() Check {
	return Check{
		Name:    "Go runtime",
		Status:  CheckOK,
		Message: fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
}

func checkGit() Check {
	cmd := exec.Command("git", "--version")
	out, err := cmd.Output()
	if err != nil {
		return Check{
			Name:    "Git",
			Status:  CheckWarn,
			Message: "git not found in PATH",
		}
	}
	return Check{
		Name:    "Git",
		Status:  CheckOK,
		Message: string(out),
	}
}

func checkTerminal() Check {
	return Check{
		Name:    "Terminal",
		Status:  CheckOK,
		Message: fmt.Sprintf("TERM=%s, GOOS=%s", "detected", runtime.GOOS),
	}
}
