# UI Architecture

このドキュメントは、UI.md の契約を破らないための実装構造（責務分離）をまとめる。

## 目的
- セクション順序（Inputs → Info → Steps → Result → Suggestion）を常に保証する
- サブコマンド固有の描画を排除し、共通コンポーネントで契約を強制する
- 対話フローでも Inputs が重複出力されないようにする

## コンポーネント

### Frame (`internal/ui/frame.go`)
**責務**
- セクション順序と描画ルールを一元化する器
- `Inputs/Info/Steps/Result/Suggestion` の内容を保持し、順序固定で描画する

**使い方**
- `SetInputsPrompt(...)` でプロンプト行を設定
- `AppendInputsRaw(...)` でリスト/ツリー行（既に整形済み）を追加
- `SetInfo(...)` / `AppendInfoRaw(...)` などで補足情報を管理

**ポイント**
- 画面の枠組みは Frame が持つ
- 各UIは「中身だけ更新」する

### Renderer (`internal/ui/renderer.go`)
**責務**
- セクション見出し、箇条書き、ステップ、ツリー表示などの低レベル描画
- Frame から呼び出される

### Prompt Models (`internal/ui/prompt.go`)
**責務**
- 入力・選択の状態遷移とバリデーション
- `View()` は Frame を使って Inputs/Info を更新するだけにする

## 実装ルール
- 直接 `fmt.Fprintf/Printf/Println` でUI出力しない（Renderer/Frame経由）
- プロンプトの行は `Inputs` に集約し、情報は `Info` に集約する
- 独自ヘッダ（例: “Selected”）は作らず、Info にまとめる
- 対話は AltScreen を使わず、画面を奪わない（`tea.WithAltScreen` を使わない）

## 既存フローへの適用
- 連続する複数の Prompt 呼び出しは、可能な限り単一の Frame で更新する
- 既存のリスト描画は `AppendInputsRaw(...)` を使い、Frame の順序を保つ
