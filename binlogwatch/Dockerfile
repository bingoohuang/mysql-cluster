FROM golang:1.20.3 as builder
WORKDIR /build
COPY ./binlogwatch /build
ENV GOPROXY=https://goproxy.cn
RUN cd /build && go build

# FROM mysql:5.7.42
FROM mysql:8.0.33
ENV MYSQL_ROOT_PASSWORD=root
COPY --from=builder /build/binlogwatch /
COPY schema.sql /schema.sql

# 构建：docker build -t bingoohuang/binglogwatch .
# docker buildx build --platform linux/amd64,linux/arm64/v8 -t bingoohuang/binglogwatch . --push 
# 运行：docker run -it bingoohuang/binglogwatch

