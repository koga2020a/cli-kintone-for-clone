# Windows環境でのビルド方法

## 前提条件
- Gitがインストールされていること
- Go言語の開発環境が整っていること

## Go言語開発環境のセットアップ

### 1. Go言語のインストール
1. [Go言語の公式サイト](https://golang.org/dl/)からWindows用インストーラー（.msi）をダウンロード
2. ダウンロードしたインストーラーを実行
3. インストール完了後、コマンドプロンプトを開いて以下のコマンドで確認：
※powershellの場合、`$env:Path += ";C:\Program Files\Go\bin"` を実行してください
```batch
go version
```

### 2. 環境変数の設定
1. Windowsの「システムのプロパティ」を開く
2. 「環境変数」をクリック
3. システム環境変数に以下を追加：
   - `GOROOT`: Goのインストールディレクトリ（通常は `C:\Program Files\Go`）
   - `Path`: 既存のPathに `%GOROOT%\bin` を追加

### 3. VS Code（推奨エディタ）のセットアップ
1. [VS Code](https://code.visualstudio.com/)をインストール
2. VS Codeを起動し、以下の拡張機能をインストール：
   - Go（ms-vscode.go）
3. VS Code上でコマンドパレット（Ctrl+Shift+P）を開き、"Go: Install/Update Tools"を実行

## ビルド手順

### 1. 開発用フォルダの作成
```batch
mkdir -p c:\tmp\dev-cli-kintone\src
```
※ パスは任意の場所に変更可能です

### 2. GOPATH環境変数の設定
```batch
set GOPATH=c:\tmp\dev-cli-kintone
```

### 3. cli-kintoneリポジトリのクローン
```batch
cd %GOPATH%\src
git clone https://github.com/kintone/cli-kintone.git
```

### 4. 依存パッケージのインストール
```batch
cd %GOPATH%\src\cli-kintone
go get github.com/mattn/gom
..\..\bin\gom.exe -production install
```

### 5. ビルドの実行
```batch
go build -o .\build\windows-x64\cli-kintone.exe
```

## 著作権表示
Copyright(c) Cybozu, Inc.
