Output demo/fur-demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 20

# Build fur
Hide
Type "PATH=/usr/local/go/bin:$PATH go build -o /tmp/fur-demo ./cmd/fur"
Enter
Sleep 3s
Show

# Launch fur on demo directory
Type "/tmp/fur-demo demo/"
Sleep 500ms
Enter
Sleep 2s

# Navigate file list
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 400ms
Up
Sleep 400ms
Up
Sleep 1s

# Open a file
Enter
Sleep 2s

# Scroll the preview
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Sleep 1s

# Fuzzy filter
Type "/"
Sleep 500ms
Type "arch"
Sleep 1s
Enter
Sleep 2s

# Theme cycling
Type "t"
Sleep 1s
Type "t"
Sleep 1s
Type "t"
Sleep 1s

# Toggle TOC panel
Tab
Sleep 1.5s
Tab
Sleep 1s

# Exit
Type "q"
Sleep 500ms
