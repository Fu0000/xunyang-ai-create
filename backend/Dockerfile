FROM golang:1.24-alpine
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go
RUN mkdir -p uploads
EXPOSE 8092
CMD ["./server"]
