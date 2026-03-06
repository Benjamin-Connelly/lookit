package tasks

import (
	"strings"
	"testing"
)

func TestExtract(t *testing.T) {
	content := `# Project
- [x] Done task
- [ ] Pending task
- [ ] !high Important task
- [ ] Task with #tag1 #tag2
- [ ] Task @due(2024-06-15)
Normal paragraph, not a task.
`
	tasks := Extract("test.md", content)

	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	// Checked
	if !tasks[0].Checked {
		t.Error("first task should be checked")
	}
	if tasks[1].Checked {
		t.Error("second task should not be checked")
	}

	// Priority
	if tasks[2].Priority != "high" {
		t.Errorf("expected priority high, got %q", tasks[2].Priority)
	}

	// Tags
	if len(tasks[3].Tags) != 2 || tasks[3].Tags[0] != "tag1" {
		t.Errorf("expected tags [tag1, tag2], got %v", tasks[3].Tags)
	}

	// Due date
	if tasks[4].DueDate != "2024-06-15" {
		t.Errorf("expected due 2024-06-15, got %q", tasks[4].DueDate)
	}
}

func TestPending(t *testing.T) {
	all := []Task{
		{Checked: true, Text: "done"},
		{Checked: false, Text: "pending1"},
		{Checked: false, Text: "pending2"},
	}
	p := Pending(all)
	if len(p) != 2 {
		t.Errorf("expected 2 pending, got %d", len(p))
	}
}

func TestGroupByFile(t *testing.T) {
	tasks := []Task{
		{File: "a.md", Text: "t1"},
		{File: "a.md", Text: "t2"},
		{File: "b.md", Text: "t3"},
	}
	groups := GroupByFile(tasks)
	if len(groups["a.md"]) != 2 {
		t.Error("expected 2 tasks in a.md")
	}
	if len(groups["b.md"]) != 1 {
		t.Error("expected 1 task in b.md")
	}
}

func TestGroupByTag(t *testing.T) {
	tasks := []Task{
		{Text: "t1", Tags: []string{"bug", "urgent"}},
		{Text: "t2", Tags: []string{"bug"}},
		{Text: "t3"},
	}
	groups := GroupByTag(tasks)
	if len(groups["bug"]) != 2 {
		t.Errorf("expected 2 bug tasks, got %d", len(groups["bug"]))
	}
	if len(groups["urgent"]) != 1 {
		t.Errorf("expected 1 urgent task, got %d", len(groups["urgent"]))
	}
	if len(groups["untagged"]) != 1 {
		t.Errorf("expected 1 untagged task, got %d", len(groups["untagged"]))
	}
}

func TestFormatTable(t *testing.T) {
	tasks := []Task{
		{File: "a.md", Line: 1, Text: "Fix bug", Priority: "high"},
		{File: "b.md", Line: 5, Text: "Add feature", Checked: true},
	}
	out := FormatTable(tasks)
	if !strings.Contains(out, "Fix bug") {
		t.Error("table should contain task text")
	}
	if !strings.Contains(out, "Total: 2 tasks") {
		t.Error("table should show total")
	}
}

func TestFormatTableEmpty(t *testing.T) {
	out := FormatTable(nil)
	if out != "No tasks found." {
		t.Errorf("expected 'No tasks found.', got %q", out)
	}
}

func TestAggregate(t *testing.T) {
	files := map[string]string{
		"a.md": "- [ ] task a",
		"b.md": "- [ ] task b\n- [x] task c",
	}
	all := Aggregate(files)
	if len(all) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(all))
	}
}
