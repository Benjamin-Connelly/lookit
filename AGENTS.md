# AGENTS.md — lookit build instructions

## Task Tracking

Use `bd` (Beads) for all task tracking. Before starting any work:

1. Run `bd ready` to see what tasks are available
2. Run `bd update <id> --claim` to claim a task before working on it
3. Run `bd update <id> --status done` when a task is complete
4. Never work on a task that is blocked by an open dependency

Do not create new tasks unless explicitly told to. Work only the tasks
assigned in the epic for the current stage.

## Build Rules

- No CGO. Pure Go only. Must cross-compile cleanly.
- No external web frameworks. stdlib `net/http` only.
- All errors handled explicitly. No panics.
- Idiomatic Go throughout -- exported types where appropriate.
- After every task, run `go build ./...` and confirm it compiles.
- After every epic, run `go test ./...` and fix all failures before
  marking the epic done.

## Project

GitHub: github.com/Benjamin-Connelly/lookit
Language: Go
Binary name: lookit
Config dir: ~/.config/lookit/
