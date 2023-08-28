# 第一阶段：编译 Golang 代码
FROM golang:1.20 AS build

# 设置工作目录
WORKDIR /app

# 将源代码复制到容器中
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o interaction ./server/service/interaction

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server/service/interaction/interaction .
COPY --from=build /app/server/service/interaction/config.yaml ./server/service/interaction/config.yaml

# 暴露端口
EXPOSE 10003

CMD ["./interaction"]