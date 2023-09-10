
## 使用 Golang 镜像作为基础镜像
#FROM golang:latest
#
## 设置工作目录
#WORKDIR /go/src/app
#
## 复制项目文件到工作目录中
#COPY . /go/src/app
#
## 编译 Go 代码
#RUN go build -o main .
#
## 打印文件列表
#RUN ls
#
## 暴露端口（如果你的应用需要监听端口的话）
#EXPOSE 59399
#
## 运行应用
#CMD ["./main"]

FROM ubuntu:latest

# Fixed the typo in the URL to download the v2rayssfreeUrl tar.gz file
RUN wget -O v2rayssfreeUrl_2.994_linux_386.tar.gz https://github.com/fly3210/v2rayssfreeUrl/releases/download/v2.994/v2rayssfreeUrl_2.994_linux_386.tar.gz

# Fixed the version number in the tar.gz file name to extract
RUN tar -zxvf v2rayssfreeUrl_2.994_linux_386.tar.gz

# Fixed the file name to make it executable
RUN chmod +x v2rayssfreeUrl

# Fixed the config.txt URL
RUN wget -O config.txt https://raw.githubusercontent.com/fly3210/v2rayssfreeUrl/master/config.txt

EXPOSE 59399

CMD ["./v2rayssfreeUrl"]
