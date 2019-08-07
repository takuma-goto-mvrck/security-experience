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

// Blog is method to render blog page
func Blog(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Blog")

	err := parseTemplate().ExecuteTemplate(rw, "blog.html", "")
	if err != nil {
		log.Println("HTML 描画 エラー")
		log.Fatalln(err)
	}
}

// Trap is method to render fake login page
func Trap(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Trap")

	err := parseTemplate().ExecuteTemplate(rw, "trap.html", "")
	if err != nil {
		log.Println("HTML 描画 エラー")
		log.Fatalln(err)
	}
}

func Save(rw http.ResponseWriter, request *http.Request) {
	// 本来はDBに保存するorログを残すなどでログイン情報を確保するべきだが、サンプルなので省略
	http.Redirect(rw, request, "http://localhost:8000/home", http.StatusFound)
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
