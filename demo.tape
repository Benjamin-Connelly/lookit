Output demo/fur-demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 1280
Set Height 760
Set Theme "Catppuccin Mocha"
Set Padding 20
Set Framerate 30
Set TypingSpeed 40ms

# -------------------------------------------------------------
# Setup (hidden): build a fresh binary for the demo
# -------------------------------------------------------------
Hide
Type "export PATH=/usr/local/go/bin:$PATH"
Enter
Type "go build -o /tmp/fur-demo ./cmd/fur && clear"
Enter
Sleep 5s
Show

# -------------------------------------------------------------
# Launch fur on the demo directory
# -------------------------------------------------------------
Type "/tmp/fur-demo demo/"
Sleep 400ms
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 1. Browse the file tree (j/k)
# -------------------------------------------------------------
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Up
Sleep 400ms
Up
Sleep 700ms

# -------------------------------------------------------------
# 2. Fuzzy filter — `/` then typed query, Enter opens the match
# -------------------------------------------------------------
Type "/"
Sleep 300ms
Type "read"
Sleep 1s
Enter
Sleep 800ms
Enter
Sleep 2s
Down@60ms 5
Sleep 800ms
Down@60ms 5
Sleep 1.2s

# -------------------------------------------------------------
# 3. Link picker — `f` surfaces every link with j/k + Enter
# -------------------------------------------------------------
Type "f"
Sleep 1.5s
Down
Sleep 400ms
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 4. Scroll through architecture.md
# -------------------------------------------------------------
Down@60ms 10
Sleep 1.2s

# -------------------------------------------------------------
# 5. Table of contents — `t` opens the TOC panel, Enter jumps
# -------------------------------------------------------------
Type "t"
Sleep 1.8s
Down
Sleep 400ms
Down
Sleep 500ms
Enter
Sleep 1.2s
Escape
Sleep 500ms

# -------------------------------------------------------------
# 6. Navigation history — Backspace walks the breadcrumb stack
# -------------------------------------------------------------
Backspace
Sleep 1.8s

# -------------------------------------------------------------
# 7. Backlinks — `b` lists every file that references this one
# -------------------------------------------------------------
Type "b"
Sleep 2.5s
Escape
Sleep 500ms

# -------------------------------------------------------------
# 8. Data preview — CSV rendered as a table
# -------------------------------------------------------------
Escape
Sleep 300ms
Type "/"
Sleep 300ms
Type "metrics"
Sleep 900ms
Enter
Sleep 700ms
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 9. Data preview — JSON pretty-printed and highlighted
# -------------------------------------------------------------
Escape
Sleep 300ms
Type "/"
Sleep 300ms
Type "config"
Sleep 900ms
Enter
Sleep 700ms
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 10. Syntax-highlighted Go source (50+ languages via Chroma)
# -------------------------------------------------------------
Escape
Sleep 300ms
Type "/"
Sleep 300ms
Type "server"
Sleep 900ms
Enter
Sleep 700ms
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 11. Global heading jump — Ctrl+G fuzzy-picks across all files
# -------------------------------------------------------------
Ctrl+G
Sleep 500ms
Type "contrib"
Sleep 1s
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 12. Command palette — `:broken` surfaces broken link targets
# -------------------------------------------------------------
Type ":"
Sleep 400ms
Type "broken"
Sleep 1.2s
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 13. Full-text search — `/` then Tab toggles to content mode
# -------------------------------------------------------------
Escape
Sleep 300ms
Type "/"
Sleep 300ms
Tab
Sleep 400ms
Type "mermaid"
Sleep 1.6s
Enter
Sleep 900ms
Enter
Sleep 2.5s

# -------------------------------------------------------------
# 14. Theme cycling — Ctrl+T walks auto → dark → light
# -------------------------------------------------------------
Ctrl+T
Sleep 1.4s
Ctrl+T
Sleep 1.4s
Ctrl+T
Sleep 1.2s

# -------------------------------------------------------------
# 15. Help overlay — `?` toggles the keybinding reference
# -------------------------------------------------------------
Type "?"
Sleep 3s
Escape
Sleep 500ms

# Quit
Type "q"
Sleep 500ms
