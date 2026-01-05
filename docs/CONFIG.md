# gws 設定ファイル仕様（MVP）

場所:
- `$GWS_ROOT/settings.yaml`

## 例

```yaml
version: 1

defaults:
  # 空なら origin/HEAD を参照して自動検出する
  base_ref: ""
  ttl_days: 30

naming:
  workspace_id_must_be_valid_refname: true
  branch_equals_workspace_id: true

repo:
  # 省略形入力（github.com/org/repo）の解決用。MVPでは "github.com" 固定でも可
  default_host: "github.com"
  default_protocol: "https"  # "https" or "ssh"
```

## 仕様

- defaults.base_ref は新規ブランチ作成時の基点 
- ttl_days は gc の既定（--older が優先） 
- workspace_id_must_be_valid_refname=true の場合、無効 ID はエラー
