#!/bin/bash

# 起動しているコンテナを全て削除（起動有無に関わらず全言語で実行する）
cd ./Docker
ls -1 | while read line
do
    cd $line
    docker-compose down --volumes $@
    cd ../
done
