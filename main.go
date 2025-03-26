package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/howeyc/gopass"
	"github.com/kintone-labs/go-kintone"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	flags "github.com/jessevdk/go-flags"
)

// NAME of this package
const NAME = "cli-kintone"

// VERSION of this package
const VERSION = "0.14.1.for-clone"

// IMPORT_ROW_LIMIT The maximum row will be import
const IMPORT_ROW_LIMIT = 100

// EXPORT_ROW_LIMIT The maximum row will be export
const EXPORT_ROW_LIMIT = 500

// Configure of this package
type Configure struct {
	IsImport                bool     `long:"import" description:"Import data from standard input. If '-f' is specified, import data from file"`
	IsExport                bool     `long:"export" description:"Export data from kintone"`
	Domain                  string   `short:"d" long:"domain" description:"Domain name (specify FQDN)"`
	AppID                   uint64   `short:"a" long:"app-id" description:"App ID"`
	Login                   string   `short:"u" long:"user" description:"User's login name"`
	Password                string   `short:"p" long:"password" description:"User's password"`
	APIToken                string   `short:"t" long:"api-token" description:"API token"`
	GuestSpaceID            uint64   `short:"g" long:"guest-space-id" description:"Guest space ID"`
	Format                  string   `short:"o" long:"format" description:"Output format. Specify 'json' or 'csv'"`
	Encoding                string   `short:"e" long:"encoding" description:"Character encoding. Supported encodings: 'utf-8', 'utf-16', 'utf-16be-with-signature', 'utf-16le-with-signature', 'sjis', 'euc-jp', 'gbk', 'big5'"`
	BasicAuthUser           string   `short:"U" long:"basic-auth-user" description:"Basic authentication username"`
	BasicAuthPassword       string   `short:"P" long:"basic-auth-password" description:"Basic authentication password"`
	Query                   string   `short:"q" long:"query" description:"Query string"`
	Fields                  []string `short:"c" long:"fields" description:"Fields to export (comma-separated). Specify field code names"`
	FilePath                string   `short:"f" long:"file" description:"Input/Output file path for import/export"`
	FileDir                 string   `short:"b" long:"file-dir" description:"Directory for file attachments (file attachment limit is 100MB per file)"`
	DeleteAll               bool     `short:"D" long:"delete-all" description:"Delete records before insertion. Use option '-q' to specify delete conditions"`
	Line                    uint64   `short:"l" long:"line" description:"Data position index in input file"`
	Version                 bool     `short:"v" long:"version" description:"Display cli-kintone version"`
	BulkWaitSeconds         int      `long:"bulk-wait-seconds" description:"バルクリクエスト後の待機時間（秒）。デフォルトは1秒です。" default:"1"`
	BulkWaitSecondsWithFile int      `long:"bulk-wait-seconds-with-file" description:"ファイルディレクトリ指定時のバルクリクエスト後の待機時間（秒）。デフォルトは30秒です。" default:"30"`
	BulkLimitRecordOption   int      `long:"bulk-limit-record-option" description:"1回のバルクリクエストで処理できるレコード数の上限。デフォルトは10です。" default:"10"`
	DisableAutoContinue     bool     `long:"disable-auto-continue" description:"Disable auto-continue"`
}

var config Configure

// Column config
// Column config is deprecated, replace using Cell config
type Column struct {
	Code       string
	Type       string
	IsSubField bool
	Table      string
}

// Columns config
// Columns config is deprecated, replace using Row config
type Columns []*Column

// Cell config
type Cell struct {
	Code       string
	Type       string
	IsSubField bool
	Table      string
	Index      int
}

// Row config
type Row []*Cell

func init() {
	// デフォルト値の設定（必要に応じて）
	config.BulkWaitSeconds = 1
	config.BulkWaitSecondsWithFile = 30
}

