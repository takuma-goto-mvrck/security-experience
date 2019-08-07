package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takuma-goto310/security-experience/Apps/web_4/app/controllers"
	"github.com/takuma-goto310/security-experience/Apps/web_4/app/db"
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
	// LOGIN処理
	http.HandleFunc("/login", controllers.Login)
	// HOME画面
	http.HandleFunc("/home", controllers.Home)
	// COMMENT処理
	http.HandleFunc("/comment", controllers.Comment)
	// LOGOUT処理
	http.HandleFunc("/logout", controllers.Logout)
}
