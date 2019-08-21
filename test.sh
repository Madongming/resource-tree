#!/bin/sh

# 停止mysql
docker kill mysql

# 启动mysql，做为测试使用
docker run \
       --name mysql \
       --rm \
       -e MYSQL_ROOT_PASSWORD=123456 \
       -d \
       -p3306:3306 \
       mysql

# 等待服务启动
echo 'Waiting for start mysqld...'
while true;do
    mysql -h127.0.0.1 -uroot -p123456 -e "\q" > /dev/null 2>&1
    if [ $? -eq 0 ];then
	break
    fi
    sleep 2
done

# 创建库
echo 'Created database `resource`'
mysql -h127.0.0.1 -uroot -p123456 -e "create database resource default charset utf8" > /dev/null 2>&1

# 判断是否启动成功
if [ $? -ne 0 ];then
    SI=$?
    docker kill mysql
    exit $SI
fi

# Set env
export RESOURCE_TREE_MODE="test"
export CONFIG_FILE="/Users/madongming/Documents/GITREPO/mdm/resource-tree/conf/config.yml"
export LOG_CONFIG_FILE="/Users/madongming/Documents/GITREPO/mdm/resource-tree/conf/test.xml"

# 开始测试
go clean -testcache
go test -v ./model
go test -v ./dao
