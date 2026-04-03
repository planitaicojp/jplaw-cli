# jplaw-cli

e-Gov法令API Version 2 を操作するためのコマンドラインツール。日本の法令（法律・政令・省令等）の検索・取得・閲覧ができます。

## インストール

### Homebrew (macOS / Linux)

```bash
brew install planitaicojp/tap/jplaw
```

### Scoop (Windows)

```powershell
scoop bucket add planitaicojp https://github.com/planitaicojp/scoop-bucket
scoop install jplaw
```

### Go install

```bash
go install github.com/planitaicojp/jplaw-cli@latest
```

### リリースバイナリ

[GitHub Releases](https://github.com/planitaicojp/jplaw-cli/releases) からダウンロード:

```bash
# Linux (amd64)
curl -Lo jplaw https://github.com/planitaicojp/jplaw-cli/releases/latest/download/jplaw-cli_Linux_amd64.tar.gz
tar xzf jplaw-cli_Linux_amd64.tar.gz
chmod +x jplaw
sudo mv jplaw /usr/local/bin/

# macOS (Apple Silicon)
curl -Lo jplaw https://github.com/planitaicojp/jplaw-cli/releases/latest/download/jplaw-cli_Darwin_arm64.tar.gz
tar xzf jplaw-cli_Darwin_arm64.tar.gz
chmod +x jplaw
sudo mv jplaw /usr/local/bin/
```

## 使い方

### 法令検索

キーワードで法令本文を全文検索:

```bash
jplaw search "個人情報"
jplaw search "労働基準" --law-type 法律 --limit 20
jplaw search "契約" --asof 2025-04-01 --era 令和
```

### 法令一覧

条件を指定して法令一覧を取得:

```bash
jplaw list --law-type 法律
jplaw list --law-type 政令 --era 令和
jplaw list --title "個人情報"
jplaw list --promulgation-from 2024-01-01
```

### 法令本文取得

法令IDまたは法令番号で本文を取得:

```bash
jplaw get 405AC0000000088                    # 法令IDで取得
jplaw get "平成十五年法律第五十七号"           # 法令番号で取得
jplaw get 405AC0000000088 --elm 第一条        # 特定条文のみ
jplaw get 405AC0000000088 --asof 2020-04-01   # 過去の時点で取得
jplaw get 405AC0000000088 --file html -o law.html  # HTML形式でダウンロード
```

### 改正履歴

法令の改正履歴を取得:

```bash
jplaw history 405AC0000000088
jplaw history "平成十五年法律第五十七号" --format json
```

### 添付ファイル

法令の添付ファイル（別表等）をダウンロード:

```bash
jplaw attachment <リビジョンID>
jplaw attachment <リビジョンID> --src "別表第一.jpg" -o ./attachments/
```

## 出力フォーマット

`--format` フラグで出力形式を指定:

| フォーマット | 説明 | 対象コマンド |
|------------|------|------------|
| `table` | テーブル表示（デフォルト） | search, list, history |
| `json` | JSON出力 | 全コマンド |
| `text` | 可読テキスト（デフォルト） | get |

```bash
jplaw list --law-type 法律 --limit 5 --format json | jq '.laws[].revision_info.law_title'
```

## 設定

### 設定ファイル

`~/.config/jplaw/config.yaml`:

```yaml
format: table
base_url: https://laws.e-gov.go.jp/api/2
```

### 環境変数

| 変数 | 説明 |
|------|------|
| `JPLAW_FORMAT` | デフォルト出力フォーマット |
| `JPLAW_BASE_URL` | APIベースURL |
| `JPLAW_CONFIG_DIR` | 設定ディレクトリ |
| `JPLAW_NO_COLOR` | 色出力を無効化 |
| `JPLAW_VERBOSE` | 詳細出力 |

優先順位: コマンドフラグ > 環境変数 > config.yaml > デフォルト値

## グローバルフラグ

| フラグ | 説明 |
|--------|------|
| `--format` | 出力フォーマット (table, json, text) |
| `--verbose` | HTTPリクエスト/レスポンスのデバッグ出力 |
| `--no-color` | 色出力を無効化 |

## シェル補完

```bash
# Bash
jplaw completion bash > /etc/bash_completion.d/jplaw

# Zsh
jplaw completion zsh > "${fpath[1]}/_jplaw"

# Fish
jplaw completion fish > ~/.config/fish/completions/jplaw.fish
```

## 対象API

- **提供元**: デジタル庁 / 総務省
- **API**: [e-Gov法令API Version 2](https://laws.e-gov.go.jp/api/2)
- **仕様書**: [Swagger UI](https://laws.e-gov.go.jp/api/2/swagger-ui) / [Redoc](https://laws.e-gov.go.jp/api/2/redoc/)
- **認証**: 不要（APIキーなしで利用可能）

## ライセンス

MIT
