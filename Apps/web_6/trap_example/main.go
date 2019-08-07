package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takuma-goto310/security-experience/Apps/web_trap/app/controllers"
)

func main() {
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
	// CSRF処理
	http.HandleFunc("/csrf", controllers.Csrf)
}
