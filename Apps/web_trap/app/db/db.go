package db

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

// InitDB is method to initialize DB
func InitDB() {
	log.Println("========== DB 初期化 開始 ==========")

	// DB接続
	log.Println("== DB 接続 ==")
	conn := getConnectionString()
	db, err := sql.Open("mysql", conn)
	if err != nil {
		outputErrorLog("DB 接続 エラー", err)
	}
	log.Println("== DB 接続 成功 ==")

	// 既存のテーブル削除
	log.Println("== DB 削除 ==")
	db.Exec("DROP TABLE users")
	log.Println("== DB 削除 成功 ==")

	// テーブル生成
	log.Println("== DB テーブル 作成 ==")
	_, err = db.Exec(`CREATE TABLE users ( 
		id INTEGER AUTO_INCREMENT PRIMARY KEY, 
		account VARCHAR(32) NOT NULL, 
		name VARCHAR(32) NOT NULL, 
		password VARCHAR(32) NOT NULL
		)`)
	if err != nil {
		outputErrorLog("usersテーブル作成失敗", err)
	}
	log.Println("== DB テーブル 作成 完了 ==")

	// データ投入
	log.Println("== DB データ投入 ==")
	_, err = db.Exec(`INSERT INTO users 
		(account, name, password)
		VALUES
		('test', 'テスト太郎', 'test'),
		('ichiro_yamada', '山田一郎', 'yamadA01')
	`)
	if err != nil {
		outputErrorLog("データ投入 失敗", err)
	}
	log.Println("== DB データ投入 完了 ==")

	log.Println("========== DB 初期化 完了 ==========")
}

// CloseDB is method to close DB connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// ConnDB is method to get DB connection
func ConnDB() *sql.DB {
	if db == nil {
		conn := getConnectionString()
		var err error
		db, err = sql.Open("mysql", conn)
		if err != nil {
			outputErrorLog("DB 接続 エラー", err)
		}
	}
	return db
}

// getConnectionString is method to generate DB connection string
func getConnectionString() string {
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s", "root", "security", "tcp", "db", "3306", "security", "?parseTime=true")
}

// output error log and stop app
func outputErrorLog(message string, err error) {
	log.Println(message)
	log.Fatalln(err)
}
