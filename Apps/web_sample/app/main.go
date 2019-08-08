package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takuma-goto-mvrck/security-experience/Apps/web_sample/app/controllers"
	"github.com/takuma-goto-mvrck/security-experience/Apps/web_sample/app/db"
)

func main() {
	db.InitDB()
	setRoute()
	// 指定したポートをListen
	http.ListenAndServe(":8080", nil)
}

// ルーティング設定
func setRoute() {
	// 静的ファイルのルーティング
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	// TOP画面
	http.HandleFunc("/index", controllers.Index)
	// 一覧画面
	http.HandleFunc("/list", controllers.List)
}
