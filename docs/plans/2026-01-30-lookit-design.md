# Lookit - Beautiful Code & File Browser

**Date:** 2026-01-30
**Author:** Benjamin Connelly
**Status:** Approved

## Overview

Lookit is a beautiful local development server for browsing code, markdown, and files. Born from frustration with VS Code's markdown preview and files downloading instead of viewing, lookit provides a polished, modern interface for exploring any codebase.

## Core Philosophy

- **Personal first:** Optimize for solo developer workflow
- **Beautiful by default:** Modern, polished UI inspired by GitHub, Vercel, Linear
- **Zero config:** Just run `lookit` and it works
- **Smart defaults:** Respects .gitignore, syntax highlights everything, handles binaries gracefully

## Project Structure

```
lookit/
├── src/
│   ├── index.js          # Main server
│   ├── fileHandler.js    # File type detection & routing
│   ├── templates/        # HTML templates
│   │   ├── base.js       # Base HTML with beautiful defaults
│   │   ├── directory.js  # Directory listing (grid/list view)
│   │   ├── code.js       # Code viewer (GitHub-style)
│   │   ├── markdown.js   # Markdown renderer
│   │   └── binary.js     # Binary file preview card
│   ├── styles.js         # Modern CSS (Tailwind-inspired)
│   └── utils.js          # Helper functions
├── bin/
│   └── lookit.js         # CLI entry point
├── test/
│   └── fixtures/         # Test files
├── docs/
│   └── plans/            # Design documents
├── package.json
├── README.md
├── LICENSE               # MIT
└── .gitignore
```

## File Type Handling

### Categories

1. **Directories** → Modern grid/list layout with file icons
2. **Markdown** (.md, .mdx) → Rendered HTML with syntax highlighting
3. **Code files** (40+ extensions) → Syntax highlighting with highlight.js
4. **Images** (.png, .jpg, .svg, .webp) → Native display with metadata
5. **Binary files** → Preview card (no rendering)
6. **Unknown** → Treat as binary

### Detection Flow

```
Request → File exists? → Detect type → Route to handler → Render template
```

### Binary File Handling

Display card with:
- File icon (🔒 binaries, 📦 archives, 📄 unknown)
- File name and extension
- File size (human-readable)
- Last modified date
- "Copy path" button (copies absolute path)
- "Open in file manager" button
- NO attempt to render content

### .gitignore Support

- Parse `.gitignore` in current directory and all parents
- Filter directory listings by default
- `--all` flag shows ignored files with visual indicator
- Uses `ignore` npm package (battle-tested)

## CLI Interface

### Command Syntax
```bash
lookit [directory] [options]
```

### Options
```
--port, -p <number>    Port to listen on (default: 7777) 🍀
--host, -h <address>   Host to bind to (default: 127.0.0.1)
--open, -o             Open browser automatically
--all, -a              Show all files including .gitignore entries
--no-https             Use HTTP only
--https-only           Require HTTPS, fail if no certs
--quiet, -q            Suppress startup messages
--help                 Show help
--version, -v          Show version
```

### Startup Output
```
🔍 lookit v1.0.0

📂 Serving: /Users/ben/projects/myapp
🌐 Address:  https://127.0.0.1:7777
🔒 Security: HTTPS (TLS)
👁️  Filters:  .gitignore enabled (use --all to show all)

Press Ctrl+C to stop.
```

## Design Principles

### Visual Design
- Modern, minimal aesthetic (GitHub, Vercel, Linear inspired)
- Beautiful typography and spacing
- Smooth transitions and interactions
- Mobile-responsive
- Dark theme by default (light mode later)

### Directory Listings
- Grid/list toggle view
- File icons based on type
- Sort options: name, size, date
- Search/filter box (live filtering)
- Breadcrumb navigation

## Technology Stack

- **Runtime:** Node.js (>=18.0.0)
- **Server:** Built-in HTTP/HTTPS
- **Markdown:** markdown-it + markdown-it-highlightjs
- **Syntax Highlighting:** highlight.js (40+ languages)
- **Binary Detection:** isbinaryfile package
- **Gitignore Parsing:** ignore package

## Dependencies

```json
{
  "markdown-it": "^14.0.0",
  "markdown-it-highlightjs": "^4.1.0",
  "highlight.js": "^11.9.0",
  "isbinaryfile": "^5.0.0",
  "ignore": "^5.3.0"
}
```

## Distribution

### npm Package
- Name: `lookit`
- Author: Benjamin Connelly
- License: MIT
- Keywords: code-viewer, markdown-viewer, file-browser, dev-server

### GitHub Repository
- URL: github.com/Benjamin-Connelly/lookit
- Topics: nodejs, developer-tools, code-viewer, markdown, file-browser
- Enable issues and discussions
- GitHub Actions for automated npm publish

## Roadmap

### v1.0 (Initial Release)
- [x] Directory browsing with .gitignore support
- [x] Markdown rendering
- [x] Code syntax highlighting (40+ languages)
- [x] Binary file handling
- [x] Image viewing
- [x] Beautiful, modern UI
- [x] HTTPS support

### v1.1 (Near Future)
- [ ] Audio/video player (HTML5 player with beautiful controls)
- [ ] Light/dark theme toggle
- [ ] File search improvements
- [ ] Performance optimizations for large repos

### v2.0 (Future)
- [ ] Git integration (branch info, commit details)
- [ ] Diff viewer
- [ ] File tree sidebar
- [ ] Configuration file support
- [ ] Plugin system

## Success Criteria

- ✅ Can run `lookit` in any directory and immediately browse
- ✅ All text files display beautifully with syntax highlighting
- ✅ Binary files don't crash or display garbage
- ✅ .gitignore respected by default
- ✅ UI is polished and pleasant to use
- ✅ Zero configuration required
- ✅ Published to npm and GitHub

## Migration from mdserve

1. Copy mdserve code to new repo
2. Rename binary and package
3. Refactor into modular structure
4. Add new features (binary handling, .gitignore)
5. Implement beautiful UI templates
6. Update documentation
7. Publish to npm

---

**Next Steps:** Set up git repository and create implementation plan
