# alpine:3.15.0

# 参数
ARG DOCKER_IMAGE_HASH

FROM alpine@sha256:$DOCKER_IMAGE_HASH
LABEL author="tickstep"
LABEL email="tickstep@outlook.com"
LABEL version="1.0"
LABEL description="webdav file server for aliyun cloud drive"

# 时区
ENV TZ=Asia/Shanghai
# 手动下载tzdata安装包，注意要下载对应架构的： https://dl-cdn.alpinelinux.org/alpine/v3.15/community/
RUN apk add -U tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建运行目录
RUN mkdir -p /home/app
VOLUME /home/app
WORKDIR /home/app
RUN cd /home/app

# 创建配置文件目录
RUN mkdir /home/app/config

# 复制文件
COPY ./docker/app.sh /home/app/app.sh
RUN chmod +x /home/app/app.sh

COPY ./out/binary_files/aliyunpan /home/app
#RUN chmod +x /home/app/aliyunpan

EXPOSE 23077
ENV ALIYUNPAN_CONFIG_DIR=/home/app/config
ENV ALIYUNPAN_AUTH_USER=admin
ENV ALIYUNPAN_AUTH_PASSWORD=admin
ENV ALIYUNPAN_PAN_DRIVE=File
ENV ALIYUNPAN_PAN_DIR=/
ENV ALIYUNPAN_REFRESH_TOKEN=""
ENV ALIYUNPAN_TRANSFER_URL_TYPE=0
ENV ALIYUNPAN_BLOCK_SIZE=10240

ENTRYPOINT ./app.sh