```markdown:build_and_test.ps1
# 自動ビルド・試験スクリプトの作成・実行ガイド

このドキュメントでは、Windows環境における自動ビルドおよび試験のスクリプト作成と実行方法について説明します。

## 前提条件

- **PowerShell** がインストールされていること。
- **Go** がインストールされ、`go` コマンドが使用可能であること。
- 必要なパスや認証情報が適切に設定されていること。

## スクリプトの作成

以下の手順でビルドおよび試験の自動化スクリプトを作成します。

1. 任意のディレクトリに `build_and_test.ps1` という名前でファイルを作成します。
2. 以下の内容をファイルに貼り付けて保存します。

```powershell:build_and_test.ps1
# トランスクリプトの開始
$logFile = ".\build_test_log.txt"
Start-Transcript -Path $logFile

Write-Host "ビルドを開始します..."

# ビルドの実行
$buildOutput = & go build -o .\build\windows-x64\cli-kintone.exe 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "ビルドが成功しました。"
} else {
    Write-Host "ビルドが失敗しました。"
    Write-Host $buildOutput
    exit 1
}

Write-Host "試験を開始します..."

# 試験の実行
$testOutput = & .\build\windows-x64\cli-kintone.exe --export -d=6rzyuy6y7jv1.cybozu.com -a=13 -t=KE8AYxx3PJ5Ln3UX5Ai0YkQBNC5g5Ci0iUSaIo4j -b=build_test_downloadRecord\13_新しいアプリ\attachmentFiles -e=utf-8
if ($LASTEXITCODE -eq 0) {
    Write-Host "試験が成功しました。"
} else {
    Write-Host "試験が失敗しました。"
    Write-Host $testOutput
    exit 1
}

Write-Host "すべての処理が完了しました。"
Stop-Transcript

# 一時ファイルとディレクトリの削除
if (Test-Path $logFile) {
    Remove-Item $logFile -Force
}

if (Test-Path ".\build_test_downloadRecord") {
    Remove-Item ".\build_test_downloadRecord" -Recurse -Force
}
```

## スクリプトの内容説明

### ログの記録

```powershell:build_and_test.ps1
# トランスクリプトの開始
$logFile = ".\build_test_log.txt"
Start-Transcript -Path $logFile
```

- ビルドおよび試験のログを `build_test_log.txt` に記録します。

### ビルドの実行

```powershell:build_and_test.ps1
Write-Host "ビルドを開始します..."

# ビルドの実行
$buildOutput = & go build -o .\build\windows-x64\cli-kintone.exe 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "ビルドが成功しました。"
} else {
    Write-Host "ビルドが失敗しました。"
    Write-Host $buildOutput
    exit 1
}
```

- `go build` コマンドを実行してプロジェクトをビルドします。
- ビルドが成功した場合は成功メッセージを表示し、失敗した場合はエラーメッセージを表示してスクリプトを終了します。

### 試験の実行

```powershell:build_and_test.ps1
Write-Host "試験を開始します..."

# 試験の実行
$testOutput = & .\build\windows-x64\cli-kintone.exe --export -d=6rzyuy6y7jv1.cybozu.com -a=13 -t=**************************************** -b=build_test_downloadRecord\13_新しいアプリ\attachmentFiles -e=utf-8
if ($LASTEXITCODE -eq 0) {
    Write-Host "試験が成功しました。"
} else {
    Write-Host "試験が失敗しました。"
    Write-Host $testOutput
    exit 1
}

Write-Host "すべての処理が完了しました。"
Stop-Transcript
```

- ビルドされた実行ファイルを使用して試験を実行します。
- 試験の結果に応じて成功または失敗のメッセージを表示します。

### 一時ファイルとディレクトリの削除

```powershell:build_and_test.ps1
# 一時ファイルとディレクトリの削除
if (Test-Path $logFile) {
    Remove-Item $logFile -Force
}

if (Test-Path ".\build_test_downloadRecord") {
    Remove-Item ".\build_test_downloadRecord" -Recurse -Force
}
```

- ログファイルおよび試験で使用した一時ディレクトリを削除します。

## スクリプトの実行方法

1. **PowerShellを管理者として実行**します。
2. スクリプトが保存されているディレクトリに移動します。

   ```powershell
   cd パス\to\スクリプトディレクトリ
   ```

3. スクリプトの実行ポリシーを設定します（必要に応じて）。

   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

4. スクリプトを実行します。

   ```powershell
   .\build_and_test.ps1
   ```

## 注意事項

- スクリプト内のパラメータ（例：ドメイン、アプリID、認証トークンなど）は環境に合わせて適切に設定してください。
- ビルドおよび試験に必要なファイルやディレクトリが存在することを確認してください。
- スクリプト実行後、必要に応じてログファイルや一時ファイルのバックアップを保存してください。

## トラブルシューティング

- **ビルドが失敗する場合**:
  - Goの環境設定が正しいか確認してください。
  - ソースコードにエラーがないか確認してください。

- **試験が失敗する場合**:
  - 接続先のドメインや認証情報が正しいか確認してください。
  - 必要なリソースが適切に配置されているか確認してください。

- **スクリプトが実行できない場合**:
  - 実行ポリシーが適切に設定されているか確認してください。
  - スクリプトファイルのパスが正しいか確認してください。

以上で、自動ビルドおよび試験スクリプトの作成と実行に関するガイドを終わります。必要に応じてスクリプトをカスタマイズし、プロジェクトのニーズに合わせて最適化してください。

```
