package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/takuma-goto-mvrck/security-experience/Apps/web_5/app/db"
	"github.com/takuma-goto-mvrck/security-experience/Apps/web_5/app/models"
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

// Login is method to authenticate user
func Login(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Login")

	// フォームデータのパース
	err := request.ParseForm()
	if err != nil {
		toIndex(rw, request, "フォーム パース 失敗", err)
		return
	}

	// リクエストデータ取得
	account := request.Form.Get("account")
	password := request.Form.Get("password")
	log.Println("ユーザ：", account)

	// ユーザデータ取得しモデルデータに変換
	dbm := db.ConnDB()
	user := new(models.User)
	row := dbm.QueryRow("select name, account from users where account = \"" + account + "\" and password = \"" + password + "\"")
	if err := row.Scan(&user.Name, &user.Account); err != nil {
		toIndex(rw, request, "ユーザ データ変換 失敗", err)
		return
	}

	log.Println("認証 成功")

	// 既存のセッション情報を取得
	_, err = dbm.Exec("delete from sessions where account = ?", user.Account)
	if err != nil {
		// セッション削除に失敗したらTOP画面に遷移
		toIndex(rw, request, "セッション 削除 失敗", err)
		return
	}

	// セッション情報を新規登録
	sessionID, err := generateRandomSessionID(32)
	if err != nil {
		// セッらション生成に失敗したらTOP画面表示
		toIndex(rw, request, "セッション 生成 失敗", err)
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
		toIndex(rw, request, "セッション 保存 失敗", err)
		return
	}
	log.Println("新規作成したセッションID：", sessionID)

	// Cookieにセット
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

	// フォームデータのパース
	err := request.ParseForm()
	if err != nil {
		toIndex(rw, request, "フォーム パース 失敗", err)
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

// Search is method to search product
func Search(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Search")

	// 検索文字列取得
	query := request.URL.Query().Get("query")
	log.Println("クエリ：", query)

	// ユーザ情報取得
	user := getAccount(rw, request)

	tmpl := parseTemplate()

	// 商品検索
	message := ""
	dbm := db.ConnDB()
	rows, err := dbm.Query("select number, name, image, price, stock from products where name like \"%" + query + "%\" order by name")
	if err != nil {
		// 商品検索失敗した場合、商品検索画面にメッセージ表示
		outputErrorLog("商品 検索 失敗", err)
		message = "検索に失敗しました。" + err.Error()
		renderSearch(tmpl, rw, request, user, query, message, nil)
		return
	}

	// モデルに格納
	products := []models.Product{}
	for rows.Next() {
		product := models.Product{}
		if err = rows.Scan(&product.Number, &product.Name, &product.Image, &product.Price, &product.Stock); err != nil {
			// 商品データの格納に失敗した場合、商品検索画面にメッセージ表示
			outputErrorLog("商品 格納 失敗", err)
			message = "検索に失敗しました。" + err.Error()
			renderSearch(tmpl, rw, request, user, query, message, nil)
			return
		}
		products = append(products, product)
	}

	if len(products) == 0 {
		message = "検索結果なし"
	}
	log.Println("検索件数：", len(products), "件")

	// 商品検索画面表示
	renderSearch(tmpl, rw, request, user, query, message, products)
}

// Logout is method to delete session
func Logout(rw http.ResponseWriter, request *http.Request) {

	log.Println("call Logout")

	// Cookieからセッション情報取得
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		toIndex(rw, request, "Cookie 取得 失敗", err)
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

// TOP画面に遷移
func toIndex(rw http.ResponseWriter, request *http.Request, message string, err error) {
	outputErrorLog(message, err)
	clearCookie(rw)
	http.Redirect(rw, request, "/index", http.StatusFound)
}

func getAccount(rw http.ResponseWriter, request *http.Request) *models.User {
	// クッキーからセッション情報取得しDB照合
	sessionID, err := request.Cookie(sessionIDName)
	if err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		toIndex(rw, request, "Cookie 取得 失敗", err)
		return nil
	}
	log.Println("クッキー 取得 成功")
	log.Println("セッション情報：", sessionID.Value)

	session := new(models.Session)
	dbm := db.ConnDB()
	row := dbm.QueryRow("select sessionID, account from sessions where sessionID = ?", sessionID.Value)
	if err = row.Scan(&session.SessionID, &session.Account); err != nil {
		// セッション情報取得に失敗した場合TOP画面に遷移
		toIndex(rw, request, "セッション 取得 失敗", err)
		return nil
	}

	// セッションを所有するアカウントの情報を取得
	user := new(models.User)
	row = dbm.QueryRow("select account, name from users where account = ?", session.Account)
	if err = row.Scan(&user.Account, &user.Name); err != nil {
		// ユーザの照合に失敗した場合TOP画面に遷移
		toIndex(rw, request, "ユーザ 照合 失敗", err)
		return nil
	}

	return user
}

// 検索画面表示
func renderSearch(tmpl *template.Template, rw http.ResponseWriter, request *http.Request, user *models.User, query string, message string, products []models.Product) {

	log.Println("検索画面表示")

	err := tmpl.ExecuteTemplate(rw, "search.html", struct {
		User     *models.User
		Query    string
		Message  string
		Products []models.Product
	}{
		User:     user,
		Query:    query,
		Message:  message,
		Products: products,
	})
	if err != nil {
		// 描画エラーの場合、HOME画面表示
		outputErrorLog("HTML 描画 エラー", err)
		http.Redirect(rw, request, "/home", http.StatusFound)
	}
}
