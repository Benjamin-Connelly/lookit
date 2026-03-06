package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
}

func TestRegisterAndRun(t *testing.T) {
	r := NewRegistry()
	called := false
	r.Register(Hook{
		Name:  "test-hook",
		Point: HookBeforeRender,
		Fn: func(ctx *HookContext) error {
			called = true
			ctx.Content = "modified"
			return nil
		},
	})

	ctx := &HookContext{Content: "original"}
	if err := r.Run(HookBeforeRender, ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !called {
		t.Error("hook was not called")
	}
	if ctx.Content != "modified" {
		t.Errorf("expected modified content, got %q", ctx.Content)
	}
}

func TestRunEmptyPoint(t *testing.T) {
	r := NewRegistry()
	ctx := &HookContext{Content: "unchanged"}
	if err := r.Run(HookAfterRender, ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if ctx.Content != "unchanged" {
		t.Error("content should not change with no hooks")
	}
}

func TestMakeHookFn(t *testing.T) {
	hc := HookConfig{
		Prepend: "HEADER\n",
		Append:  "\nFOOTER",
		Replace: []ReplaceRule{{Old: "foo", New: "bar"}},
	}
	fn := makeHookFn(hc)
	ctx := &HookContext{Content: "hello foo world"}
	if err := fn(ctx); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(ctx.Content, "HEADER\n") {
		t.Error("expected HEADER prepend")
	}
	if !strings.HasSuffix(ctx.Content, "\nFOOTER") {
		t.Error("expected FOOTER append")
	}
	if !strings.Contains(ctx.Content, "bar") || strings.Contains(ctx.Content, "foo") {
		t.Error("expected foo replaced with bar")
	}
}

func TestParseHookPoint(t *testing.T) {
	tests := []struct {
		input   string
		want    HookPoint
		wantErr bool
	}{
		{"beforerender", HookBeforeRender, false},
		{"before_render", HookBeforeRender, false},
		{"afterrender", HookAfterRender, false},
		{"after_render", HookAfterRender, false},
		{"beforeindex", HookBeforeIndex, false},
		{"afterindex", HookAfterIndex, false},
		{"onnavigate", HookOnNavigate, false},
		{"on_navigate", HookOnNavigate, false},
		{"invalid", 0, true},
	}
	for _, tt := range tests {
		got, err := parseHookPoint(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseHookPoint(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("parseHookPoint(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestLoadPluginsNoDir(t *testing.T) {
	dir := t.TempDir()
	// No plugins dir should return empty registry
	r, err := LoadPlugins(dir)
	if err != nil {
		t.Fatalf("LoadPlugins: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestLoadPluginsFromFile(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	os.MkdirAll(pluginDir, 0o755)

	yaml := `name: test-plugin
description: A test plugin
hooks:
  - point: before_render
    prepend: "<!-- injected -->\n"
`
	os.WriteFile(filepath.Join(pluginDir, "test.yaml"), []byte(yaml), 0o644)

	r, err := LoadPlugins(dir)
	if err != nil {
		t.Fatalf("LoadPlugins: %v", err)
	}

	ctx := &HookContext{Content: "<p>hello</p>"}
	r.Run(HookBeforeRender, ctx)
	if !strings.HasPrefix(ctx.Content, "<!-- injected -->") {
		t.Errorf("expected prepend, got %q", ctx.Content)
	}
}
