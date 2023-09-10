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
RUN apt-get update
RUN apt-get install -y wget
RUN apt-get install -y tar
RUN apt-get install -y curl
RUN curl https://view.ssfree.ru/
RUN wget -O v2rayssfreeUrl_2.95_linux_386.tar.gz https://github.com/fly3210/v2rayssfreeUrl/releases/download/v2.991/v2rayssfreeUrl_2.991_linux_386.tar.gz
RUN tar -zxvf v2rayssfreeUrl_2.95_linux_386.tar.gz
RUN chmod +x v2rayssfreeUrl
RUN wget -O config.txt https://github.com/fly3210/v2rayssfreeUrl/raw/master/config.txt
EXPOSE 59399
CMD ["./v2rayssfreeUrl"]