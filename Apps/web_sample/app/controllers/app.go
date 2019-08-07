package controllers

import (
	"html/template"
	"log"
	"net/http"
)

// Index is method for TOP page
func Index(rw http.ResponseWriter, request *http.Request) {
	// テンプレート取得
	template, err := template.ParseFiles("./app/views/index.html", "./app/views/header.html")
	if err != nil {
		log.Println("HTMLパースエラー")
		log.Fatalln(err)
	}
	// テンプレート描画
	err = template.Execute(rw, "index.html")
	if err != nil {
		log.Println("テンプレート描画エラー")
		log.Fatalln(err)
	}
}

// List is method for LIST page
func List(rw http.ResponseWriter, request *http.Request) {

	rw.Write([]byte("hello"))
}
