Output demo-apply.mp4
Set Width 1200
Set Height 720
Set FontSize 18
Set TypingSpeed 60ms
Set CursorBlink true

# セットアップ（録画の主役は gion apply）
Sleep 1s
Type 'export GION_ROOT="$(mktemp -d)"'
Enter
Sleep 1s

Type "gion init"
Enter
Sleep 2s

Type 'cp vhs/fixtures/demo-apply-gion.yaml "$GION_ROOT/gion.yaml"'
Enter
Sleep 1s

Type 'cd "$GION_ROOT"'
Enter
Sleep 0.5s

Type "clear"
Enter
Sleep 0.5s

# 事前に用意した gion.yaml を適用する
Type "cat gion.yaml | yq"
Enter
Sleep 2s

Type "gion apply"
Enter
Sleep 1s

Type "y"
Sleep 0.5s
Enter
Sleep 10s

# できあがり確認（最低限）
Type "ls"
Enter
Sleep 1s

Type "ls workspaces"
Enter
Sleep 1s

Type "ls workspaces/PROJ-123"
Enter
Sleep 1s

Type "ls workspaces/PROJ-456"
Enter
Sleep 1s

Type "cd workspaces/PROJ-123/gion && git status -sb"
Enter
Sleep 2s

Type "cd ../../PROJ-456/gion && git status -sb"
Enter
Sleep 2s
