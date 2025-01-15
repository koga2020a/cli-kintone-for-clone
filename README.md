cli-kintone
==========

cli-kintoneは、kintoneアプリのデータをエクスポート・インポートするためのコマンドラインユーティリティです。

```
ⓘ このツールは git://github.com/kintone/cli-kintone から移行されました
```

## バージョン

0.14.1

## ダウンロード

以下のバイナリがダウンロード可能です：

- Windows
- Linux
- Mac OS X

https://github.com/kintone-labs/cli-kintone/releases

## 使用方法
```text
    使用方法:
        cli-kintone [オプション]

    アプリケーションオプション:
            --import  標準入力からデータをインポートします。"-f"が指定された場合は、ファイルからインポートします
            --export  kintoneのデータを標準出力にエクスポートします
        -d=           ドメイン名（FQDNを指定）
        -a=           アプリID（デフォルト: 0）
        -u=           ログインユーザー名
        -p=           ユーザーパスワード
        -t=           APIトークン
        -g=           ゲストスペースID（デフォルト: 0）
        -o=           出力形式。'json'または'csv'を指定（デフォルト: csv）
        -e=           文字エンコーディング（デフォルト: utf-8）
                      以下のエンコーディングのみサポート（フィールドコードとデータ共に）:
                      'utf-8', 'utf-16', 'utf-16be-with-signature', 'utf-16le-with-signature', 'sjis', 'euc-jp', 'gbk', 'big5'
        -U=           Basic認証ユーザー名
        -P=           Basic認証パスワード
        -q=           クエリ文字列
        -c=           エクスポートするフィールド（カンマ区切り）。フィールドコード名を指定
        -f=           入力ファイルパス
        -b=           添付ファイルディレクトリ
        -D            挿入前にレコードを削除。オプション"-q"で削除条件を指定可能
        -l=           入力ファイル内のデータ開始位置（デフォルト: 1）
        -v, --version cli-kintoneのバージョン

    ヘルプオプション:
        -h, --help    このヘルプメッセージを表示
```

## 使用例
注意:
* Windowsデバイスを使用する場合は、cli-kintone.exeを指定してください
* 事前にcli-kintoneへのPATHを環境に合わせて設定してください

### アプリから全カラムをエクスポート
```
cli-kintone --export -a <アプリID> -d <FQDN> -t <APIトークン>
```

### 指定したカラムをShift-JISエンコードでCSVファイルにエクスポート
```
cli-kintone --export -a <アプリID> -d <FQDN> -e sjis -c "$id, name1, name2" -t <APIトークン> > <出力ファイル>
```

### ファイルからアプリにインポート
```
cli-kintone --import -a <アプリID> -d <FQDN> -e sjis -t <APIトークン> -f <入力ファイル>
```

レコード番号フィールド（$id）またはキーフィールド（*記号付きフィールドコード、例："\*mykeyfield"）を含むファイルをインポートする場合、レコードは更新および/または追加されます。

- $id（またはキーフィールド）の値が既存のレコード番号と一致する場合、そのレコードは更新されます
- $id（またはキーフィールド）の値が空の場合、新しいレコードが追加されます
- $id（またはキーフィールド）の値が既存のレコード番号と一致しない場合、インポートは停止しエラーが発生します
- ファイルに$id（またはキーフィールド）列が存在しない場合、新しいレコードのみが追加され、更新は行われません

### 添付ファイルを./mydownloadsディレクトリにエクスポート・ダウンロード
```
cli-kintone --import -a <アプリID> -d <FQDN> -t <APIトークン> -b mydownloads
```

### ./myuploadsディレクトリから添付ファイルをインポート・アップロード
> :warning: 警告
>- "-b"フラグが指定されていない場合、CSVの添付ファイルフィールドの値に関係なく、添付ファイルフィールドはスキップされkintoneで更新されません
>
>- "-b"フラグが指定され、CSVの添付ファイルフィールドの値が空（空白または""）の場合、kintoneにインポート後の添付ファイルフィールドのデータは削除されます
>   - 全ての添付ファイルを削除する条件：
>       - "-b"フラグにディレクトリパスが指定されている
>       - CSVファイルに添付ファイル列が必要
>       - CSVのディレクトリパスが空
>       - 添付ファイルは任意
>   - 一部の添付ファイルを削除し一部を更新する条件：
>       - "-b"フラグにディレクトリパスが指定されている
>       - CSVファイルに添付ファイル列が必要
>       - 更新する添付ファイルのみ必要
>
>例：添付ファイルフィールドのファイルを削除するCSVファイル
>```
>"$id","Name","Department","File"
>"1","Adam Clark","Planning",""
>"2","Sarah Jones","HR",""
>```
>&nbsp;

```
cli-kintone --import -a <アプリID> -d <FQDN> -t <APIトークン> -b myuploads -f <入力ファイル>
```

### キーを指定して一括更新するインポート
> :warning: 警告
>
>CSVインポートファイルに"$id"とキーフィールドの両方が指定されている場合、「"$id"フィールドと更新キーフィールドを同時に指定することはできません」というエラーメッセージが表示されます。

一括更新のキーは、入力ファイル内でフィールドコード名の前に*を付けて指定する必要があります。
例：「"update_date","*id","status"」

```
cli-kintone --import -a <アプリID> -d <FQDN> -e sjis -t <APIトークン> -f <入力ファイル>
```

### 入力ファイルの25行目からCSVをインポート
```
cli-kintone --import -a <アプリID> -d <FQDN> -t <APIトークン> -f <入力ファイル> -l 25
```

### 標準入力（stdin）からインポート
```
printf "name,age\nJohn,37\nJane,29" | cli-kintone --import -a <アプリID> -d <FQDN> -t <APIトークン>
```

## 制限事項
* 添付ファイルフィールドへのアップロードファイルサイズは1ファイルあたり10MBまでです
* cli-kintoneではクライアント証明書を使用できません
* 以下のレコードデータは取得できません：グループフィールド、スペース、ラベル、罫線、関連レコード、ステータス、担当者、カテゴリー

## エンコード/デコードの制限
* Windowsコマンドプロンプトでは「文字化け」のような文字が正しく表示されない場合があります。
  これは中国語および日本語の文字とWindowsコマンドプロンプトの互換性の問題によるものです。
  * 中国語（繁体字/簡体字）：gbkまたはbig5エンコードでエクスポートしても正しく表示されません
  * 日本語：sjisまたはeuc-jpエンコードでエクスポートしても正しく表示されません

  この場合、以下のようにutf-8エンコードを指定してデータを表示してください：
  ```
  cli-kintone --export -a <アプリID> -d <FQDN> -e utf-8
  ```
  *この問題はWindowsコマンドプロンプトでの表示時のみ発生します。他の方法でのデータのインポート/エクスポートは、gbk、big5、sjis、euc-jpエンコードで正常に動作します。

## 基本的な使用方法のドキュメント
- 英語: https://kintone.dev/en/tutorials/tool-guides/features-of-the-command-line-tool/
- 日本語: https://developer.cybozu.io/hc/ja/articles/202957070

## ビルド方法

必要条件

- Go 1.15.15
- パッケージのクローンに必要なGitとMercurial

[Mac OS X/Linux](./docs/BuildForMacLinux.md)

[Windows](./docs/BuildForWindows.md)

## ライセンス

GPL v2

## 著作権

Copyright(c) Cybozu, Inc.
