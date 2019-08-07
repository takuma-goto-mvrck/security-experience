package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/takuma-goto310/security-experience/Apps/web_8/app/db"
	"github.com/takuma-goto310/security-experience/Apps/web_8/app/models"
)

// Cookieに格納するセッション情報のキー
const sessionIDName = "sid"

// Index is method to render Top page.
func Index(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Index")

	tmpl := parseTemplate()
	err := tmpl.ExecuteTemplate(rw, "index.html", "")
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
		log.Fatalln("強制終了")
	}
}

// Enter is method to render Login page
func Enter(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Enter")

	// クエリ文字列からURL取得
	url := request.URL.Query().Get("url")
	log.Println(url)

	// LOGIN画面表示
	tmpl := parseTemplate()
	err := tmpl.ExecuteTemplate(rw, "enter.html", struct {
		URL string
	}{
		URL: url,
	})
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
		log.Fatalln("強制終了")
	}
}

// Login is method to authenticate user
func Login(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Login")

	// フォームデータのパース
	err := request.ParseForm()
	if err != nil {
		toEnter(rw, request, "フォーム パース 失敗", err)
		return
	}

	// リクエストデータ取得
	account := request.Form.Get("account")
	password := request.Form.Get("password")
	url := request.Form.Get("url")
	log.Println("ユーザ：", account)
	log.Println("リダイレクト先：", url)

	// ユーザデータ取得しモデルデータに変換
	dbm := db.ConnDB()
	user := new(models.User)
	row := dbm.QueryRow("select account, name, password from users where account = ?", account)
	if err := row.Scan(&user.Account, &user.Name, &user.Password); err != nil {
		toEnter(rw, request, "ユーザ データ変換 失敗", err)
		return
	}

	// ユーザのパスワード認証
	if user.Password != password {
		// 認証に失敗した場合、TOP画面表示
		toEnter(rw, request, "ユーザ パスワード照合 失敗", err)
		return
	}

	log.Println("認証 成功")

	// 既存のセッション情報を取得
	_, err = dbm.Exec("delete from sessions where account = ?", user.Account)
	if err != nil {
		// セッション削除に失敗したらTOP画面に遷移
		toEnter(rw, request, "セッション 削除 失敗", err)
		return
	}

	// セッション情報を新規登録
	sessionID, err := generateRandomSessionID(32)
	if err != nil {
		// セッらション生成に失敗したらTOP画面表示
		toEnter(rw, request, "セッション 生成 失敗", err)
		return
	}
	now := time.Now()
	_, err = dbm.Exec(`insert into sessions
		(sessionID, account, expireDate)
		values
		(?, ?, ?)`,
		sessionID,
		user.Account,
		now.Add(1*time.Hour),
	)
	if err != nil {
		// セッション情報保存に失敗したらTOP画面表示
		toEnter(rw, request, "セッション 保存 失敗", err)
		return
	}
	log.Println("新規作成したセッションID：", sessionID)

	// Cookieにセット
	cookie := &http.Cookie{
		Name:  sessionIDName,
		Value: sessionID,
	}
	http.SetCookie(rw, cookie)

	// 指定されたURLに遷移
	http.Redirect(rw, request, url, http.StatusFound)
}

// Home is method to render Home page.
func Home(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Home")

	// フォームデータのパース
	err := request.ParseForm()
	if err != nil {
		toEnter(rw, request, "フォーム パース 失敗", err)
	}

	// ユーザ情報取得
	user := getAccount(rw, request)
	log.Println("ユーザ情報：", user)

	// ユーザの照合ができればHOME画面表示
	tmpl := parseTemplate()
	err = tmpl.ExecuteTemplate(rw, "home.html", struct {
		User *models.User
	}{
		User: user,
	})
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
	}
}

// Logout is method to delete session
func Logout(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Logout")

	// Cookieからセッション情報取得
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		toEnter(rw, request, "Cookie 取得 失敗", err)
		return
	}
	log.Println("Cookie 取得 成功")
	log.Println("セッション情報：", sessionID.Value)

	// セッション情報を削除
	dbm := db.ConnDB()
	_, err = dbm.Exec("delete from sessions where sessionID = ?", sessionID.Value)
	if err != nil {
		log.Println("セッション 削除 失敗")
	} else {
		log.Println("セッション 削除 成功")
		log.Println("削除したセッションID：", sessionID.Value)
	}

	// CookieクリアしてTOP画面表示
	clearCookie(rw)
	http.Redirect(rw, request, "/enter", http.StatusFound)
}

// Preview is method to render preview page
func Preview(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Preview")

	tmpl := parseTemplate()
	err := tmpl.ExecuteTemplate(rw, "preview.html", "")
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
		log.Fatalln("強制終了")
	}
}

// Product is method to render preview page
func Product(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Product")

	tmpl := parseTemplate()
	err := tmpl.ExecuteTemplate(rw, "product.html", "")
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
		log.Fatalln("強制終了")
	}
}

// generate random sessionID
func generateRandomSessionID(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// parse HTML
func parseTemplate() *template.Template {
	tmpl, err := template.ParseGlob("./app/views/*.html")
	if err != nil {
		outputErrorLog("HTML パース 失敗", err)
		log.Fatalln("強制終了")
	}
	return tmpl
}

// delete session information in Cookie
func clearCookie(rw http.ResponseWriter) {
	log.Println("Cookie初期化")
	cookie := &http.Cookie{
		Name:  sessionIDName,
		Value: "",
	}
	http.SetCookie(rw, cookie)
}

// output error log
func outputErrorLog(message string, err error) {
	log.Println(message)
	log.Println(err)
}

// LOGIN画面に遷移
func toEnter(rw http.ResponseWriter, request *http.Request, message string, err error) {
	outputErrorLog(message, err)
	clearCookie(rw)
	http.Redirect(rw, request, "/enter", http.StatusFound)
}

// セッション情報からユーザ情報を引き当てる
func getAccount(rw http.ResponseWriter, request *http.Request) *models.User {
	// クッキーからセッション情報取得しDB照合
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		toEnter(rw, request, "Cookie 取得 失敗", err)
		return nil
	}
	log.Println("クッキー 取得 成功")
	log.Println("リクエストのセッション情報：", sessionID.Value)

	session := new(models.Session)
	dbm := db.ConnDB()
	row := dbm.QueryRow("select sessionID, account from sessions where sessionID = ?", sessionID.Value)
	if err = row.Scan(&session.SessionID, &session.Account); err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		toEnter(rw, request, "セッション 取得 失敗", err)
		return nil
	}
	log.Println("DBのセッション情報(セッションID)：", session.SessionID)
	log.Println("DBのセッション情報(アカウント)：", session.Account)

	// セッションを所有するアカウントの情報を取得
	user := new(models.User)
	row = dbm.QueryRow("select account, name from users where account = ?", session.Account)
	if err = row.Scan(&user.Account, &user.Name); err != nil {
		// ユーザの照合に失敗した場合TOP画面に遷移
		toEnter(rw, request, "ユーザ 照合 失敗", err)
		return nil
	}

	return user
}
