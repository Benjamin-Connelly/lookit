//go:build integration

package remote

import (
	"testing"
)

// Integration tests require SSH access to a test server.
// Run with: go test -tags=integration ./internal/remote/ -v
//
// These tests expect:
// - SSH config for "command" host in ~/.ssh/config
// - ~/test_md directory on the remote with markdown files

const testHost = "command"
const testPath = "~/test_md"

func TestIntegration_Connect(t *testing.T) {
	target := Target{Host: testHost, Path: testPath}
	conn := NewConn(target)

	if err := conn.Connect(); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer conn.Close()

	if conn.State() != ConnConnected {
		t.Errorf("state = %v, want Connected", conn.State())
	}

	// Path should be resolved (~ expanded)
	resolved := conn.Target()
	if resolved.Path == testPath {
		t.Error("path should be expanded from ~")
	}
	if resolved.Path == "" {
		t.Error("resolved path is empty")
	}
	t.Logf("Resolved path: %s", resolved.Path)

	// SFTP should be available
	if conn.SFTP() == nil {
		t.Error("SFTP client should not be nil after connect")
	}
}

func TestIntegration_ConnectAndList(t *testing.T) {
	target := Target{Host: testHost, Path: testPath}
	conn := NewConn(target)

	if err := conn.Connect(); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer conn.Close()

	resolved := conn.Target()
	entries, err := conn.SFTP().ReadDir(resolved.Path)
	if err != nil {
		t.Fatalf("ReadDir(%s): %v", resolved.Path, err)
	}

	if len(entries) == 0 {
		t.Error("expected files in remote test_md directory")
	}

	for _, e := range entries {
		t.Logf("  %s (%d bytes)", e.Name(), e.Size())
	}
}

func TestIntegration_Reconnect(t *testing.T) {
	target := Target{Host: testHost, Path: testPath}
	conn := NewConn(target)

	if err := conn.Connect(); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Close and reconnect
	conn.Close()
	if conn.State() != ConnDisconnected {
		t.Errorf("state after close = %v, want Disconnected", conn.State())
	}

	if err := conn.Reconnect(); err != nil {
		t.Fatalf("Reconnect: %v", err)
	}
	defer conn.Close()

	if conn.State() != ConnConnected {
		t.Errorf("state after reconnect = %v, want Connected", conn.State())
	}

	// SFTP should work after reconnect
	resolved := conn.Target()
	entries, err := conn.SFTP().ReadDir(resolved.Path)
	if err != nil {
		t.Fatalf("ReadDir after reconnect: %v", err)
	}
	if len(entries) == 0 {
		t.Error("no files after reconnect")
	}
}

func TestIntegration_BuildSSHConfig(t *testing.T) {
	target := Target{Host: testHost, Path: testPath}
	conn := NewConn(target)

	sshCfg, err := conn.buildSSHConfig()
	if err != nil {
		t.Fatalf("buildSSHConfig: %v", err)
	}

	if sshCfg.User == "" {
		t.Error("SSH config user should not be empty")
	}
	if len(sshCfg.Auth) == 0 {
		t.Error("SSH config should have at least one auth method")
	}
	if sshCfg.HostKeyCallback == nil {
		t.Error("SSH config should have host key callback")
	}
	t.Logf("SSH user: %s, auth methods: %d", sshCfg.User, len(sshCfg.Auth))
}
