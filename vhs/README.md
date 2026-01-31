# vhs

このディレクトリは、ブログ/README 用のデモ動画（mp4 / gif）を `vhs` で生成するための `.tape` を置きます。

## 使い方

```sh
vhs vhs/demo-apply.tape
vhs vhs/demo.tape
vhs vhs/demo-repo.tape
```

## tapes

- `vhs/demo-apply.tape`: 事前に用意した `gion.yaml` を `gion apply` で反映（概要/仕組み向け）
- `vhs/demo.tape`: `manifest add (issue)` → apply → giongo → rm（長尺の通しデモ）
- `vhs/demo-repo.tape`: `manifest add (repo)` の入力例（短尺）
