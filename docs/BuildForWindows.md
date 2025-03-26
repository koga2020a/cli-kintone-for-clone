# Windows環境でのcli-kintoneビルドガイド

## 前提条件
- Git がインストール済みであること
- Go 1.21以上の開発環境が整っていること

## 開発環境のセットアップ

### 1. Go言語のインストール
1. [Go公式サイト](https://golang.org/dl/)からWindows用インストーラーをダウンロード
2. インストーラーを実行し、指示に従ってインストール
3. インストール確認:
   ```batch
   go version
   ```

### 2. 環境変数の設定
1. 「システムのプロパティ」→「環境変数」を開く
2. システム環境変数に以下を追加:
   - `GOROOT`: Goのインストールディレクトリ（通常は `C:\Program Files\Go`）
   - `Path`: 既存のPathに `%GOROOT%\bin` を追加

### 3. 推奨開発環境 (VS Code)
1. [VS Code](https://code.visualstudio.com/)をインストール
2. VS Codeで「Go」拡張機能(ms-vscode.go)をインストール
3. コマンドパレット(Ctrl+Shift+P)で "Go: Install/Update Tools" を実行

## ビルド手順

### A. 従来のビルド方法
1. 開発用フォルダの作成
   ```batch
   mkdir -p c:\tmp\dev-cli-kintone\src
   ```
2. GOPATH環境変数の設定
   ```batch
   set GOPATH=c:\tmp\dev-cli-kintone
   ```
3. cli-kintoneリポジトリのクローン
   ```batch
   cd %GOPATH%\src
   git clone https://github.com/kintone/cli-kintone.git
   ```
4. 依存パッケージのインストール
   ```batch
   cd %GOPATH%\src\cli-kintone
   go get github.com/mattn/gom
   ..\..\bin\gom.exe -production install
   ```
5. ビルドの実行
   ```batch
   go build -o .\build\windows-x64\cli-kintone.exe
   ```

### B. モダンなビルド方法（推奨）
- 自動ビルド・テストスクリプトを利用できます。詳細は `docs/build_windows.md` を参照してください。

## 使用方法の例
- その他のオプションやドキュメントは [公式サイト](https://github.com/kintone/cli-kintone) を参照してください。

## 著作権表示
- オリジナル: Copyright(c) Cybozu, Inc.
- 現バージョン: koga2020a
