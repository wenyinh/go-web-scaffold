# 第一阶段：构建
FROM golang:latest AS build

# 设置工作目录
WORKDIR /app

# 复制go.mod go.sum
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=build /app/main .

# 复制配置文件
COPY --from=build /app/config.yaml .

COPY --from=build /app/wait-for-it.sh .

RUN chmod +x wait-for-it.sh

EXPOSE 8088
