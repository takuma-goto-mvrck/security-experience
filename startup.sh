#!/bin/bash

source ./app.conf

# 環境変数設定
docker_dir_key="$1_docker_dir"
docker_dir=$(eval echo '$'$docker_dir_key)
readonly ENV_FILE="./Docker/$docker_dir/.env"
echo SERVICE_NAME=$1 > $ENV_FILE

# 登録されたタスクを実行
PRE_IFS=$IFS
task_key="$1_tasks"
tasks=$(eval echo '$'{$task_key[@]})
echo $tasks
IFS=$','
for task in ${tasks[@]}; do
    eval $task
done

# Docker起動
services_key="$1_docker_compose_services"
services=$(eval echo '$'{$services_key[@]})
cd ./Docker/$docker_dir
eval "docker-compose build $services"
eval "docker-compose up -d $services"

IFS=$PRE_IFS
