# 第一阶段：编译 Golang 代码
FROM golang:1.20 AS build

# 设置工作目录
WORKDIR /app

# 将源代码复制到容器中
COPY . .
RUN GOOS=linux go build -o api ./server/service/api

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server/service/api/api .
COPY --from=build /app/server/service/api/config.yaml ./server/service/api/config.yaml
COPY --from=build /app/server/service/api/IP2LOCATION-LITE-DB11.BIN ./deployment/ip_info/IP2LOCATION-LITE-DB11.BIN

# 暴露端口
EXPOSE 8080

CMD ["./api"]
