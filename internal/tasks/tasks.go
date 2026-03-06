package tasks

import (
	"regexp"
	"strings"
)

// Task represents a TODO/FIXME extracted from markdown files.
type Task struct {
	File     string
	Line     int
	Text     string
	Checked  bool
	Priority string // from context or tags
}

var taskPattern = regexp.MustCompile(`^(\s*[-*]\s+\[([xX ])\]\s+)(.+)$`)

// Extract finds all task items in markdown content.
func Extract(filePath, content string) []Task {
	var tasks []Task
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		matches := taskPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		checked := matches[2] == "x" || matches[2] == "X"
		text := strings.TrimSpace(matches[3])

		tasks = append(tasks, Task{
			File:    filePath,
			Line:    i + 1,
			Text:    text,
			Checked: checked,
		})
	}

	return tasks
}

// Aggregate collects tasks from multiple files.
func Aggregate(fileContents map[string]string) []Task {
	var all []Task
	for path, content := range fileContents {
		all = append(all, Extract(path, content)...)
	}
	return all
}

// Pending returns only unchecked tasks.
func Pending(tasks []Task) []Task {
	var pending []Task
	for _, t := range tasks {
		if !t.Checked {
			pending = append(pending, t)
		}
	}
	return pending
}
