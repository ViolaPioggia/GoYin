FROM alpine:latest

WORKDIR /app

COPY  ./chat .
COPY  ./config.yaml ./server/service/chat/config.yaml

# 暴露端口
EXPOSE 10004

CMD ["./chat"]