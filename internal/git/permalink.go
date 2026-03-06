package git

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

// PermalinkStyle determines the URL format for generated links.
type PermalinkStyle int

const (
	PermalinkGitHub PermalinkStyle = iota
	PermalinkGitLab
	PermalinkBitbucket
)

// Permalink generates a permanent URL to a file at a specific commit.
func (r *Repo) Permalink(filePath string, line int) (string, error) {
	remotes, err := r.Remotes()
	if err != nil {
		return "", err
	}

	url, ok := remotes["origin"]
	if !ok {
		return "", fmt.Errorf("no origin remote found")
	}

	hash, err := r.HeadHash()
	if err != nil {
		return "", err
	}

	style := detectStyle(url)
	baseURL := normalizeRemoteURL(url)

	switch style {
	case PermalinkGitHub:
		link := fmt.Sprintf("%s/blob/%s/%s", baseURL, hash, filePath)
		if line > 0 {
			link += fmt.Sprintf("#L%d", line)
		}
		return link, nil
	case PermalinkGitLab:
		link := fmt.Sprintf("%s/-/blob/%s/%s", baseURL, hash, filePath)
		if line > 0 {
			link += fmt.Sprintf("#L%d", line)
		}
		return link, nil
	default:
		return fmt.Sprintf("%s/src/%s/%s", baseURL, hash, filePath), nil
	}
}

// CopyPermalink generates a permalink and copies it to the clipboard.
func (r *Repo) CopyPermalink(filePath string, line int) (string, error) {
	link, err := r.Permalink(filePath, line)
	if err != nil {
		return "", err
	}

	if err := clipboard.WriteAll(link); err != nil {
		return link, fmt.Errorf("copying to clipboard: %w (link: %s)", err, link)
	}
	return link, nil
}

func detectStyle(url string) PermalinkStyle {
	if strings.Contains(url, "github.com") {
		return PermalinkGitHub
	}
	if strings.Contains(url, "gitlab") {
		return PermalinkGitLab
	}
	return PermalinkBitbucket
}

func normalizeRemoteURL(url string) string {
	// Convert SSH URLs to HTTPS
	url = strings.TrimSuffix(url, ".git")
	if strings.HasPrefix(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}
	return url
}
