package git

import (
	"fmt"
	"path/filepath"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Repo wraps a go-git repository with convenience methods.
type Repo struct {
	repo     *gogit.Repository
	root     string
	worktree *gogit.Worktree
}

// FileStatus represents the git status of a single file.
type FileStatus struct {
	Path     string
	Staging  StatusCode
	Worktree StatusCode
}

// StatusCode maps to go-git status codes.
type StatusCode byte

const (
	Unmodified StatusCode = ' '
	Modified   StatusCode = 'M'
	Added      StatusCode = 'A'
	Deleted    StatusCode = 'D'
	Renamed    StatusCode = 'R'
	Copied     StatusCode = 'C'
	Untracked  StatusCode = '?'
)

// Open opens the git repository at or above the given path.
func Open(path string) (*Repo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	repo, err := gogit.PlainOpenWithOptions(absPath, &gogit.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("opening repo: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("getting worktree: %w", err)
	}

	return &Repo{
		repo:     repo,
		root:     absPath,
		worktree: wt,
	}, nil
}

// IsRepo returns true if path is inside a git repository.
func IsRepo(path string) bool {
	_, err := Open(path)
	return err == nil
}

// Branch returns the current branch name.
func (r *Repo) Branch() (string, error) {
	head, err := r.repo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}

// Status returns the working tree status for all files.
func (r *Repo) Status() ([]FileStatus, error) {
	status, err := r.worktree.Status()
	if err != nil {
		return nil, err
	}

	var files []FileStatus
	for path, s := range status {
		files = append(files, FileStatus{
			Path:     path,
			Staging:  StatusCode(s.Staging),
			Worktree: StatusCode(s.Worktree),
		})
	}
	return files, nil
}

// Remotes returns the list of remote names and URLs.
func (r *Repo) Remotes() (map[string]string, error) {
	remotes, err := r.repo.Remotes()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, remote := range remotes {
		cfg := remote.Config()
		if len(cfg.URLs) > 0 {
			result[cfg.Name] = cfg.URLs[0]
		}
	}
	return result, nil
}

// HeadHash returns the HEAD commit hash.
func (r *Repo) HeadHash() (string, error) {
	head, err := r.repo.Head()
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}

// Tags returns all tag names.
func (r *Repo) Tags() ([]string, error) {
	tags, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}

	var names []string
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		names = append(names, ref.Name().Short())
		return nil
	})
	return names, err
}

// Root returns the repository root path.
func (r *Repo) Root() string {
	return r.root
}
