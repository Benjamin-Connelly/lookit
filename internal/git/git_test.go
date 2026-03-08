package git

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func setupTestRepo(t *testing.T) (*Repo, string) {
	t.Helper()
	dir := t.TempDir()

	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(dir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test"), 0o644); err != nil {
		t.Fatal(err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := wt.Add("README.md"); err != nil {
		t.Fatal(err)
	}
	_, err = wt.Commit("initial commit", &gogit.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@test.com", When: time.Now()},
	})
	if err != nil {
		t.Fatal(err)
	}

	repoCacheMu.Lock()
	delete(repoCache, dir)
	repoCacheMu.Unlock()

	r, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	return r, dir
}

func TestOpen_ValidRepo(t *testing.T) {
	r, _ := setupTestRepo(t)
	if r == nil {
		t.Fatal("expected non-nil Repo")
	}
	if r.repo == nil {
		t.Fatal("expected non-nil underlying repo")
	}
}

func TestOpen_NotARepo(t *testing.T) {
	dir := t.TempDir()

	repoCacheMu.Lock()
	delete(repoCache, dir)
	repoCacheMu.Unlock()

	_, err := Open(dir)
	if err == nil {
		t.Fatal("expected error opening non-repo directory")
	}
}

func TestIsRepo(t *testing.T) {
	_, dir := setupTestRepo(t)
	if !IsRepo(dir) {
		t.Error("expected IsRepo=true for git directory")
	}

	nonGit := t.TempDir()
	repoCacheMu.Lock()
	delete(repoCache, nonGit)
	repoCacheMu.Unlock()

	if IsRepo(nonGit) {
		t.Error("expected IsRepo=false for non-git directory")
	}
}

func TestBranch(t *testing.T) {
	r, _ := setupTestRepo(t)
	branch, err := r.Branch()
	if err != nil {
		t.Fatal(err)
	}
	// go-git PlainInit defaults to "master"
	if branch != "master" {
		t.Errorf("expected branch 'master', got %q", branch)
	}
}

func TestBranchList(t *testing.T) {
	r, _ := setupTestRepo(t)
	branches, err := r.BranchList()
	if err != nil {
		t.Fatal(err)
	}
	if len(branches) == 0 {
		t.Fatal("expected at least one branch")
	}
	found := false
	for _, b := range branches {
		if b == "master" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'master' in branch list, got %v", branches)
	}
}

func TestStatus_Clean(t *testing.T) {
	r, _ := setupTestRepo(t)
	files, err := r.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Errorf("expected clean status (0 files), got %d files", len(files))
	}
}

func TestStatus_Modified(t *testing.T) {
	r, dir := setupTestRepo(t)

	// Modify the committed file
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Modified"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := r.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("expected at least one modified file")
	}

	found := false
	for _, f := range files {
		if f.Path == "README.md" && f.Worktree == Modified {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected README.md with worktree Modified status, got %+v", files)
	}
}

func TestIsClean(t *testing.T) {
	r, dir := setupTestRepo(t)

	clean, err := r.IsClean()
	if err != nil {
		t.Fatal(err)
	}
	if !clean {
		t.Error("expected clean repo after commit")
	}

	// Dirty it
	if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}

	clean, err = r.IsClean()
	if err != nil {
		t.Fatal(err)
	}
	if clean {
		t.Error("expected dirty repo after adding untracked file")
	}
}

func TestLog(t *testing.T) {
	r, _ := setupTestRepo(t)
	commits, err := r.Log(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if commits[0].Message != "initial commit" {
		t.Errorf("expected message 'initial commit', got %q", commits[0].Message)
	}
	if commits[0].Author != "Test" {
		t.Errorf("expected author 'Test', got %q", commits[0].Author)
	}
	if len(commits[0].Hash) != 40 {
		t.Errorf("expected 40-char hash, got %q", commits[0].Hash)
	}
}

func TestHeadHash(t *testing.T) {
	r, _ := setupTestRepo(t)
	hash, err := r.HeadHash()
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) != 40 {
		t.Errorf("expected 40-char hex hash, got %q (len=%d)", hash, len(hash))
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{40}$`, hash)
	if !matched {
		t.Errorf("hash %q does not match hex pattern", hash)
	}
}

func TestTags_Empty(t *testing.T) {
	r, _ := setupTestRepo(t)
	tags, err := r.Tags()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Errorf("expected no tags in fresh repo, got %v", tags)
	}
}

func TestRemotes_Empty(t *testing.T) {
	r, _ := setupTestRepo(t)
	remotes, err := r.Remotes()
	if err != nil {
		t.Fatal(err)
	}
	if len(remotes) != 0 {
		t.Errorf("expected no remotes in local repo, got %v", remotes)
	}
}

func TestRoot(t *testing.T) {
	r, dir := setupTestRepo(t)
	root := r.Root()
	if root != dir {
		t.Errorf("expected root %q, got %q", dir, root)
	}
}
