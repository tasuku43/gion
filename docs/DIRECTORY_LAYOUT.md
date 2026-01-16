# ディレクトリ設計（確定）

## ルート
- `GWS_ROOT` があればそれを使う
- 無ければ `~/gws`
- `templates.yaml` は `GWS_ROOT` 直下に配置

## 配下構造（固定）
- `$GWS_ROOT/bare/`   : repo store（bare repo）
- `$GWS_ROOT/workspaces/`     : workspace (AI)
- `$GWS_ROOT/workspaces/<ID>/`:
    - `<alias>/`        : worktree 作業ディレクトリ
    - `.gws/metadata.json` : workspace のメタデータ（description など）

## repo store のパス
- `$GWS_ROOT/bare/<host>/<owner>/<repo>.git`
