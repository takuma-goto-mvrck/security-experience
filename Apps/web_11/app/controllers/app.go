package controllers

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// History is a type that has transition path information
type History struct {
	Cwd  string
	Path string
}

// 各ユーザ用のファイル置き場
const baseDirPath = "storage/files/"

// Index is method to render Top page.
func Index(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Index")

	tmpl := parseTemplate()

	// TOP画面表示
	toIndex(tmpl, rw, request)
}

// Home is method to render Home page.
func Home(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Home")

	tmpl := parseTemplate()

	// フォームデータのパース
	err := request.ParseForm()
	if err != nil {
		outputErrorLog("フォーム パース 失敗", err)
	}

	// ユーザID取得
	userID, err := getUserIDFromPath(tmpl, rw, request)
	if err != nil {
		log.Println("アクセス 権限 エラー")
		toIndex(tmpl, rw, request)
		return
	}

	// リクエストデータ取得
	query := request.Form.Get("query")

	// 履歴を生成（必要ではないかもしれないけど、操作性悪いので追加）
	historyNames := strings.Split(query, "/")
	path := []string{}
	histories := make([]History, len(historyNames))
	for i, historyName := range historyNames {
		log.Println("履歴：", historyName)
		path = append(path, historyName)
		histories[i].Cwd = historyName
		histories[i].Path = strings.Join(path, "/")
	}
	if query != "" {
		query += "/"
	}

	// ファイル一覧取得
	dirPath := baseDirPath + userID + "/" + query
	log.Println("検索クエリ：", query)
	log.Println("ディレクトリパス：", dirPath)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		// ファイル一覧取得に失敗
		log.Println("ファイル一覧 取得 失敗")
		log.Println(err)
		files = nil
		histories = nil
	}

	// HOME画面表示
	err = tmpl.ExecuteTemplate(rw, "home.html", struct {
		UserID    string
		Query     string
		DirPath   string
		Files     []os.FileInfo
		Histories []History
	}{
		UserID:    userID,
		Query:     query,
		DirPath:   dirPath,
		Files:     files,
		Histories: histories,
	})
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
	}
}

// Browse is method to browse user's file
func Browse(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Browse")

	tmpl := parseTemplate()

	// ユーザID取得
	userID, err := getUserIDFromPath(tmpl, rw, request)
	if err != nil {
		log.Println("アクセス 権限 エラー")
		toIndex(tmpl, rw, request)
		return
	}

	// リクエストデータ取得
	query := request.URL.Query().Get("query")

	// ファイル表示
	filePath := baseDirPath + userID + "/" + query
	log.Println("表示するファイルパス：", filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		outputErrorLog("ファイル 表示 失敗", err)
	}

	rw.Write(data)
}

// parse HTML
func parseTemplate() *template.Template {
	tmpl, err := template.ParseGlob("./app/views/*.html")
	if err != nil {
		outputErrorLog("HTML パース 失敗", err)
	}
	return tmpl
}

// output error log
func outputErrorLog(message string, err error) {
	log.Println(message)
	log.Fatalln(err)
}

// render top page
func toIndex(tmpl *template.Template, rw http.ResponseWriter, request *http.Request) {
	err := tmpl.ExecuteTemplate(rw, "index.html", "")
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
	}
}

// get user id
func getUserIDFromPath(tmpl *template.Template, rw http.ResponseWriter, request *http.Request) (string, error) {

	requestPath := request.URL.Path
	log.Println("アクセスURL：", requestPath)
	splitedPath := strings.Split(requestPath, "/")

	// 簡単にするため、ユーザIDが1以外の場合はTOP画面表示
	if len(splitedPath) < 2 || splitedPath[2] != "1" {
		err := errors.New("アクセス 権限 エラー")
		return "", err
	}
	return splitedPath[2], nil
}
