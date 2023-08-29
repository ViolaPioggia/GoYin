FROM alpine:latest

WORKDIR /app

COPY  ./interaction .
COPY  ./config.yaml ./server/service/interaction/config.yaml

# 暴露端口
EXPOSE 10003

CMD ["./interaction"]