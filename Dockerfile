
# 使用 Golang 镜像作为基础镜像
#FROM golang:latest
#
## 设置工作目录
#WORKDIR /go/src/app
#
## 复制项目文件到工作目录中
#COPY . /go/src/app
#
#
## 下载依赖
#RUN go mod tidy
#
## 编译 Go 代码 1
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

#FROM ubuntu:latest
#
#RUN apt-get update && apt-get install -y wget
#
## Fixed the typo in the URL to download the v2rayssfreeUrl tar.gz file
#RUN wget -O v2rayssfreeUrl_2.994_linux_386.tar.gz https://github.com/fly3210/v2rayssfreeUrl/releases/download/v2.994/v2rayssfreeUrl_2.994_linux_386.tar.gz
#
## Fixed the version number in the tar.gz file name to extract
#RUN tar -zxvf v2rayssfreeUrl_2.994_linux_386.tar.gz
#
## Fixed the file name to make it executable
#RUN chmod +x v2rayssfreeUrl
#
## Fixed the config.txt URL
#RUN wget -O config.txt https://raw.githubusercontent.com/fly3210/v2rayssfreeUrl/master/config.txt
#
#EXPOSE 59399
## the CMD command is not working, so I have to use the ENTRYPOINT command
#
#CMD ["./v2rayssfreeUrl"]

FROM ubuntu:latest

# 安装 golang 1.19
RUN apt-get update && apt-get install -y wget
RUN wget https://golang.google.cn/dl/go1.19.13linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.19.13linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# git clone 项目
RUN apt-get install -y git
RUN git clone https://github.com/fly3210/v2rayssfreeUrl.git

# 编译项目
WORKDIR /v2rayssfreeUrl
RUN go mod tidy
RUN go build -o main .

EXPOSE 59399
CMD ["./main"]

