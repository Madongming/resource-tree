#!/bin/sh

readonly CONTAINER_NAME="resource-tree-test-mysql"
readonly HOST_NAME=`hostname`
readonly RUN_DOCKER_MYSQL_DAEMON="docker run \
       --name $CONTAINER_NAME \
       --rm \
       -e MYSQL_ROOT_PASSWORD=123456 \
       -d \
       -p3306:3306 \
       mysql"
readonly RUN_DOCKER_MYSQL_CMD="docker run \
       --rm \
       mysql \
       mysql \
       -h$HOST_NAME \
       -uroot \
       -p123456"

echo 'Start vet check.'
go vet ./...
echo 'Finished.'
echo ''

# 停止mysql
docker kill $CONTAINER_NAME > /dev/null 2>&1

echo 'Start mysql.'
# 启动mysql，做为测试使用
$RUN_DOCKER_MYSQL_DAEMON

# 等待服务启动
echo 'Waiting for start mysqld...'
while true;do
    $RUN_DOCKER_MYSQL_CMD -e "\q" > /dev/null 2>&1
    if [ $? -eq 0 ];then
	break
    fi
    sleep 2
done
echo 'Started!'
echo ''

# 创建库
echo 'Created database `resource`'
$RUN_DOCKER_MYSQL_CMD -e "create database resource default charset utf8" > /dev/null 2>&1

# 判断是否启动成功
if [ $? -ne 0 ];then
    SI=$?
    docker kill $CONTAINER_NAME > /dev/null 2>&1
    exit $SI
fi
echo 'Created.'
echo ''

# Set env
export RESOURCE_TREE_MODE="test"
export CONFIG_FILE="${PWD}/conf/config.yml"
export LOG_CONFIG_FILE="${PWD}/conf/test.xml"

# 开始测试
echo 'Cleaning test cache...'
go clean -testcache > /dev/null 2>&1
echo 'Clean up.'
echo ''

echo 'Start test:'
go test -v ./model
go test -v ./dao
