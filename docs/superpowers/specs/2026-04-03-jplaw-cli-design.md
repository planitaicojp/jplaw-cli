# jplaw-cli 設計書

## 概要

e-Gov法令API Version 2 を操作するためのGoベースCLIツール。
日本の法令（法律・政令・省令等）の検索・取得・閲覧を行う。

- **対象API**: https://laws.e-gov.go.jp/api/2 (認証不要)
- **参照プロジェクト**: conoha-cli (同一ディレクトリ構成パターン)
- **Goモジュール**: `github.com/planitaicojp/jplaw-cli`

## 決定事項

| 項目 | 決定 |
|------|------|
| APIクライアント生成方式 | 手動実装（oapi-codegen不使用） |
| API バージョン | v2のみ |
| サブコマンド構成 | ユーザー視点グルーピング（5コマンド） |
| デフォルト出力 | 本文→テキスト変換、一覧→テーブル |
| 法令種別・時代名の言語 | 日本語（法律、政令、令和 等） |
| 設定ファイル | `~/.config/jplaw/config.yaml` + 環境変数（環境変数優先） |
| 依存ライブラリ | cobra + yaml.v3（最小構成） |

## ディレクトリ構造

```
jplaw-cli/
├── main.go                          # cmd.Execute() 呼び出しのみ
├── go.mod                           # github.com/planitaicojp/jplaw-cli
├── Makefile                         # build, test, lint, install, clean
├── .goreleaser.yaml                 # マルチプラットフォームリリース
├── .golangci.yml                    # リンター設定
├── CLAUDE.md
├── README.md                        # 日本語
├── cmd/
│   ├── root.go                      # ルートコマンド、グローバルフラグ、Execute()
│   ├── version.go                   # jplaw version
│   ├── completion.go                # jplaw completion bash/zsh/fish
│   ├── cmdutil/
│   │   ├── client.go                # NewClient() — api.Client生成
│   │   └── format.go                # GetFormat() — フォーマット決定ロジック
│   ├── search/                      # jplaw search <keyword>
│   │   └── search.go
│   ├── list/                        # jplaw list [filters]
│   │   └── list.go
│   ├── get/                         # jplaw get <id> [--asof, --elm]
│   │   └── get.go
│   ├── history/                     # jplaw history <id>
│   │   └── history.go
│   └── attachment/                  # jplaw attachment <revision_id>
│       └── attachment.go
├── internal/
│   ├── api/
│   │   ├── client.go                # ベースHTTPクライアント (Get, Do, retry)
│   │   ├── laws.go                  # GET /laws
│   │   ├── keyword.go               # GET /keyword
│   │   ├── lawdata.go               # GET /law_data + /law_file
│   │   ├── revisions.go             # GET /law_revisions
│   │   └── attachment.go            # GET /attachment
│   ├── model/
│   │   ├── law.go                   # LawInfo, LawsResponse
│   │   ├── revision.go              # RevisionInfo, RevisionsResponse
│   │   ├── keyword.go               # KeywordResponse, Sentence
│   │   ├── lawdata.go               # LawDataResponse
│   │   ├── enums.go                 # LawType, Era, CategoryCd 等
│   │   └── error.go                 # ErrorInfo（APIエラーレスポンス）
│   ├── config/
│   │   ├── config.go                # Config struct, Load(), Save()
│   │   └── env.go                   # 環境変数処理 (JPLAW_*)
│   ├── output/
│   │   ├── formatter.go             # Formatter インターフェース + New()
│   │   ├── table.go                 # テーブルフォーマット
│   │   ├── json.go                  # JSONフォーマット
│   │   └── text.go                  # 法令本文用テキストフォーマット
│   ├── errors/
│   │   ├── errors.go                # APIError, ValidationError 等
│   │   └── exitcodes.go             # 終了コードマッピング
│   └── lawtext/
│       ├── converter.go             # JSON/XML → 可読テキスト変換
│       └── structure.go             # 編・章・節・条・項・号 構造解析
└── test/
    └── fixtures/                    # テスト用APIレスポンスサンプル
```

## サブコマンド設計

### jplaw search <keyword>

キーワードで法令本文を全文検索。`GET /keyword` にマッピング。

```
jplaw search "個人情報"
jplaw search "労働基準" --law-type 法律 --limit 20
jplaw search "契約" --asof 2025-04-01 --era 令和
```

| フラグ | 説明 | デフォルト |
|--------|------|------------|
| `--law-type` | 法令種別（法律、政令、省令等、複数可） | 全て |
| `--era` | 時代（明治〜令和） | - |
| `--asof` | 時点指定（YYYY-MM-DD） | - |
| `--category` | 分類コード | - |
| `--limit` | 結果数 | 100 |
| `--offset` | 開始位置 | 0 |
| `--format` | 出力フォーマット（table/json） | table |

デフォルト出力: 法令番号、法令名、マッチ文章（ハイライト）のテーブル。

### jplaw list

法令一覧取得。`GET /laws` にマッピング。

