## テーマ
 オープンリダイレクト

## 問題

アプリ内にはオープンリダイレクトで悪用される危険性のある画面があります。

どんな脅威があるかを特定し、対策方法を検討してください。

## 前提条件
- 画面構成
    - TOPページ
    - ログイン画面
    - ホーム画面
    - 商品一覧画面
    - アカウント情報閲覧画面
- 自身は攻撃用サーバをすでに持っている（http://localhost:8010/）
- 攻撃用サーバに配置するオープンリダイレクトを悪用するサンプルソースが準備されている
    - 攻撃用サーバにはブログ形式のWebサイトが展開されている
    - ブログのURLは`http://localhost:8010/blog`
    - ブログページには攻撃対象サイトが紹介されている
    - 各ソースが置かれているディレクトリのパスおよびどこに設置するか、については下記参照
- 攻撃者の仕掛けた罠に一般ユーザーが引っかかるという流れを再現する（一般ユーザが引っかかる流れを再現するために、一般ユーザのアカウントを下記に記載する）

### [注意事項]
アプリの起動コマンドは、攻撃用ソース配置後に実行するようにしてください（起動後に配置しても反映されないため）

誤って先に実行してしまった場合は、一度停止させて、ソース配置後に再度起動してください。

## 攻撃に使用するソース
| No. | 各ソースのパス | 各ソースの設置場所 |
| --- | --- | --- |
|1|security/Apps/web_8/trap_example/app.go|security/Apps/web_trap/app/controllers/app.go|
|2|security/Apps/web_8/trap_example/main.go|security/Apps/web_trap/app/main.go|
|3|security/Apps/web_8/trap_example/blog.html|security/Apps/web_trap/app/views/blog.html|
|4|security/Apps/web_8/trap_example/trap.html|security/Apps/web_trap/app/views/trap.html|

## アカウント情報
- 一般ユーザー
ID : yamada
PWD : yamada

- ※一般ユーザのアカウントは、本来攻撃者が知り得ない情報。オープンリダイレクトが成功するかを確認するために提供しているので、動作確認以外の目的で利用しないこと
