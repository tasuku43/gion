Output demo-issue.mp4
Set Width 800
Set Height 600
Set FontSize 18
Set TypingSpeed 60ms
Set CursorBlink true

# Run command
Sleep 1s
Type "gion manifest add"
Enter
Sleep 1s

# Select mode: issue
Type "issue"
Sleep 0.5s
Enter
Sleep 1s

# Select repo
Type "backend"
Sleep 0.5s
Enter
Sleep 2s

# Select issues
Type "10"
Sleep 0.5s
Enter
Sleep 0.5s
Type "11"
Sleep 0.5s
Enter
Sleep 0.5s

# Finish issue selection
Type "done"
Sleep 0.5s
Enter
Sleep 1s

# Set branch suffixes per issue
Type "-example-config"
Sleep 0.5s
Enter
Sleep 0.5s
Type "-graceful-shutdown"
Sleep 0.5s
Enter
Sleep 3s

# Confirm apply
Enter
Sleep 1s
Enter
Sleep 3s
