FROM golang:1.20.4 AS builder

WORKDIR /build
COPY . /build
ENV GOPROXY=https://goproxy.cn
RUN unset HTTPS_PROXY;unset HTTP_PROXY;unset http_proxy;unset https_proxy;go mod download
RUN unset HTTPS_PROXY;unset HTTP_PROXY;unset http_proxy;unset https_proxy;go build -o rewriter main.go


FROM golang:1.20.4 as plugins

WORKDIR /build
ADD ./plugins/ .
ENV GOPROXY=https://goproxy.cn
RUN unset HTTPS_PROXY;unset HTTP_PROXY;unset http_proxy;unset https_proxy;cd /build/print&&go build --buildmode=plugin -o print.so print.go
RUN unset HTTPS_PROXY;unset HTTP_PROXY;unset http_proxy;unset https_proxy;cd /build/regexp_query_selection &&go build --buildmode=plugin -o regexp_query_selection.so regexp_query_selection.go
RUN unset HTTPS_PROXY;unset HTTP_PROXY;unset http_proxy;unset https_proxy;cd /build/result_remove &&go build --buildmode=plugin -o result_remove.so result_remove.go

FROM alpine
WORKDIR /
COPY --from=builder /build/rewriter /rewriter
COPY --from=plugins /build/* /plugins/

CMD ["./rewriter"]