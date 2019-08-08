package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/takuma-goto-mvrck/security-experience/Apps/web_6/app/db"
	"github.com/takuma-goto-mvrck/security-experience/Apps/web_6/app/models"
)

// Cookieに格納するセッション情報のキー
const sessionIDName = "sid"

// template内で使用する関数の登録
var funcMap = template.FuncMap{
	"contains": strings.Contains,
}

// Index is method to render Top page.
func Index(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Index")

	err := parseTemplate().ExecuteTemplate(rw, "index.html", "")
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
		log.Fatalln("エラーのため強制終了")
	}
}

// Login is method to authenticate user
func Login(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Login")

	tmpl := parseTemplate()

	// フォームデータのパース
	err := request.ParseForm()
	if err != nil {
		toIndex(tmpl, rw, request, "フォーム パース 失敗", err)
		return
	}

	// リクエストデータ取得
	account := request.Form.Get("account")
	password := request.Form.Get("password")
	log.Println("ユーザ：", account)

	// ユーザデータ取得しモデルデータに変換
	dbm := db.ConnDB()
	user := new(models.User)
	row := dbm.QueryRow("select id, account, name, password from users where account = ?", account)
	if err = row.Scan(&user.ID, &user.Name, &user.Account, &user.Password); err != nil {
		// 変換に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "ユーザ データ変換 失敗", err)
		return
	}

	// ユーザのパスワード認証
	if user.Password != password {
		// 認証に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "ユーザ パスワード照合 失敗", err)
		return
	}

	log.Println("認証 成功")

	// 新規セッション情報生成
	sessionID, err := generateRandomSessionID(32)
	if err != nil {
		// セッション生成に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "セッション 生成 失敗", err)
		return
	}
	log.Println("生成したセッションID：", sessionID)

	// 生成した新規セッション情報をDBに保存
	now := time.Now()
	result, err := dbm.Exec(`INSERT INTO sessions
		(sessionID, account, expireDate)
		VALUES
		(?, ?, ?)
		`, sessionID, account, now.Add(1*time.Hour))
	num, err := result.RowsAffected()
	if err != nil || num == 0 {
		// セッション情報保存に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "セッション 保存 失敗", err)
		return
	}

	log.Println("セッション 保存 成功")

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

	tmpl := parseTemplate()

	// クッキーからセッション情報取得しDB照合
	user, msg, err := checkSession(tmpl, rw, request)
	if err != nil {
		toIndex(tmpl, rw, request, msg, err)
		return
	}

	// ユーザの照合ができればコメント取得
	dbm := db.ConnDB()
	comments := []models.Comment{}
	rows, err := dbm.Query("select id, comment from comments order by id")
	if err != nil {
		// コメント取得に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "コメント 取得 失敗", err)
		return
	}
	for rows.Next() {
		comment := models.Comment{}
		if err = rows.Scan(&comment.ID, &comment.Comment); err != nil {
			toIndex(tmpl, rw, request, "コメント 変換 失敗", err)
			return
		}
		log.Println("コメント：", comment.Comment)
		comments = append(comments, comment)
	}
	log.Println("コメント 取得 成功")

	// HOME画面表示
	err = tmpl.ExecuteTemplate(rw, "home.html", struct {
		User     *models.User
		Comments []models.Comment
	}{
		User:     user,
		Comments: comments,
	})
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
	}
}

// Comment is method to register comment
func Comment(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Comment")

	// フォームをパース
	err := request.ParseForm()
	if err != nil {
		outputErrorLog("フォーム パース 失敗", err)
		http.Redirect(rw, request, "/home", http.StatusFound)
		return
	}

	// 送信されたコメント取得
	comment := request.Form.Get("comment")
	log.Println("コメント：", comment)

	// コメント保存
	dbm := db.ConnDB()
	// 簡単のため、影響があった件数の判定はしない
	_, err = dbm.Exec(`insert into comments 
		(comment)
		values
		(?)
		`, comment)
	if err != nil {
		// コメント保存失敗
		outputErrorLog("コメント 保存 失敗", err)
	}

	http.Redirect(rw, request, "/home", http.StatusFound)
}

// Mail is method to render Mail page
func Mail(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Mail")

	err := parseTemplate().ExecuteTemplate(rw, "email.html", "")
	if err != nil {
		outputErrorLog("HTML 描画 エラー", err)
		log.Fatalln("エラーのため強制終了")
	}
}

// ChangeMail is method to change mail address
func ChangeMail(rw http.ResponseWriter, request *http.Request) {

	log.Println("call ChangeMail")

	tmpl := parseTemplate()

	// セッション情報の確認
	user, msg, err := checkSession(tmpl, rw, request)
	if err != nil {
		toIndex(tmpl, rw, request, msg, err)
		return
	}
	log.Println("変更するユーザ：", user.ID)

	// 変更するメールアドレス取得
	err = request.ParseForm()
	if err != nil {
		// パースに失敗したらエラー
		toIndex(tmpl, rw, request, "フォーム パース エラー", err)
		return
	}
	email := request.Form.Get("email")
	log.Println("変更するメールアドレス：", email)

	// メールアドレスの変更
	dbm := db.ConnDB()
	result, err := dbm.Exec(`update users
		set email = ?
		where id = ?
	`, email, user.ID)
	if err != nil {
		toIndex(tmpl, rw, request, "メールアドレス 変更 SQLエラー", err)
		return
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		toIndex(tmpl, rw, request, "メールアドレス 変更 行取得失敗", err)
		return
	}
	if rowNum < 1 {
		err = errors.New("SQL実行は終了したが結果件数0件のためエラー")
		toIndex(tmpl, rw, request, "メールアドレス 変更", err)
		return
	}
	log.Println("メールアドレス 変更 成功")

	http.Redirect(rw, request, "/home", http.StatusFound)
}

// Logout is method to delete session
func Logout(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Logout")

	tmpl := parseTemplate()

	// Cookieからセッション情報取得
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "Cookie 取得 失敗", err)
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
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("./app/views/*.html")
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
	log.Println(err)
}

// render Top page
func toIndex(tmpl *template.Template, rw http.ResponseWriter, request *http.Request, message string, err error) {
	outputErrorLog(message, err)
	clearCookie(rw)
	http.Redirect(rw, request, "/index", http.StatusFound)
}

// checkSession is method to check valid session
func checkSession(tmpl *template.Template, rw http.ResponseWriter, request *http.Request) (*models.User, string, error) {
	// クッキーからセッション情報取得しDB照合
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得失敗
		return nil, "Cookie 取得 失敗", err
	}
	log.Println("クッキー 取得 成功")
	log.Println("セッション情報：", sessionID.Value)

	session := new(models.Session)
	dbm := db.ConnDB()
	row := dbm.QueryRow("select sessionID, account from sessions where sessionID = ?", sessionID.Value)
	if err = row.Scan(&session.SessionID, &session.Account); err != nil {
		// セッション情報取得失敗
		return nil, "セッション 取得 失敗", err
	}
	log.Println("DBからのセッション情報取得成功")

	// セッションを所有するアカウントが存在するかチェック
	user := new(models.User)
	row = dbm.QueryRow("select id, account, name, email from users where account = ?", session.Account)
	if err = row.Scan(&user.ID, &user.Account, &user.Name, &user.Email); err != nil {
		// ユーザの照合失敗
		return nil, "ユーザ 照合 失敗", err
	}
	log.Println("ユーザ 照合 成功")

	return user, "", nil
}
