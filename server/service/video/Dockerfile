# 第一阶段：编译 Golang 代码
FROM golang:1.20 AS build

# 设置工作目录
WORKDIR /app

# 将源代码复制到容器中
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o video ./server/service/video

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server/service/video/video .
COPY --from=build /app/server/service/video/config.yaml ./server/service/video/config.yaml

# 暴露端口
EXPOSE 10005

CMD ["./video"]