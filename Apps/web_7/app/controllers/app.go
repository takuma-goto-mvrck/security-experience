package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/takuma-goto310/security-experience/Apps/web_7/app/db"
	"github.com/takuma-goto310/security-experience/Apps/web_7/app/models"
)

// Cookieに格納するセッション情報のキー
const sessionIDName = "sid"

// Index is method to render Top page.
func Index(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Index")

	tmpl := parseTemplate()

	// Cookieからセッション情報取得
	log.Println("セッション情報取得")
	cookie, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報がないorエラーの場合、TOP画面表示
		log.Println("Cookie 取得 失敗")
		// TOP画面表示
		err = tmpl.ExecuteTemplate(rw, "index.html", "")
		if err != nil {
			outputErrorLog("HTML 描画 エラー", err)
		}
		return
	}
	log.Println("セッション情報：", cookie.Value)

	// セッション情報があればHOME画面に飛ばしてユーザ認証させる
	sessionID := cookie.Value
	if sessionID != "" {
		http.Redirect(rw, request, "/home", http.StatusFound)
		return
	}

	// セッション情報がなければTOP画面表示
	err = tmpl.ExecuteTemplate(rw, "index.html", "")
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
	}
}

// Login is method to authenticate user
func Login(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Login")

	// フォームデータのパース
	err := request.ParseForm()
	if err != nil {
		outputErrorLog("フォーム パース 失敗", err)
	}

	// リクエストデータ取得
	account := request.Form.Get("account")
	password := request.Form.Get("password")
	log.Println("ユーザ：", account)

	// ユーザデータ取得しモデルデータに変換
	dbm := db.ConnDB()
	user := new(models.User)
	row := dbm.QueryRow("select account, name, password from users where account = ?", account)
	if err = row.Scan(&user.Name, &user.Account, &user.Password); err != nil {
		outputErrorLog("ユーザ データ変換 失敗", err)
	}

	// ユーザのパスワード認証
	if user.Password != password {
		log.Println("ユーザ パスワード照合 失敗")
		http.Redirect(rw, request, "/index", http.StatusFound)
		return
	}

	log.Println("認証 成功")

	// 認証が通ったら、セッション情報をDBに保存
	sessionID := generateSessionID(account)
	log.Println("生成したセッションID：", sessionID)
	now := time.Now()
	result, err := dbm.Exec(`INSERT INTO sessions
		(sessionID, account, expireDate)
		VALUES
		(?, ?, ?)
		`, sessionID, account, now.Add(1*time.Hour))
	num, err := result.RowsAffected()
	if err != nil || num == 0 {
		outputErrorLog("セッション データ保存 失敗", err)
	}

	log.Println("セッション データ保存 成功")

	// クッキーにセッション情報付与
	cookie := &http.Cookie{
		Name:  sessionIDName,
		Value: sessionID,
	}
	http.SetCookie(rw, cookie)

	// HOME画面に遷移
	http.Redirect(rw, request, "/home", http.StatusFound)
}

// Home is method to render Home page.
func Home(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Home")

	// クッキーからセッション情報取得しDB照合
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		log.Println("Cookie 取得 失敗")
		log.Println(err)
		http.Redirect(rw, request, "/index", http.StatusFound)
		return
	}
	log.Println("クッキー 取得 成功")
	log.Println("セッション情報：", sessionID.Value)

	session := new(models.Session)
	dbm := db.ConnDB()
	row := dbm.QueryRow("select sessionID, account from sessions where sessionID = ?", sessionID.Value)
	if err = row.Scan(&session.SessionID, &session.Account); err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		log.Println("セッション 取得 失敗")
		log.Println(err)
		// Cookieクリア
		clearCookie(rw)
		http.Redirect(rw, request, "/index", http.StatusFound)
		return
	}

	// セッションを所有するアカウントが存在するかチェック
	user := new(models.User)
	row = dbm.QueryRow("select account, name from users where account = ?", session.Account)
	if err = row.Scan(&user.Account, &user.Name); err != nil {
		// ユーザの照合に失敗した場合TOP画面に遷移
		log.Println("ユーザ 照合 失敗")
		log.Println(err)
		http.Redirect(rw, request, "/index", http.StatusFound)
		return
	}

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
		log.Println("Cookie 取得 失敗")
		log.Println(err)
		// Cookieクリア
		clearCookie(rw)
		http.Redirect(rw, request, "/index", http.StatusFound)
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
	http.Redirect(rw, request, "/index", http.StatusFound)
}

// generate sessionID
func generateSessionID(data string) string {
	hashedSessionID := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hashedSessionID[:])
}

// parse HTML
func parseTemplate() *template.Template {
	tmpl, err := template.ParseGlob("./app/views/*.html")
	if err != nil {
		outputErrorLog("HTML パース 失敗", err)
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
	log.Fatalln(err)
}
