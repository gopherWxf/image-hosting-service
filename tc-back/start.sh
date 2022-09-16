#!/bin/bash
echo
echo ============= fastdfs ==============
# 关闭已启动的 tracker 和 storage
./fastdfs.sh stop
# 启动 tracker 和 storage
./fastdfs.sh all
# 重启所有的 cgi程序
echo
echo ============= nginx ==============
./nginx.sh stop
# 启动nginx
./nginx.sh start
# 关闭redis
echo
echo ============= redis ==============
./redis.sh stop
# 启动redis
./redis.sh start
echo
echo ============= Go ==============
go build -o image-hosting-service-back
# 启动redis
./image-hosting-service-back
