package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/kintone-labs/go-kintone"
)

// テスト用のKintoneアプリ接続情報を環境変数から設定
func newAppTest(id uint64) *kintone.App {
	return &kintone.App{
		Domain:   os.Getenv("KINTONE_DOMAIN"),
		User:     os.Getenv("KINTONE_USERNAME"),
		Password: os.Getenv("KINTONE_PASSWORD"),
		AppId:    id,
	}
}

func TestRequest(t *testing.T) {
	// バルクリクエストの初期化
	bulkReq := &BulkRequests{}
	app := newAppTest(16) // AppID: 16のテスト用アプリを作成
	bulkReq.Requests = make([]*BulkRequestItem, 0)

	// レコード追加（POST）処理
	/// INSERT
	records := make([]*kintone.Record, 0)
	record1 := kintone.NewRecord(map[string]interface{}{
		"Text": kintone.SingleLineTextField("test 11!"),
		"_2":   kintone.SingleLineTextField("test 21!"),
	})
	// ... 2つ目のレコード作成と追加 ...
	records = append(records, record1)
	record2 := kintone.NewRecord(map[string]interface{}{
		"Text": kintone.SingleLineTextField("test 22!"),
		"_2":   kintone.SingleLineTextField("test 22!"),
	})
	records = append(records, record2)
	dataPOST := &DataRequestRecordsPOST{app.AppId, records}
	postRecords := &BulkRequestItem{"POST", "/k/v1/records.json", dataPOST}

	bulkReq.Requests = append(bulkReq.Requests, postRecords)

	// レコード更新（PUT）処理
	/// UPDATE
	recordsUpdate := make([]interface{}, 0)
	// ID:4902のレコードを更新
	recordsUpdate1 := kintone.NewRecordWithId(4902, map[string]interface{}{
		"Text": kintone.SingleLineTextField("test NNN!"),
		"_2":   kintone.SingleLineTextField("test MMM!"),
	})
	// ... ID:4903のレコード更新処理 ...
	fmt.Println(recordsUpdate1)
	recordsUpdate = append(recordsUpdate, &DataRequestRecordPUT{ID: recordsUpdate1.Id(),
		Record: recordsUpdate1})

	recordsUpdate2 := kintone.NewRecordWithId(4903, map[string]interface{}{
		"Text": kintone.SingleLineTextField("test 123!"),
		"_2":   kintone.SingleLineTextField("test 234!"),
	})
	recordsUpdate = append(recordsUpdate, &DataRequestRecordPUT{
		ID: recordsUpdate2.Id(), Record: recordsUpdate2})

	dataPUT := &DataRequestRecordsPUT{app.AppId, recordsUpdate}
	putRecords := &BulkRequestItem{"PUT", "/k/v1/records.json", dataPUT}

	bulkReq.Requests = append(bulkReq.Requests, putRecords)

	// バルクリクエストの実行とエラー確認
	rs, err := bulkReq.Request(app)

	if err != nil {
		t.Error(" failed", err)
	} else {
		t.Log(rs)
	}
}
