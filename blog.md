# Git worktreeの作成・移動・片付けをラクにする「gion」を作りました

AIエージェントで並行開発を回すようになって、Git worktreeに辿り着きました。  
でも worktree だけだと「作業場所の運用ルール」までは自分で作る必要があって、手癖になるまで時間がかかりそうでした。  
さらに増えると「今どこで何してるんだっけ？」の認知負荷が上がり、「これ消していいんだっけ？」が怖くなって片付けが止まる。  
だから、複数のworktreeをタスク（workspace）単位で束ねて管理できるように、gionというツールを作ってみました。

## 概要

gion は、Git worktree を「タスク（workspace）単位」で扱うための小さなCLIです。

- 作る：`gion manifest add` → `gion plan` → `gion apply`
- 移動：`giongo` で検索して移動
- 片付け：`gion manifest gc` / `gion manifest rm`

GitHub: https://github.com/tasuku43/gion

<!-- XのPostはる -->

## 仕組み（gion.yaml と manifestサブコマンドの関係）

gion の中心は `gion.yaml` です。ここに「こうなっていてほしい（望ましい状態）」を書きます。  
`gion manifest` は、その `gion.yaml` を更新するための入口です。  
そして `gion plan` で差分を確認して、`gion apply` で実体（作業場所）を揃える、という流れになっています。

用語だけ補足すると、**Git worktree** はブランチ（や特定コミット）をチェックアウトした作業用ディレクトリです。  
一方、ここで言う **workspace** は「タスク単位の箱」で、その中に1つ以上のworktree（必要なら複数リポジトリのworktree）を束ねて扱います。

イメージとしては、だいたい次のようなディレクトリ構造になります（`GION_ROOT` 配下だけを触る前提）。

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

<!-- Plan が出て、その後に確認プロンプトが出る画面（「いきなり実行しない」を一枚で）-->
*Planの確認プロンプト（いきなり実行しない）*

---

## 作る（Planで差分を見て、Applyでまとめて作る）

workspaceを「作る」操作は、`gion manifest add` コマンド（入口 / mode）か、`gion.yaml` の直接編集で行います。  
どちらの場合も、まず “望ましい状態” を宣言して `gion plan` で差分（何が作られるか）を確認し、納得できたら `gion apply` でまとめて作る——という流れです。

### 入口（mode）は4つある

入口は `repo` / `issue` / `review` / `preset` の4つです。  
始め方に合わせて入口を選べるだけで、行き着く先は同じで、最終的には `gion.yaml` に「こうしたい」を積んでいきます。

まず入口をインタラクティブに選ぶなら、これだけでOKです。
<!-- 画像を挟む -->
*入口の選択（repo/issue/review/preset）*

### issue / review（まとめて積んで、一括で作る）

Issue（やPR）を複数選んで `gion.yaml` に積み、`gion plan` で差分を見てから、`gion apply` は1回だけ。並行開発の「机をまとめて出す」がかなりラクになります。

<!-- GIFを挟む -->
*issue/reviewをまとめて選んで、一括で作る*

※ `--issue` / `--review` を使う場合は `gh` CLI が必要です（GitHub前提）。

### repo（workspaceを一つ作る）

とにかく最短で1つ作るなら `repo` が一番シンプルです。リポジトリとworkspace IDを指定して追加し、`gion plan` で作成内容を確認してから `gion apply` します。

<!-- 画像を挟む -->
*repoを1つ追加して、Planで確認する*

### preset（複数repoをworkspaceに束ねる）

workspaceは「タスク単位の箱」なので、backend + frontend + docs みたいに複数repoを束ねたくなります。presetを作っておけば、次からは `--preset` でまとめて積めます。

<!-- 画像を挟む -->
*presetで複数repoをまとめて宣言する*

### YAML直編集 vs manifest

`gion.yaml` は直接編集してもOKです。特に、すでにあるinventoryを「まとめて整える」用途に向いています。  
たとえば ブランチ名を直したいとき、複数workspaceを同時に削除・作成したいとき、既存の定義を更新しつつ整理したいとき、などです。

直編集のあとに `gion plan` を叩くと、削除・作成・更新がまとめて一覧できるので「何が起きるか」を落ち着いて確認できます。確認できたら `gion apply` で反映、という流れ自体は `gion manifest add` と同じです。

<!-- 画像: Planで「削除・作成・更新」が同時に表示されるスクショ -->
*削除・作成・更新が同時に出るPlanの例*

---

## 移動する（workspace/worktreeを検索して移動する）

worktreeが増えてくると、「あの作業どこでやってたっけ？」を思い出す時間が地味に効いてきます。  
移動は `giongo` を使います（brew/miseで入れると `gion` と一緒に入ってきます）。  
これは状態を一切変えず、目的地を選ぶところまでを担当します。

※ `giongo` 自体はそのまま使えますが、選んだ場所に `cd` までしたい場合は bash/zsh 側で関数でラップします（README参照）。

### 検索して選んで移動する

`giongo` は workspace と worktree をまとめて一覧し、検索で絞って選べます。

<!-- GIF: workspace/worktreeを検索して移動する -->
*workspace/worktreeを検索して移動する（giongo）*

---

## 削除する（gcで安全に回収して、rmは止まりながら消す）

worktreeが増えてくると、「これってもう消していいんだっけ？」と立ち止まることがあると思います。  
gionはこの片付けを、`gion manifest gc` と `gion manifest rm` の2つに分けて扱います。

### gion manifest gc（自動・保守的に回収）

`gion manifest gc` は「高い確度で安全に消せるものだけ」をまとめて候補にします。  
たとえば、デフォルトブランチにマージ済みのものは回収できる一方で、判断が難しい（未コミット/未push/状態が読めない等）ものは基本的に対象外です。作っただけでコミットが無いworkspaceも、うっかり消さないように外します。

<!-- 画像: gion manifest gc の結果（回収候補と除外が分かる） -->
*gcの結果（回収される/されないが一目で分かる）*

### gion manifest rm（手動・ガードレール付きで消す）

一方 `gion manifest rm` は「人間が消したいもの」を選ぶための入口です。選択自体はインタラクティブにできて、実行前に `Plan` で削除が出るので、そこで落ち着いて確認してから進めます。  

<!-- 画像: gion manifest rm → Plan（risk/sync）→ 確認プロンプト -->
*rmのPlan（risk/sync）と確認プロンプトの例*

---

## おわりに

インストール手順と使い方はGitHubのREADMEにまとめています。よければ覗いて、手元で一度触ってみてください！

https://github.com/tasuku43/gion
