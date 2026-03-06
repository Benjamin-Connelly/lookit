# Lookit

Zero-config local dev server for browsing code, markdown, and files. Dark theme, syntax highlighting (50+ languages), git-aware directory listings, HTTPS via mkcert.

**Port:** 7777 (auto-increments if in use, up to 10 retries)

## Tech Stack

- **Runtime:** Node.js >= 18, pure CommonJS (`require`), no transpilation
- **Dependencies:** markdown-it + markdown-it-highlightjs for markdown rendering, highlight.js for syntax highlighting, isbinaryfile for binary detection, ignore for .gitignore parsing
- **No framework** -- raw `http`/`https` server, no Express, no build step
- **CLI entry:** `bin/lookit.js` -> `src/index.js`

## Directory Structure

```
bin/lookit.js           # CLI entrypoint (shebang wrapper)
src/
  index.js              # Server setup, arg parsing, request routing, port auto-increment
  fileHandler.js        # Route requests to the right template (dir/code/markdown/binary)
  gitHandler.js         # Git status, branch info, per-file commit data
  commands.js           # Management commands: --list, --stop, --stop-all
  instanceManager.js    # Multi-instance tracking (register/unregister/list/clean)
  utils.js              # Arg parser, cert helpers, .gitignore filtering, escapeHtml
  styles.js             # CSS-in-JS for all page templates
  templates/
    base.js             # HTML shell (head, styles, layout)
    directory.js         # Directory listing with git badges, file icons
    code.js              # Syntax-highlighted source view
    markdown.js          # GitHub-style rendered markdown
    binary.js            # Preview card with download link
test/
  comprehensive-test.sh # Bash integration tests
  test-git-features.js  # Git integration tests
  test-security-fix.js  # Path traversal / security regression tests
  fixtures/             # Test data
```

## Architecture

Request flow: `handleRequest()` in index.js validates the path (traversal guard), stats the file, then delegates to `handleFile()` in fileHandler.js. That function decides the template based on file type: directory listing, markdown render, syntax-highlighted code, media passthrough, or binary download card.

All HTML is server-rendered via template functions in `src/templates/`. Styles live in `styles.js` as exported string constants -- no external CSS files.

Git integration shells out to `git` CLI (via `execFileSync` / `execFile`) -- no libgit2 or nodegit dependency.

## Conventions

- YAGNI -- only build what's needed. No speculative abstractions.
- Functional, declarative style. Small focused modules.
- No test framework -- tests are plain Node scripts and bash.
- Templates return HTML strings, no JSX or template engine.
- All paths validated against CWD to prevent traversal.

## Quick Reference

```bash
# Run locally (dev)
node bin/lookit.js                  # serve CWD on :7777
node bin/lookit.js --port 3000      # custom port
node bin/lookit.js --open           # serve and open browser

# Management
node bin/lookit.js --list           # show running instances
node bin/lookit.js --stop 7778      # stop instance on port
node bin/lookit.js --stop-all       # stop all instances

# Run as global CLI
npm install -g .                    # link globally
lookit                              # then use anywhere

# Test
bash test/comprehensive-test.sh
node test/test-git-features.js
node test/test-security-fix.js
```

## Gotchas

- No `npm test` script defined -- run test files directly (see above).
- Styles are JS strings in `src/styles.js`, not CSS files. Easy to miss when changing UI.
- Git operations use `execFileSync` which blocks the event loop. Fine for a personal tool, would need async refactor for concurrent users.
- HTTPS auto-configures if certs exist at `~/.config/lookit/localhost.pem` and `localhost-key.pem`. No flag needed -- just having the files enables it. Use `--no-https` to force HTTP.
- Port auto-increment means the server may not be on 7777 if something else claimed it. Check console output.
- `.gitignore` filtering is on by default; `--all` disables it.
