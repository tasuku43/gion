Output demo.mp4
Set Width 1200
Set Height 720
Set FontSize 18
Set TypingSpeed 60ms
Set CursorBlink true

Sleep 1s
Type "gion manifest add"
Enter
Sleep 1s

Type "issue"
Sleep 0.5 s
Enter
Sleep 1s

Type "backend"
Sleep 0.5 s
Enter
Sleep 2s

Enter
Sleep 0.5s
Enter
Sleep 1s

Type "done"
Sleep 1s
Enter
Sleep 3s

Enter
Sleep 1s
Enter
Sleep 3s

Type "y"
Sleep 0.5s
Enter
Sleep 5s

Type 'eval "$(giongo init)"'
Sleep 0.5s
Enter
Sleep 1s
Type "giongo"
Sleep 0.5s
Enter
Sleep 1s

Type "iss"
Sleep 0.5s
Type "11"
Sleep 0.5s
Enter

Type "pwd"
Enter
Sleep 3s

Type 'echo "demo" >> README.md'
Sleep 0.5s
Enter

Type "cd $GION_ROOT"
Sleep 1s
Enter

Type "cat gion.yaml | yq"
Sleep 1s
Enter
Sleep 3s

Type "clear"
Enter

Type "gion manifest rm"
Sleep 1s
Enter
Sleep 1s
Enter
Sleep 1s
Enter
Sleep 1s

Type "done"
Sleep 1s
Enter
Sleep 3s

Type "y"
Sleep 0.5s
Enter
Sleep 3s
