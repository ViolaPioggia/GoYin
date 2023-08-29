FROM alpine:latest

WORKDIR /app

COPY  ./user .
COPY  ./config.yaml ./server/service/user/config.yaml

# 暴露端口
EXPOSE 10001

CMD ["./user"]