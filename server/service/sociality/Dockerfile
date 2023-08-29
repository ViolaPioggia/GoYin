FROM alpine:latest

WORKDIR /app

COPY  ./sociality .
COPY  ./config.yaml ./server/service/sociality/config.yaml

# 暴露端口
EXPOSE 10002

CMD ["./sociality"]