# Changelog

## v1.0.0

Initial stable release of fur — a dual-mode markdown navigator with TUI and web interfaces.

### Features
- **TUI mode**: Split-pane Bubble Tea interface with fuzzy file search, markdown rendering (Glamour), syntax highlighting (Chroma), and inter-document link navigation
- **Web mode**: stdlib `net/http` server with Goldmark markdown rendering, SSE live reload, security headers, ETag caching
- **MCP server**: Model Context Protocol server exposing 5 tools for AI agent integration
- **Remote browsing**: SSH/SFTP support with ssh-agent, key files, and `~/.ssh/config` integration
- **Link graph**: Bidirectional link tracking with backlinks, broken link detection, and DOT/JSON graph output
- **Full-text search**: Bleve-based search with BM25 scoring, plus fuzzy filename matching
- **Task extraction**: TODO/checkbox extraction with priority markers, tags, and due dates
- **Plugin system**: YAML-based hooks for content transformation (prepend/append/replace)
- **50+ language highlighting**: Chroma-powered syntax highlighting in both TUI and web modes
- **Git integration**: go-git for status, branches, log, and permalink generation (GitHub/GitLab/Bitbucket/Gitea/Codeberg)
- **Man pages**: Embedded man page installer for all subcommands
- **Shell completions**: Bash, Zsh, and Fish completion generation
- **Per-project config**: `.fur.toml`/`.fur.yaml` with automatic discovery (walks up from CWD)
- **Environment diagnostics**: `fur doctor` with 9 checks and colored output

### Distribution
- Homebrew tap: `brew install Benjamin-Connelly/fur/fur`
- Nix flake: `nix run github:Benjamin-Connelly/fur`
- Go install: `go install github.com/Benjamin-Connelly/fur/cmd/fur@v1.0.0`
- Pure Go, no CGO — cross-compiles to linux/darwin on amd64/arm64

### Security
- Path traversal protection via `Index.ValidatePath()` (shared by web and MCP)
- Content Security Policy headers
- Input sanitization on all API endpoints
- No shell-outs (pure Go throughout)
