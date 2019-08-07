package controllers

import (
	"html/template"
	"log"
	"net/http"
)

// Index is method to render Top page.
func Index(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Index")

	// TOP画面表示
	err := parseTemplate().ExecuteTemplate(rw, "index.html", "")
	if err != nil {
		log.Println("HTML 描画 エラー")
		log.Fatalln(err)
	}
}

// Csrf is method to attack specify server
func Csrf(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Csrf")

	err := parseTemplate().ExecuteTemplate(rw, "csrf.html", "")
	if err != nil {
		log.Println("HTML 描画 エラー")
		log.Fatalln(err)
	}
}

// parseTemplate is method to parse html
func parseTemplate() *template.Template {
	tmpl, err := template.ParseGlob("./app/views/*.html")
	if err != nil {
		log.Println("HTML パース 失敗")
		log.Fatalln(err)
	}

	return tmpl
}
