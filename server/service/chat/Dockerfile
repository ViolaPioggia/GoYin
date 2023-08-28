FROM alpine:latest

WORKDIR /

COPY ./chat .
COPY ./config.yaml ./server/service/chat/config.yaml

# 暴露端口
EXPOSE 10004

RUN ./chat

