FROM golang:alpine AS builder

WORKDIR /build
COPY . /build
ENV GOPROXY=https://goproxy.cn
RUN go build -o rewriter main.go

FROM alpine
WORKDIR /
COPY --from=builder /build/rewriter /rewriter

CMD ["./rewriter"]