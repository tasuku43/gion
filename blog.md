AIエージェントで並行開発を回すようになって、Git worktree に辿り着きました。  
ただ、並行が増えるほど worktree も増えて、「どこに作ろう？」「移動がめんどい...」「これ削除していいんだっけ？」が増えてきました。  
そこで、作る・移動・片付けをいい感じにまとめたくて gion を作りました。

gion は Git worktree を「タスク（workspace）単位」で扱う小さな CLI です。  
`gion.yaml` に望ましい状態を書き、`gion apply` で差分（Plan）を確認しつつ作業場所を揃えます。  
このGIFは、YAML直編集で入った作成・削除・更新を `gion apply`（内部で Plan 表示→確認→Apply）で反映する例です。

![作成・削除・更新](https://storage.googleapis.com/zenn-user-upload/64d7ae3ea0a3-20260131.gif)

## 概要

コア機能は“作る/移動/片付け”の3つです。

- 作る：`gion.yaml` を手動（直）編集、または `gion manifest add`（対話式）で更新 → `gion apply`（Plan→確認→Apply）
- 移動：`giongo` で検索して移動（状態は変えない）
- 片付け：`gion.yaml` を手動（直）編集、または `gion manifest gc`（自動）/ `gion manifest rm`（対話式）で更新 → `gion apply`（Plan→確認→Apply）

:::message
※ `gion manifest` は `gion m` / `gion man` と短縮できます！
:::

https://github.com/tasuku43/gion

## 仕組み（gion.yaml と manifestサブコマンドの関係）

gion の中心は `gion.yaml` です。ここに「こうなっていてほしい（望ましい状態）」を書きます。  
`gion manifest` は、その `gion.yaml` を更新するための入口です（直接編集してもOKです）。

ポイントは一貫していて、**やることは常に `gion.yaml` の更新 → `gion apply`（Plan→確認→Apply）** です。  
`gion manifest add` / `gc` / `rm` はデフォルトで `apply` まで実行します（`--no-apply` を付けると `gion.yaml` の更新だけにできます）。

用語だけ補足すると、**Git worktree** はブランチ（や特定コミット）をチェックアウトした作業用ディレクトリです。一方、ここで言う **workspace** は「タスク単位の箱」で、その中に1つ以上のworktree（必要なら複数リポジトリのworktree）を束ねて扱います。

イメージとしては、だいたい次のようなディレクトリ構造になります。

```text
GION_ROOT/
├─ gion.yaml           # 望ましい状態（inventory）
├─ bare/               # 共有のbare repoストア
└─ workspaces/         # タスク単位のworkspace
   ├─ PROJ-123/        # workspace_id（タスク）
   │  ├─ backend/      # worktree（repo: backend）
   │  ├─ frontend/
   │  └─ docs/
   └─ PROJ-456/
      └─ backend/
```

---

## 作る（ApplyでPlanを確認して、まとめて作る）

workspaceを「作る」操作は、`gion manifest add` コマンドか、`gion.yaml` の直接編集で行います。  
どちらもやることは同じで、**`gion.yaml` を更新してから `gion apply`（Plan→確認→Apply）** します。  
`gion manifest add` はこの流れを1コマンドでやってくれる入口（デフォルトで `apply` まで）で、`gion.yaml` を直接編集した場合は編集後に `gion apply` を実行します。

### 4つの作成`mode`

入口は `repo` / `issue` / `review` / `preset` の4つです。  
始め方に合わせて入口を選べるだけで、行き着く先は同じで、最終的には `gion.yaml` に「こうしたい」を積んでいきます。

![入口の選択（repo/issue/review/preset）](https://storage.googleapis.com/zenn-user-upload/f59efe84c584-20260131.png)

### issue / review（まとめて積んで、一括で作る）

Issue（やPR）を複数選んで `gion.yaml` に積み、`gion apply` を実行して `Plan` で差分を確認してから反映します。

![issue/reviewをまとめて選んで、一括で作る](https://storage.googleapis.com/zenn-user-upload/027b8d9c6ecf-20260131.gif)

※ `--issue` / `--review` を使う場合は `gh` CLI が必要です（GitHub前提）。

### repo（workspaceを一つ作る）

とにかく最短で1つ作るなら `repo` が一番シンプルです。リポジトリとworkspace IDを指定して追加し、`gion apply` の `Plan` で作成内容を確認してから反映します。

![repoを1つ追加して、Planで確認する](https://storage.googleapis.com/zenn-user-upload/36fcce70fba4-20260131.png)

### preset（複数repoをworkspaceに束ねる）

workspaceは「タスク単位の箱」なので、backend + frontend + docs みたいに複数repoを束ねたくなります。presetを作っておけば、次からはそれらをまとめて一つのworkspaceを作成できます。

![presetを作成](https://storage.googleapis.com/zenn-user-upload/e715690715a9-20260131.png)
![presetで複数repoをまとめて宣言する](https://storage.googleapis.com/zenn-user-upload/c86453de43b2-20260131.png)

### YAML直編集

`gion.yaml` は直接編集も可能です。
たとえば ブランチ名を直したいとき、複数workspaceを同時に削除・作成したいとき、既存の定義を更新しつつ整理したいとき、などです。

直編集のあとに `gion apply` を実行すると、まず削除・作成・更新がまとめて `Plan` に出るので「何が起きるか」を落ち着いて確認できます。納得できたらそのまま `Apply` で反映できます。

![削除・作成・更新](https://storage.googleapis.com/zenn-user-upload/271e0d40813c-20260131.png)

---

## 移動する（workspace/worktreeを検索して移動する）

worktreeが増えてくると、「あの作業どこでやってたっけ？」を思い出す時間が増えてきます。
移動は `giongo` を使います（brew/miseで入れると `gion` と一緒に入ってきます）。  
これは状態を一切変えず、目的地を選ぶところまでを担当します。

`giongo` は workspace と worktree をまとめて一覧し、検索で絞って選べます。
workspace を選べばその workspace ディレクトリへ、worktree を選べばそのリポジトリの作業ディレクトリへ移動できます。

![workspace/worktreeを検索して移動する（giongo）](https://storage.googleapis.com/zenn-user-upload/55a6eb8dd780-20260131.gif)

:::message
`giongo` 自体はそのまま使えますが、選んだ場所に `cd` までしたい場合は bash/zsh 側で初期設定が必要です（README参照）。
:::

---

## 削除する（gcで安全に回収して、rmは止まりながら消す）

worktreeが増えてくると、「これってもう消していいんだっけ？」と立ち止まることがあると思います。  
gionはこの片付けを、`gion manifest gc` と `gion manifest rm` の2つに分けて扱います。

### gion manifest gc（自動・保守的に回収）

`gion manifest gc` は「高い確度で安全に消せるものだけ」をまとめて候補にします。  
たとえば、デフォルトブランチにマージ済みのものは回収できる一方で、判断が難しい（未コミット/未push/状態が読めない等）ものは基本的に対象外です。作っただけでコミットが無いworkspaceも、うっかり消さないように外します。

![gcの結果（回収される/されないが一目で分かる）](https://storage.googleapis.com/zenn-user-upload/592d914cb8dd-20260131.png)

### gion manifest rm（手動・ガードレール付きで消す）

一方 `gion manifest rm` は「人間が消したいもの」を選ぶための入口です。選択はインタラクティブにできて、最後に確認してから進めます。

選択画面では workspace に軽いタグが付きます（`[dirty]` / `[unpushed]` / `[diverged]` / `[unknown]`）。

- `dirty`：未コミット変更がある（未追跡ファイルやコンフリクト含む）
- `unpushed`：upstream より手元が ahead（未pushのコミットがある）
- `diverged`：ahead かつ behind（手元と upstream が分岐している）
- `unknown`：upstream が無い / detached HEAD などで判定できない

たとえば `dirty` の場合、`Plan` 側で `risk: dirty (...)` や変更ファイルが見えるので、削除前にサッと確認できます（具体的な見え方はGIF参照）。

![rmのPlan（risk/sync）と確認プロンプトの例](https://storage.googleapis.com/zenn-user-upload/f542f7f94a3f-20260201.gif)

### YAML直編集

削除も `gion.yaml` を直接編集して行えます。  
直編集の場合も `gion apply` の `Plan` で “消えるもの” が明示され、最後に確認プロンプトが出るので、うっかり消しすぎる事故を減らせます（不安ならそこで `n` で止めればOKです）。

---

## おわりに

インストール手順（Homebrew / mise 対応）と使い方はGitHubのREADMEにまとめています。よければ覗いて、手元で一度触ってみてください！

https://github.com/tasuku43/gion