func getFields(app *kintone.App) (map[string]*kintone.FieldInfo, error) {
	fields, err := app.Fields()
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func getSupportedFields(app *kintone.App) (map[string]*kintone.FieldInfo, error) {
	fields, err := getFields(app)
	if err != nil {
		return nil, err
	}
	for key, field := range fields {
		switch field.Type {
		case "STATUS_ASSIGNEE", "CATEGORY", "STATUS":
			delete(fields, key)
		default:
			continue
		}
	}
	return fields, nil
}

// set column information from fieldinfo
// This function is deprecated, replace using function getCell
func getColumn(code string, fields map[string]*kintone.FieldInfo) *Column {
	// initialize values
	column := Column{Code: code, IsSubField: false, Table: ""}

	if code == "$id" {
		column.Type = kintone.FT_ID
		return &column
	} else if code == "$revision" {
		column.Type = kintone.FT_REVISION
		return &column
	} else {
		// is this code the one of sub field?
		for _, val := range fields {
			if val.Code == code {
				column.Type = val.Type
				return &column
			}
			if val.Type == kintone.FT_SUBTABLE {
				for _, subField := range val.Fields {
					if subField.Code == code {
						column.IsSubField = true
						column.Type = subField.Type
						column.Table = val.Code
						return &column
					}
				}
			}
		}
	}

	// the code is not found
	column.Type = "UNKNOWN"
	return &column
}

func containtString(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// set Cell information from fieldinfo
// function replace getColumn so getColumn is invalid name
func getCell(code string, fields map[string]*kintone.FieldInfo) *Cell {
	// initialize values
	cell := Cell{Code: code, IsSubField: false, Table: ""}

	if code == "$id" {
		cell.Type = kintone.FT_ID
		return &cell
	} else if code == "$revision" {
		cell.Type = kintone.FT_REVISION
		return &cell
	} else {
		// is this code the one of sub field?
		for _, val := range fields {
			if val.Code == code {
				cell.Type = val.Type
				return &cell
			}
			if val.Type == kintone.FT_SUBTABLE {
				for _, subField := range val.Fields {
					if subField.Code == code {
						cell.IsSubField = true
						cell.Type = subField.Type
						cell.Table = val.Code
						return &cell
					}
				}
			}
		}
	}

	// the code is not found
	cell.Type = "UNKNOWN"
	return &cell
}

func getEncoding() encoding.Encoding {
	switch config.Encoding {
	case "utf-8":
		return nil
	case "utf-16":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case "utf-16be-with-signature":
		return unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	case "utf-16le-with-signature":
		return unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	case "euc-jp":
		return japanese.EUCJP
	case "sjis":
		return japanese.ShiftJIS
	case "gbk":
		return simplifiedchinese.GBK
	case "big5":
		return traditionalchinese.Big5
	default:
		return nil
	}
}

// 新しい関数を追加
func getConsoleWriter() io.Writer {
	if runtime.GOOS == "windows" {
		// Windowsの場合はUTF-8で出力できるようにする
		stdout := os.Stdout
		if err := setConsoleUTF8(stdout); err == nil {
			return stdout
		}
		// UTF-8設定に失敗した場合はShift-JISに変換
		return transform.NewWriter(os.Stdout, japanese.ShiftJIS.NewEncoder())
	}
	// その他のOSではそのまま標準出力を使用
	return os.Stdout
}

// Windows環境でUTF-8出力を有効にする
func setConsoleUTF8(file *os.File) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	// Windows用のコンソールをUTF-8モードに設定
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleOutputCP")
	r, _, err := proc.Call(65001) // 65001 is the code page for UTF-8
	if r == 0 {
		return err
	}
	return nil
}

func validateConfig() error {
	if config.BulkWaitSeconds < 0 {
		return errors.New("bulk-wait-seconds は0以上の整数でなければなりません")
	}
	if config.BulkWaitSecondsWithFile < 0 {
		return errors.New("bulk-wait-seconds-with-file は0以上の整数でなければなりません")
	}
	if config.BulkLimitRecordOption <= 0 {
		return errors.New("bulk-limit-record-option は1以上の整数でなければなりません")
	}
	return nil
}

func main() {
	var err error

	parser := flags.NewParser(&config, flags.Default)
	parser.ShortDescription = "cli-kintone は kintone と対話するための CLI ツールです。"
	parser.LongDescription = "cli-kintone を使用すると、kintone アプリケーションからデータをインポートおよびエクスポートできます。"

	// 引数がない場合、ヘルプを表示して終了
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	_, err = parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			log.SetFlags(0)
			printHelp()
			os.Exit(0)
		}
		// ヘルプやバージョン以外のエラーの場合、ヒントを表示
		fileExecute := os.Args[0]
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		fmt.Printf("\n詳細については '%s --help' を試してください。\n", fileExecute)
		os.Exit(1)
	}

	if config.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if (config.IsImport || config.IsExport) && (config.AppID == 0 || (config.APIToken == "" && (config.Domain == "" || config.Login == ""))) {
		printHelp()
		os.Exit(1)
	}

	if !strings.Contains(config.Domain, ".") {
		config.Domain += ".cybozu.com"
	}

	// カンマ区切りのフィールドをサポート
	var cols []string
	if len(config.Fields) > 0 {
		for _, field := range config.Fields {
			curField := strings.Split(field, ",")
			cols = append(cols, curField...)
		}
		config.Fields = nil
		for _, col := range cols {
			curFieldString := strings.TrimSpace(col)
			if curFieldString != "" {
				config.Fields = append(config.Fields, curFieldString)
			}
		}
	}

	var app *kintone.App
	if config.BasicAuthUser != "" && config.BasicAuthPassword == "" {
		fmt.Printf("ベーシック認証のパスワード: ")
		pass, _ := gopass.GetPasswd()
		config.BasicAuthPassword = string(pass)
	}

	if config.APIToken == "" {
		if config.Password == "" {
			fmt.Printf("パスワード: ")
			pass, _ := gopass.GetPasswd()
			config.Password = string(pass)
		}

		app = &kintone.App{
			Domain:       config.Domain,
			User:         config.Login,
			Password:     config.Password,
			AppId:        config.AppID,
			GuestSpaceId: config.GuestSpaceID,
		}
	} else {
		app = &kintone.App{
			Domain:       config.Domain,
			ApiToken:     config.APIToken,
			AppId:        config.AppID,
			GuestSpaceId: config.GuestSpaceID,
		}
	}

	if config.BasicAuthUser != "" {
		app.SetBasicAuth(config.BasicAuthUser, config.BasicAuthPassword)
	}

	app.SetUserAgentHeader(NAME + "/" + VERSION + " (" + runtime.GOOS + " " + runtime.GOARCH + ")")

	// 出力先の決定
	var writer io.Writer
	if config.IsExport {
		if config.Query != "" {
			if config.FilePath != "" {
				file, err := os.Create(config.FilePath)
				if err != nil {
					log.Fatalf("出力ファイルの作成に失敗しました: %v", err)
				}
				defer file.Close()
				writer = getWriter(file)
			} else {
				writer = getWriter(getConsoleWriter())
			}
			err = exportRecordsWithQuery(app, config.Fields, writer)
		} else {
			fields := config.Fields
			isAppendIdCustome := false
			if len(config.Fields) > 0 && !containtString(config.Fields, "$id") {
				fields = append(fields, "$id")
				isAppendIdCustome = true
			}
			if config.FilePath != "" {
				file, err := os.Create(config.FilePath)
				if err != nil {
					log.Fatalf("出力ファイルの作成に失敗しました: %v", err)
				}
				defer file.Close()
				writer = getWriter(file)
			} else {
				writer = getWriter(getConsoleWriter())
			}
			err = exportRecordsBySeekMethod(app, writer, fields, isAppendIdCustome)
		}
	}

	// インポートおよびエクスポートの処理
	if config.IsImport && config.IsExport {
		log.Fatal("オプション --import と --export は同時に指定できません！")
	}

	if config.IsImport {
		if config.FilePath == "" {
			err = importFromCSV(app, os.Stdin)
		} else {
			err = importDataFromFile(app)
		}
	}

	// 設定のバリデーション
	if err := validateConfig(); err != nil {
		log.Fatalf("設定エラー: %v", err)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func importDataFromFile(app *kintone.App) error {
	var file *os.File
	var err error
	file, err = os.Open(config.FilePath)
	if err == nil {
		defer file.Close()
		err = importFromCSV(app, file)
	}
	return err
}

// 新しいヘルプ表示関数を修正
func printHelp() {
	log.SetFlags(0) // タイムスタンプなどのプレフィックスを無効化

	log.Output(2, "使用方法:")
	log.Output(2, "  cli-kintone [オプション]")
	log.Output(2, "")
	log.Output(2, "オプション:")
	log.Output(2, "  -i, --import                    標準入力からデータをインポートします。'-f' が指定されている場合はファイルからデータをインポートします")
	log.Output(2, "  -e, --export                    kintoneデータをエクスポートします")
	log.Output(2, "  -d, --domain=DOMAIN             ドメイン名（FQDNを指定）")
	log.Output(2, "  -a, --app-id=APP-ID            アプリID")
	log.Output(2, "  -u, --user=USER                 ユーザーのログイン名")
	log.Output(2, "  -p, --password=PASSWORD         ユーザーのパスワード")
	log.Output(2, "  -t, --api-token=TOKEN          APIトークン")
	log.Output(2, "  -g, --guest-space-id=ID        ゲストスペースID")
	log.Output(2, "  -o, --format=FORMAT            出力形式。'json' または 'csv' を指定")
	log.Output(2, "  -e, --encoding=ENCODING        文字エンコーディング")
	log.Output(2, "                                 サポートされているエンコーディング:")
	log.Output(2, "                                 'utf-8', 'utf-16', 'utf-16be-with-signature',")
	log.Output(2, "                                 'utf-16le-with-signature', 'sjis', 'euc-jp',")
	log.Output(2, "                                 'gbk', 'big5'")
	log.Output(2, "  -U, --basic-auth-user=USER     ベーシック認証のユーザー名")
	log.Output(2, "  -P, --basic-auth-password=PASS ベーシック認証のパスワード")
	log.Output(2, "  -q, --query=QUERY              クエリ文字列")
	log.Output(2, "  -c, --fields=FIELDS            エクスポートするフィールド（カンマ区切り）")
	log.Output(2, "  -f, --file=FILE                csvファイル 入力・出力ファイルパス")
	log.Output(2, "  -b, --file-dir=DIR             添付ファイルのディレクトリ（各添付ファイルの上限は100MBです）")
	log.Output(2, "  -D, --delete-all               挿入前にレコードを削除します")
	log.Output(2, "  -l, --line=LINE                入力ファイル内のデータの位置インデックス")
	log.Output(2, "  -v, --version                  cli-kintone のバージョンを表示します")
	log.Output(2, "  -h, --help                     このヘルプを表示します")
	log.Output(2, "  --bulk-wait-seconds=SECONDS             バルクリクエスト後の待機時間（秒）。デフォルトは1秒です。")
	log.Output(2, "  --bulk-wait-seconds-with-file=SECONDS   ファイルディレクトリ指定時のバルクリクエスト後の待機時間（秒）。デフォルトは30秒です。")
	log.Output(2, "  --bulk-limit-record-option=OPTION        1回のバルクリクエストで処理できるレコード数の上限。デフォルトは10です。")
	log.Output(2, "  --disable-auto-continue                 エラー発生時に次のレコードの処理を継続しない（デフォルトでは継続する）")

	log.Output(2, "詳細については https://github.com/kintone/cli-kintone を参照してください")
}