```
jplaw list --law-type 法律
jplaw list --law-type 政令 --era 令和
jplaw list --title "個人情報"
jplaw list --promulgation-from 2024-01-01
```

| フラグ | 説明 | デフォルト |
|--------|------|------------|
| `--law-type` | 法令種別（複数可） | 全て |
| `--era` | 時代 | - |
| `--title` | 法令名部分一致 | - |
| `--promulgation-from/to` | 公布日範囲 | - |
| `--category` | 分類コード | - |
| `--repeal-status` | 廃止状態 | - |
| `--limit` | 結果数 | 100 |
| `--offset` | 開始位置 | 0 |
| `--format` | 出力フォーマット（table/json） | table |

デフォルト出力: 法令ID、法令番号、法令名、種別、公布日のテーブル。

### jplaw get <法令IDまたは法令番号>

法令本文取得。`GET /law_data` + `GET /law_file` にマッピング。

```
jplaw get 405AC0000000088
jplaw get "平成十五年法律第五十七号"
jplaw get 405AC0000000088 --elm 第一条
jplaw get 405AC0000000088 --asof 2020-04-01
jplaw get 405AC0000000088 --file html
```

| フラグ | 説明 | デフォルト |
|--------|------|------------|
| `--asof` | 時点指定 | 最新 |
| `--elm` | 特定条文指定（例: 第一条, MainProvision） | 全体 |
| `--format` | 出力フォーマット（text/json/xml） | text |
| `--file` | ファイルダウンロード（xml/json/html/rtf/docx） | - |
| `--output` / `-o` | ファイル保存先（`--file`使用時） | stdout |

デフォルト出力: 法令本文をlawtextで変換した可読テキスト（インデントで構造表現）。

### jplaw history <法令IDまたは法令番号>

改正履歴取得。`GET /law_revisions` にマッピング。

```
jplaw history 405AC0000000088
jplaw history "平成十五年法律第五十七号" --format json
```

| フラグ | 説明 | デフォルト |
|--------|------|------------|
| `--amendment-from/to` | 改正施行日範囲 | - |
| `--format` | 出力フォーマット（table/json） | table |

デフォルト出力: 改正日、改正法令名、改正法令番号、改正種別のテーブル。

### jplaw attachment <リビジョンID>

添付ファイルダウンロード。`GET /attachment` にマッピング。

```
jplaw attachment 505AC0000000088_20250401_000000000000000
jplaw attachment 505AC0000000088_20250401_000000000000000 --src "別表第一.jpg"
jplaw attachment 505AC0000000088_20250401_000000000000000 -o ./attachments/
```

| フラグ | 説明 | デフォルト |
|--------|------|------------|
| `--src` | 特定添付ファイル指定 | 全体（ZIP） |
| `--output` / `-o` | 保存先 | カレントディレクトリ |

## グローバルフラグ

| フラグ | 環境変数 | 説明 | デフォルト |
|--------|----------|------|------------|
| `--format` | `JPLAW_FORMAT` | 出力フォーマット | table (list/search/history), text (get) |
| `--no-color` | `JPLAW_NO_COLOR` | 色無効化 | false |
| `--verbose` | `JPLAW_VERBOSE` | 詳細出力（HTTPデバッグ） | false |
| `--no-input` | `JPLAW_NO_INPUT` | プロンプト無効化 | false |

優先順位: フラグ > 環境変数 > config.yaml > デフォルト値

## データフロー

```
ユーザー入力（CLIフラグ/引数）
    │
    ▼
cmd/search/search.go          ← cobra.Command、フラグパース
    │
    ├─ cmdutil.NewClient()     ← config読込 → 環境変数オーバーライド → api.Client生成
    │
    ▼
internal/api/keyword.go        ← HTTP GET、クエリパラメータ組立、レスポンスパース
    │
    ▼
internal/model/keyword.go      ← KeywordResponse structにアンマーシャル
    │
    ▼
cmd/search/search.go          ← モデルをoutput.Formatterに渡す
    │
    ├─ text/table → stdout
    ├─ json       → stdout
    └─ エラー     → stderr + 終了コード
```

## APIクライアント

### ベースクライアント (`internal/api/client.go`)

```go
type Client struct {
    HTTPClient *http.Client
    BaseURL    string        // "https://laws.e-gov.go.jp/api/2"
    UserAgent  string        // "jplaw-cli/0.1.0"
    Debug      bool
}
```

- 認証ヘッダーなし（APIキー不要）
- `Accept: application/json` デフォルト
- リトライ: 429, 5xxに対して3回、指数バックオフ
- タイムアウト: 30秒デフォルト
- デバッグモード（`--verbose`）: リクエスト/レスポンスをstderrにログ出力

### サービス別API

| ファイル | エンドポイント | メソッド |
|----------|--------------|----------|
| `laws.go` | `GET /laws` | `ListLaws(params) (LawsResponse, error)` |
| `keyword.go` | `GET /keyword` | `SearchKeyword(params) (KeywordResponse, error)` |
| `lawdata.go` | `GET /law_data/{id}` | `GetLawData(id, params) (LawDataResponse, error)` |
| `lawdata.go` | `GET /law_file/{type}/{id}` | `GetLawFile(fileType, id, params) ([]byte, error)` |
| `revisions.go` | `GET /law_revisions/{id}` | `GetRevisions(id, params) (RevisionsResponse, error)` |
| `attachment.go` | `GET /attachment/{id}` | `GetAttachment(id, params) ([]byte, error)` |

