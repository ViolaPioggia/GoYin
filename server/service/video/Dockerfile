FROM alpine:latest

WORKDIR /app

COPY  ./video .
COPY  ./config.yaml ./server/service/video/config.yaml

# 暴露端口
EXPOSE 10005

CMD ["./video"]