# 使用官方的 Go 镜像作为构建环境
FROM golang:alpine

# 复制项目文件到容器中
COPY . /mahjong-server

# 设置工作目录
WORKDIR /mahjong-server

# 下载相关依赖
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy

# 运行 go build 生成二进制文件
RUN go build -o build/mahjong-server cmd/main.go

# 暴露服务端口
EXPOSE 16548

# 设置容器启动时执行的命令
CMD ["./build/mahjong-server"]
