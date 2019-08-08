# このアプリについて
---
## アプリの目的
このアプリはセキュリティ対策意識と対策技術の向上を狙いとしたアプリです。

起動することで、様々な脆弱性を含んだwebサイトが立ち上がるようになっています。

立ち上がったwebサイトにて、使用者の皆様が脆弱性を発見、攻撃してもらうことで、実際の攻撃方法を理解し、セキュリティ対策に役立ててほしいと思います。

## 学習内容
以下の脆弱性をもつwebアプリを含んでいます。
- XSS (クロスサイトスクリプティング)
- SQLインジェクション
- OSコマンドインジェクション
- メールヘッダ・インジェクション
- オープンリダイレクト
- CSRF (クロスサイトリクエストフォージェリ)
- セッション管理の不備
- ディレクトリ・トラバーサル
- evalインジェクション

参考 :「体系的に学ぶ 安全なWebアプリケーションの作り方　脆弱性が生まれる原理と対策の実践」：著、徳丸浩

## 学習方法
アプリの起動は、下記「アプリの起動方法」に記載しています。
各脆弱性ごとにREADMEが用意されており、攻撃する上で重要な情報が記載されています。
アプリ実行前に各プロジェクトのREADMEをよく読んで学習を開始してください。

## 回答
各アプリごとに回答例を用意しています。
閲覧権限を設けていますので、回答が知りたい方は管理者に閲覧権限を申請してください。


# アプリ内容
|app_no|テーマ|README|
|---|---|---|
|web_4|XSS、エラーメッセージからの情報漏洩)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_4/README.md|
|web_5|SQL呼び出し(SQLインジェクション)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_5/README.md|
|web_6|取り消しの効かない重要な処理(CSRF)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_6/README.md|
|web_7|セッション管理の不備(セッションハイジャック、推測可能なセッションID、セッションIDの固定化)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_7/README.md|
|web_8|リダイレクト処理(オープンリダイレクタ、HTTPヘッダインジェクション)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_8/README.md|
|web_10|メール送信(メールヘッダ・インジェクション)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_10/README.md|
|web_11|ファイルアクセス(ディレクトリ・トラバーサル、意図しないファイル公開)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_11/README.md|
|web_12|OSコマンド(OSコマンド・インジェクション)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_12/README.md|
|web_15|eval(evalインジェクション)|https://github.com/takuma-goto-mvrck/security-experience/blob/master/Apps/web_15/README.md|


# アプリの起動方法
---
## 前提
[Docker](https://docs.docker.com/install/)をインストールする

## 本PJの構成
- Apps (脆弱性を含んだアプリがある。backlogのwiki「各脆弱性体験アプリのテーマ」に記載のapp_noが子階層のディレクトリ名になっている)
- Config （Apps内で使用し、Appsに依存しない設定ファイル）
- Docker （言語毎にDockerFile、docker-composeを記載）
- app.conf （Dockerコンテナ作成起動時に使用。ホスト側での事前実行タスク、Dockerディレクトリ次階層の選択、docker-composeで起動するserviceが記載されたファイル）
- startup.sh （起動シェルスクリプト）
- down.sh （停止シェルスクリプト）

## 起動・停止方法
    # 起動 ※複数のアプリを起動しないでください。起動前には必ず停止処理を行ってください。
    $ git clone https://github.com/takuma-goto-mvrck/security-experience.git
    $ cd security-experience
    $ sh ./startup.sh {$SERVICE_NAME}
    # SERVICE_NAMEには「アプリ内容」に記載した「app_no」を指定してください(例：SQLインジェクションのアプリを立ち上げたい場合　$ sh ./startup.sh web_5)
    
    # 停止
    $ sh ./down.sh
    
    # ベースのDockerイメージも含めて削除する場合
    $ sh ./down.sh --rmi all

## 起動確認
ブラウザのアドレス欄に`http://localhost:8000/index`を入力

（アプリの中には、罠用のサイトが立ち上がるものもあります。罠サイトには`http://localhost:8010/index`でアクセス可能です）

※ 起動後すぐは繋がらない場合があるため、 少し待ってからブラウザを更新してください。
    
## 注意事項
- このアプリケーションは、各脆弱性を実際に体験することで脅威への理解とその対策方法を学ぶツールとして利用しています。
- このアプリケーションを通して得た知識・技術は決して悪用しないようにしてください。
- 各アプリケーションではテーマに沿わない脆弱性が発見される場合もありますが、テーマとなる脆弱性を体験することを目的としているため、テーマ外の脆弱性については考慮しないようにお願いします。
