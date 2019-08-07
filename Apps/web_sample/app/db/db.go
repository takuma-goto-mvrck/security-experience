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
	conn := fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s", "root", "security", "tcp", "db", "3306", "security", "?parseTime=true")
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println("DB 接続 エラー")
		log.Fatalln(err)
	}
	log.Println("== DB 接続 成功 ==")

	// 既存のテーブル削除
	log.Println("== DB 削除 ==")
	db.Exec("DROP TABLE users")
	log.Println("== DB 削除 成功 ==")

	// ユーザテーブル生成
	log.Println("== DB テーブル 作成 ==")
	_, err = db.Exec("CREATE TABLE users ( id INTEGER AUTO_INCREMENT PRIMARY KEY, name VARCHAR(32) NOT NULL )")
	if err != nil {
		log.Println("ユーザテーブル作成失敗")
		log.Fatalln(err)
	}
	log.Println("== DB テーブル 作成 完了 ==")

	// データ投入
	log.Println("== DB データ投入 ==")
	_, err = db.Exec("INSERT INTO users (name) value ('Tanaka')")
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
	return db
}