## 設定管理

### config.yaml

```yaml
# ~/.config/jplaw/config.yaml
format: table
base_url: https://laws.e-gov.go.jp/api/2
```

### 環境変数

| 変数 | 説明 |
|------|------|
| `JPLAW_FORMAT` | デフォルト出力フォーマット |
| `JPLAW_BASE_URL` | APIベースURL（テスト/開発用） |
| `JPLAW_NO_COLOR` | 色無効化 |
| `JPLAW_VERBOSE` | デバッグ出力 |
| `JPLAW_NO_INPUT` | プロンプト無効化 |
| `JPLAW_CONFIG_DIR` | 設定ディレクトリ（デフォルト: `~/.config/jplaw`） |

### 優先順位

```
1. コマンドフラグ          --format json
2. 環境変数              JPLAW_FORMAT=json
3. config.yaml           format: json
4. コマンド別デフォルト    get→text, list/search/history→table
```

## エラー処理

### エラータイプと終了コード

| 終了コード | エラータイプ | 状況 |
|-----------|-------------|------|
| 0 | - | 成功 |
| 1 | `GeneralError` | 一般エラー |
| 2 | `ValidationError` | 不正なフラグ/引数 |
| 3 | `NotFoundError` | 法令ID/番号該当なし |
| 4 | `APIError` | API 4xx/5xx応答 |
| 5 | `NetworkError` | ネットワーク接続失敗 |

### stdout / stderr 分離

| 出力 | 対象 |
|------|------|
| 法令データ、検索結果、JSON | stdout |
| エラーメッセージ、警告、デバッグログ | stderr |

パイプ対応: `jplaw search "労働" --format json | jq '.items[].law_info.law_title'`

## lawtext変換

### 法令構造階層

```
編（Part）
  章（Chapter）
    節（Section）
      款（Subsection）
        条（Article）         ← 最も重要な単位
          項（Paragraph）
            号（Item）
              イロハ（Sub-item）
```

### 出力例

```
個人情報の保護に関する法律
（平成十五年法律第五十七号）

  第一章　総則

第一条（目的）
  この法律は、デジタル社会の進展に伴い個人情報の利用が著しく拡大して
  いることに鑑み、個人情報の適正な取扱いに関し、基本理念及び政府による
  基本方針の作成その他の個人情報の保護に関する施策の基本となる事項を
  定め、…
    一　個人情報の適正な取扱いに関し…
    二　個人情報の有用性に配慮しつつ…
```

### 変換ルール

- 法令名 + 法令番号をヘッダーとして表示
- 編/章/節 → インデントなし、タイトル行
- 条 → `第N条（タイトル）` 形式
- 項 → 2スペースインデント
- 号 → 4スペースインデント + `一`, `二`, `三`...
- 附則 → 区切り線の後に表示
- 別表 → 簡略表示 + 「添付ファイル参照」案内

### 実装方針

- `converter.go`: JSON full format → テキスト変換メインロジック
- `structure.go`: 法令要素タイプ定義、再帰的ツリー走査
- JSON full formatの `tag`/`attr`/`children` 構造を再帰パース

## テスト戦略

| レイヤー | 方式 | 内容 |
|----------|------|------|
| `internal/api/` | `httptest.NewServer` モック | HTTPリクエスト組立、レスポンスパース、リトライ、エラー処理 |
| `internal/model/` | ユニットテスト | JSONアンマーシャル、enum変換、日本語ラベル |
| `internal/lawtext/` | ユニットテスト + fixtures | 変換正確性、多様な法令構造カバー |
| `internal/config/` | ユニットテスト | 設定読込、環境変数優先順位 |
| `internal/output/` | ユニットテスト | フォーマッター別出力検証 |
| `internal/errors/` | ユニットテスト | エラータイプ → 終了コードマッピング |

テストフィクスチャ (`test/fixtures/`): 実際のAPIレスポンスを保存したJSON/XMLファイル。

## ビルド

### Makefile

```makefile
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-s -w -X github.com/planitaicojp/jplaw-cli/cmd.version=$(VERSION)"

build:      go build $(LDFLAGS) -o jplaw .
install:    go install $(LDFLAGS) .
test:       go test ./... -v
lint:       golangci-lint run ./...
coverage:   go test ./... -coverprofile=coverage.out
clean:      rm -f jplaw && rm -rf dist/
```

### GoReleaser

- マルチプラットフォーム: Linux, macOS, Windows
- マルチアーキテクチャ: amd64, arm64
- バージョン注入: ldflags経由

## 依存ライブラリ

```
github.com/spf13/cobra v1.10.2
gopkg.in/yaml.v3 v3.0.1
```

最小依存で標準ライブラリを最大限活用。
