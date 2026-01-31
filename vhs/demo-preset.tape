Output demo-preset.mp4
Set Width 800
Set Height 720
Set FontSize 18
Set TypingSpeed 60ms
Set CursorBlink true

# Define a preset (multiple repos)
Sleep 1s
Type "gion manifest preset add sample-app"
Enter
Sleep 1s

# Select repos (pick a couple of fetched repos)
Enter
Sleep 0.5s
Enter
Sleep 0.5s
Type "done"
Sleep 0.5s
Enter
Sleep 1s

# Static shot: gion.yaml now contains the preset
Type "cd $GION_ROOT"
Sleep 0.5s
Enter
Sleep 0.5s
Type "cat gion.yaml | yq"
Sleep 0.5s
Enter
Sleep 4s

# Apply preset: show Plan with multiple worktree creations, then cancel
Type "gion manifest add --preset sample-app SAMPLE-123"
Sleep 0.5s
Enter
Sleep 5s

# Static shot / GIF: Plan output
Sleep 4s

# Cancel at the apply prompt (should roll back the workspace entry)
Type "n"
Sleep 0.5s
Enter
Sleep 2s
