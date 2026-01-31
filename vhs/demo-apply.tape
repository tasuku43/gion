Output demo-apply.gif
Set Width 800
Set Height 600
Set FontSize 18
Set TypingSpeed 60ms
Set CursorBlink true

# 事前に用意した gion.yaml を適用する
Type "cat ~/gion/gion.yaml | yq"
Enter
Sleep 2s

Hide
Type "clear"
Enter
Show

Type "gion apply"
Sleep 0.5s
Enter
Sleep 3s

Type "y"
Sleep 0.5s
Enter
Sleep 5s

# できあがり確認（最低限）
Type "gion manifest ls"
Enter
Sleep 1s

Type "tree ~/gion/workspaces -L 2"
Enter
Sleep 3s
