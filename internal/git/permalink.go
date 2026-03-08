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
	PermalinkGitea
	PermalinkCodeberg
)

// Permalink generates a permanent URL to a file at a specific commit.
func (r *Repo) Permalink(filePath string, line int) (string, error) {
	hash, err := r.HeadHash()
	if err != nil {
		return "", fmt.Errorf("getting HEAD for permalink: %w", err)
	}

	baseURL, style, err := r.remoteBaseURL()
	if err != nil {
		return "", err
	}

	return buildFileLink(baseURL, style, hash, filePath, line, 0), nil
}

// PermalinkForBranch generates a URL using a branch name instead of a commit hash.
func (r *Repo) PermalinkForBranch(filePath string, branch string, line int) (string, error) {
	baseURL, style, err := r.remoteBaseURL()
	if err != nil {
		return "", err
	}

	return buildFileLink(baseURL, style, branch, filePath, line, 0), nil
}

// PermalinkForRange generates a permanent URL to a file with a line range.
func (r *Repo) PermalinkForRange(filePath string, startLine, endLine int) (string, error) {
	hash, err := r.HeadHash()
	if err != nil {
		return "", fmt.Errorf("getting HEAD for permalink range: %w", err)
	}

	baseURL, style, err := r.remoteBaseURL()
	if err != nil {
		return "", err
	}

	return buildFileLink(baseURL, style, hash, filePath, startLine, endLine), nil
}

// FileURL generates a non-permanent URL using the current branch.
func (r *Repo) FileURL(filePath string) (string, error) {
	branch, err := r.Branch()
	if err != nil {
		return "", fmt.Errorf("getting branch for file URL: %w", err)
	}

	baseURL, style, err := r.remoteBaseURL()
	if err != nil {
		return "", err
	}

	return buildFileLink(baseURL, style, branch, filePath, 0, 0), nil
}

// CopyPermalink generates a permalink and copies it to the clipboard.
// Returns the link string even if clipboard is unavailable.
func (r *Repo) CopyPermalink(filePath string, line int) (string, error) {
	link, err := r.Permalink(filePath, line)
	if err != nil {
		return "", err
	}

	if clipErr := clipboard.WriteAll(link); clipErr != nil {
		// Clipboard unavailable -- return the link anyway.
		return link, nil
	}
	return link, nil
}

// remoteBaseURL returns the normalized base URL and detected style for the origin remote.
func (r *Repo) remoteBaseURL() (string, PermalinkStyle, error) {
	remotes, err := r.Remotes()
	if err != nil {
		return "", 0, fmt.Errorf("reading remotes: %w", err)
	}

	url, ok := remotes["origin"]
	if !ok {
		return "", 0, fmt.Errorf("no origin remote found")
	}

	style := detectStyle(url)
	baseURL := normalizeRemoteURL(url)
	return baseURL, style, nil
}

// buildFileLink constructs a file URL for the given forge style.
// If endLine > startLine, a line range fragment is appended.
func buildFileLink(baseURL string, style PermalinkStyle, ref, filePath string, startLine, endLine int) string {
	var link string

	switch style {
	case PermalinkGitLab:
		link = fmt.Sprintf("%s/-/blob/%s/%s", baseURL, ref, filePath)
	case PermalinkBitbucket:
		link = fmt.Sprintf("%s/src/%s/%s", baseURL, ref, filePath)
	case PermalinkGitea, PermalinkCodeberg:
		link = fmt.Sprintf("%s/src/commit/%s/%s", baseURL, ref, filePath)
	default: // GitHub
		link = fmt.Sprintf("%s/blob/%s/%s", baseURL, ref, filePath)
	}

	// Force code view for markdown files (rendered view hides line numbers)
	if isMarkdownFile(filePath) && startLine > 0 {
		link += "?plain=1"
	}

	link += lineFragment(style, startLine, endLine)
	return link
}

// lineFragment returns the URL fragment for line references.
func lineFragment(style PermalinkStyle, startLine, endLine int) string {
	if startLine <= 0 {
		return ""
	}

	switch style {
	case PermalinkBitbucket:
		// Bitbucket uses #lines-N or #lines-N:M
		if endLine > startLine {
			return fmt.Sprintf("#lines-%d:%d", startLine, endLine)
		}
		return fmt.Sprintf("#lines-%d", startLine)
	case PermalinkGitea, PermalinkCodeberg:
		// Gitea/Codeberg use #L1-L5
		if endLine > startLine {
			return fmt.Sprintf("#L%d-L%d", startLine, endLine)
		}
		return fmt.Sprintf("#L%d", startLine)
	default: // GitHub, GitLab
		// GitHub: #L1-L5, GitLab: #L1-5
		if endLine > startLine {
			if style == PermalinkGitLab {
				return fmt.Sprintf("#L%d-%d", startLine, endLine)
			}
			return fmt.Sprintf("#L%d-L%d", startLine, endLine)
		}
		return fmt.Sprintf("#L%d", startLine)
	}
}

func detectStyle(url string) PermalinkStyle {
	lower := strings.ToLower(url)
	switch {
	case strings.Contains(lower, "github.com"):
		return PermalinkGitHub
	case strings.Contains(lower, "gitlab"):
		return PermalinkGitLab
	case strings.Contains(lower, "codeberg.org"):
		return PermalinkCodeberg
	case strings.Contains(lower, "gitea"):
		return PermalinkGitea
	case strings.Contains(lower, "bitbucket"):
		return PermalinkBitbucket
	default:
		return PermalinkGitHub
	}
}

func isMarkdownFile(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown") ||
		strings.HasSuffix(lower, ".mdown") || strings.HasSuffix(lower, ".mkd")
}

func normalizeRemoteURL(url string) string {
	url = strings.TrimSuffix(url, ".git")

	// ssh://git@host:port/path or ssh://git@host/path
	if strings.HasPrefix(url, "ssh://") {
		url = strings.TrimPrefix(url, "ssh://")
		// Remove user@ prefix
		if at := strings.Index(url, "@"); at >= 0 {
			url = url[at+1:]
		}
		// Handle port: host:port/path -> host/path
		if colon := strings.Index(url, ":"); colon >= 0 {
			slash := strings.Index(url[colon:], "/")
			if slash >= 0 {
				// port exists before slash -- strip port
				url = url[:colon] + url[colon+slash:]
			}
		}
		return "https://" + url
	}

	// git://host/path
	if strings.HasPrefix(url, "git://") {
		url = strings.TrimPrefix(url, "git://")
		// Remove user@ prefix if present
		if at := strings.Index(url, "@"); at >= 0 {
			url = url[at+1:]
		}
		return "https://" + url
	}

	// git@host:path (SCP-style)
	if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")
		// Replace first colon with slash (host:org/repo -> host/org/repo)
		if colon := strings.Index(url, ":"); colon >= 0 {
			url = url[:colon] + "/" + url[colon+1:]
		}
		return "https://" + url
	}

	// Already HTTPS or HTTP
	return url
}
