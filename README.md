# cli-kintone

kintoneアプリのデータをエクスポート・インポートするためのコマンドラインツールです。

> **Note**  
> このリポジトリは、[オリジナルのcli-kintone](https://github.com/kintone/cli-kintone)（現在は開発停止中）をフォークし、現代のgoLang環境に対応するように更新したものです。主な改善点：
> - 最新のgoLangバージョンに対応
> - 依存パッケージを最新化
> - パフォーマンスの改善
> - バグ修正

## 📚 ドキュメント

詳細なドキュメントは以下のウェブサイトをご覧ください：
https://cli.kintone.dev/

## 🚀 インストール

npmを使用してインストール：

```bash
npm install -g @kintone/cli-kintone
```

## 📋 基本的な使用例

### レコードのエクスポート
```bash
cli-kintone --export -a <アプリID> -d <ドメイン> -t <APIトークン>
```

### レコードのインポート
```bash
cli-kintone --import -a <アプリID> -d <ドメイン> -t <APIトークン> -f <入力ファイル>
```

### 主なオプション

| オプション | 説明 | デフォルト値 | 必須 |
|------------|------|------------|------|
| `-i, --import` | 標準入力またはファイルからデータをインポート | - | - |
| `-e, --export` | kintoneからデータをエクスポート | - | - |
| `-a, --app-id` | アプリIDを指定 | - | ✅ |
| `-d, --domain` | kintoneドメインを指定（例：company.kintone.com） | - | ✅ |
| `-u, --user` | ログインユーザー名 | - | APIトークンがない場合✅ |
| `-p, --password` | ログインパスワード | - | APIトークンがない場合✅ |
| `-t, --api-token` | APIトークン | - | ユーザー認証がない場合✅ |
| `-g, --guest-space-id` | ゲストスペースID | - | - |
| `-o, --format` | 出力形式（json または csv） | csv | - |
| `-e, --encoding` | 文字エンコーディング（utf-8, utf-16, utf-16be-with-signature, utf-16le-with-signature, sjis, euc-jp, gbk, big5） | utf-8 | - |
| `-U, --basic-auth-user` | Basic認証のユーザー名 | - | - |
| `-P, --basic-auth-password` | Basic認証のパスワード | - | - |
| `-q, --query` | クエリ文字列 | - | - |
| `-c, --fields` | エクスポートするフィールド（カンマ区切り） | - | - |
| `-f, --file` | 入出力ファイルパス | - | - |
| `-b, --file-dir` | 添付ファイルのディレクトリ（上限100MB/ファイル） | - | - |
| `-D, --delete-all` | インポート前にレコードを削除 | false | - |
| `-l, --line` | 入力ファイル内のデータ位置インデックス | - | - |
| `--bulk-wait-seconds` | バルクリクエスト後の待機時間（秒） | 1 | - |
| `--bulk-wait-seconds-with-file` | ファイル添付時のバルクリクエスト後の待機時間（秒） | 30 | - |
| `--bulk-limit-record-option` | 1回のバルクリクエストで処理するレコード数の上限 | 10 | - |

### 認証方法

以下のいずれかの認証方法を選択できます：

1. APIトークン認証
```bash
cli-kintone --export -a <アプリID> -d <ドメイン> -t <APIトークン>
```

2. ユーザー認証
```bash
cli-kintone --export -a <アプリID> -d <ドメイン> -u <ユーザー名> -p <パスワード>
```

3. Basic認証付きの場合
```bash
cli-kintone --export -a <アプリID> -d <ドメイン> -t <APIトークン> \
  -U <Basic認証ユーザー> -P <Basic認証パスワード>
```

### 高度な使用例

```bash
# 特定のフィールドのみをエクスポート
cli-kintone --export -a 1234 -d company.kintone.com -t "XXX" --fields "field1,field2"

# クエリを指定してエクスポート
cli-kintone --export -a 1234 -d company.kintone.com -t "XXX" --query "created_time > \"2023-01-01\""

# Basic認証付きでインポート
cli-kintone --import -a 1234 -d company.kintone.com -t "XXX" -f data.csv \
  --basic-auth-username "user" --basic-auth-password "pass"

# JSONフォーマットでエクスポート
cli-kintone --export -a 1234 -d company.kintone.com -t "XXX" --format json
```

## 🔧 主な機能

- CSVまたはJSONフォーマットでのデータのインポート/エクスポート
- 添付ファイルの取り扱い
- バルクリクエストのサポート
- 複数の文字エンコーディングのサポート
- Basic認証対応

## ⚠️ 制限事項

- 1リクエストあたりの最大レコード数：100件
- 添付ファイルの最大サイズ：100MB/ファイル
- クライアント証明書は非対応

詳細な制限事項については[公式ドキュメント](https://cli.kintone.dev/)をご確認ください。

## 🤝 コントリビューション

バグ報告や機能要望は[Issues](https://github.com/kintone/cli-kintone/issues)にて受け付けています。

## 📜 ライセンス

MIT License

## 📝 注記

このプロジェクトは、Cybozu, Inc.による[オリジナルのcli-kintone](https://github.com/kintone/cli-kintone)を基に、現代のgolang環境でビルド可能なように改良されたバージョンです。オリジナルのプロジェクトは開発が停止していますが、このバージョンでは以下の改善が行われています：

- 大きいサイズの添付ファイルのアップロード時に、より確実なデータ処理
- バルクリクエストの制御方法の改善

## 👥 作者

- オリジナル: Copyright (c) Cybozu, Inc.
- 現バージョン: [koga2020a](https://github.com/koga2020a)

## 🔨 ビルド方法

### 必要条件
- Go 1.21以上

### ソースからビルド
```bash
# リポジトリのクローン
git clone https://github.com/koga2020a/cli-kintone.git
cd cli-kintone

# 依存パッケージのインストール
go mod download

# ビルド
go build -o cli-kintone

# インストール（オプション）
go install
```

### バイナリのダウンロード
[Releases](https://github.com/koga2020a/cli-kintone/releases)ページから、各プラットフォーム向けのビルド済みバイナリをダウンロードできます：
- Windows (x64, ARM64)

### クロスコンパイル(Windows 64bitのみ)
```bash
# Windows (64bit) 向けビルド
GOOS=windows GOARCH=amd64 go build -o cli-kintone.exe
```
