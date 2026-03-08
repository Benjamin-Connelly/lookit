package tui

import (
	"strings"
	"testing"
)

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()
	// Verify critical bindings exist
	if len(km.Quit.Keys()) == 0 {
		t.Error("Quit should have keys")
	}
	if len(km.Up.Keys()) == 0 {
		t.Error("Up should have keys")
	}
	if len(km.Down.Keys()) == 0 {
		t.Error("Down should have keys")
	}
	if len(km.Enter.Keys()) == 0 {
		t.Error("Enter should have keys")
	}
}

func TestVimKeyMap(t *testing.T) {
	km := VimKeyMap()
	// Vim should include j/k
	hasK := false
	for _, k := range km.Up.Keys() {
		if k == "k" {
			hasK = true
		}
	}
	if !hasK {
		t.Error("vim Up should include 'k'")
	}
}

func TestEmacsKeyMap(t *testing.T) {
	km := EmacsKeyMap()
	// Emacs should use ctrl+p/ctrl+n
	hasCtrlP := false
	for _, k := range km.Up.Keys() {
		if k == "ctrl+p" {
			hasCtrlP = true
		}
	}
	if !hasCtrlP {
		t.Error("emacs Up should include 'ctrl+p'")
	}

	hasCtrlN := false
	for _, k := range km.Down.Keys() {
		if k == "ctrl+n" {
			hasCtrlN = true
		}
	}
	if !hasCtrlN {
		t.Error("emacs Down should include 'ctrl+n'")
	}

	// Search should be ctrl+s
	hasCtrlS := false
	for _, k := range km.Search.Keys() {
		if k == "ctrl+s" {
			hasCtrlS = true
		}
	}
	if !hasCtrlS {
		t.Error("emacs Search should be 'ctrl+s'")
	}
}

func TestHelp(t *testing.T) {
	km := DefaultKeyMap()
	h := Help(km)

	if h == "" {
		t.Fatal("Help should not be empty")
	}
	if !strings.Contains(h, "Key Bindings") {
		t.Error("Help should contain header")
	}
	if !strings.Contains(h, "quit") {
		t.Error("Help should list quit binding")
	}
	if !strings.Contains(h, "Visual Mode") {
		t.Error("Help should contain Visual Mode section")
	}
	if !strings.Contains(h, "Preview Search") {
		t.Error("Help should contain Preview Search section")
	}
	if !strings.Contains(h, "Heading Jump") {
		t.Error("Help should contain Heading Jump section")
	}
	if !strings.Contains(h, "Marks") {
		t.Error("Help should contain Marks section")
	}
}
