## テーマ
 CSRF攻撃

## 問題

アプリ内にはCSRF攻撃を受ける危険性のある画面があります。

どんな脅威があるかを特定し、対策方法を検討してください。

## 前提条件
- 自身のアカウントあり
- 画面構成
    - ログイン画面
    - 掲示板画面
- 攻撃用サーバはすでに準備されている（http://localhost:8010）
- 攻撃用サーバに配置するCSRF用のソースはすでに準備されている
    - ソース設置により対象ユーザのメールアドレスを攻撃者指定のメールアドレスに変更可能
    - 上記を実行するためのURLは`http://localhost:8010/csrf`
    - 各ソースが置かれているディレクトリのパスおよびどこに設置するか、については下記参照
- 攻撃者の仕掛けた罠に一般ユーザーが引っかかるという流れを再現する（一般ユーザが引っかかる流れを再現するために、一般ユーザのアカウントを下記に記載する）

### [注意事項]
アプリの起動コマンドは、CSRF用ソース配置後に実行するようにしてください（起動後に配置しても反映されないため）

誤って先に実行してしまった場合は、一度停止させて、ソース配置後に再度起動してください。

## CSRFに使用するソース
| No. | 各ソースのパス | 各ソースの設置場所 |
| --- | --- | --- |
|1|security-experience/Apps/web_6/trap_example/app.go|security-experience/Apps/web_trap/app/controllers/app.go|
|2|security-experience/Apps/web_6/trap_example/main.go|security-experience/Apps/web_trap/app/main.go|
|3|security-experience/Apps/web_6/trap_example/csrf.html|security-experience/Apps/web_trap/app/views/csrf.html|

## アカウント情報
- 攻撃者
    - アカウント名 : test
    - パスワード : test

- 一般ユーザー
    - アカウント名 : yamada
    - パスワード : yamada
    - ※注意※ 一般ユーザのアカウントは、本来攻撃者が知り得ない情報。CSRF攻撃が成功するかを確認するために提供しているので、動作確認以外の目的で利用しないこと