FROM golang:alpine3.18 AS builder

WORKDIR /app

ENV GOPROXY https://goproxy.cn,direct

# 如果 go.mod 没有变化就会使用缓存
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o server  ./

FROM alpine

WORKDIR /app

ENV TZ Asia/Shanghai

COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml /app/conf
COPY --from=builder /app/deploy /app/deploy

VOLUME C:/data/conf:/conf

EXPOSE 8001
ENTRYPOINT ["./server"]
CMD ["-config", "/conf"]