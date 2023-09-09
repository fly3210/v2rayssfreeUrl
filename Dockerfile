# 使用 Golang 镜像作为基础镜像
FROM golang:latest

# 设置工作目录
WORKDIR /go/src/app

# 复制项目文件到工作目录中
COPY . /go/src/app

# 编译 Go 代码
RUN go build -o main .

# 打印文件列表
RUN ls

# 暴露端口（如果你的应用需要监听端口的话）
EXPOSE 59399

# 运行应用
CMD ["./main"]