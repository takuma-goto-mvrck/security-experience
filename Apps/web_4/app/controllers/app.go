package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/takuma-goto310/security-experience/Apps/web_4/app/db"
	"github.com/takuma-goto310/security-experience/Apps/web_4/app/models"
)

// TemplateComment is struct for comment to render Home page
type TemplateComment struct {
	Line template.HTML
}

// Cookieに格納するセッション情報のキー
const sessionIDName = "sid"

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
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "Cookie 取得 失敗", err)
		return
	}
	log.Println("クッキー 取得 成功")
	log.Println("セッション情報：", sessionID.Value)

	session := new(models.Session)
	dbm := db.ConnDB()
	row := dbm.QueryRow("select sessionID, account from sessions where sessionID = ?", sessionID.Value)
	if err = row.Scan(&session.SessionID, &session.Account); err != nil {
		// セッション情報取得に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "セッション 取得 失敗", err)
		return
	}
	log.Println("DBからのセッション情報取得成功")

	// セッションを所有するアカウントが存在するかチェック
	user := new(models.User)
	row = dbm.QueryRow("select account, name from users where account = ?", session.Account)
	if err = row.Scan(&user.Account, &user.Name); err != nil {
		// ユーザの照合に失敗した場合、TOP画面表示
		toIndex(tmpl, rw, request, "ユーザ 照合 失敗", err)
		return
	}
	log.Println("ユーザ 照合 成功")

	// ユーザの照合ができればコメント取得
	templateComments := []TemplateComment{}
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
		tc := TemplateComment{
			Line: template.HTML(strconv.Itoa(comment.ID) + " : " + comment.Comment),
		}
		templateComments = append(templateComments, tc)
	}
	log.Println("コメント 取得 成功")

	// HOME画面表示
	err = tmpl.ExecuteTemplate(rw, "home.html", struct {
		User     *models.User
		Comments []TemplateComment
	}{
		User:     user,
		Comments: templateComments,
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
	log.Println(err)
}

// render Top page
func toIndex(tmpl *template.Template, rw http.ResponseWriter, request *http.Request, message string, err error) {
	outputErrorLog(message, err)
	clearCookie(rw)
	http.Redirect(rw, request, "/index", http.StatusFound)
}
