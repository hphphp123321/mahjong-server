#!/bin/bash

# 读取 YAML 文件并获取 port 值
PORT=$(yq e '.server.port' ./config/app.yaml)

# 输出 port 值
echo "端口号: $PORT"


# 检查 logs 目录是否存在，如果不存在则创建
if [ ! -d "./logs" ]; then
    mkdir ./logs
fi

# 检查是否有正在运行的 mahjong-server 进程，并在存在的情况下杀掉该进程
PID=$(pgrep -f "./build/mahjong-server")
if [ ! -z "$PID" ]; then
    kill $PID
fi

# 重新编译
make build

# 循环直到端口不再被占用
while nc -z localhost $PORT; do   
    echo "等待端口 $PORT 释放..."
    sleep 1
done

# 以后台模式启动 mahjong-server 并将日志输出到 logs/run.log
nohup ./build/mahjong-server > logs/run.log 2>&1 &

echo "mahjong-server 已启动。"
