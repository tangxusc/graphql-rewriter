FROM golang:alpine AS builder

WORKDIR /build
COPY . /build
ENV GOPROXY=https://goproxy.cn
RUN go build -o rewriter main.go


FROM golang:1.20.4 as plugins

WORKDIR /build
ADD ./plugins/ .
ENV GOPROXY=https://goproxy.cn
RUN cd /build/print&&go build --buildmode=plugin -o print.so print.go
RUN cd /build/regexp_query_selection &&go build --buildmode=plugin -o regexp_query_selection.so regexp_query_selection.go
RUN cd /build/result_remove &&go build --buildmode=plugin -o result_remove.so result_remove.go

FROM alpine
WORKDIR /
COPY --from=builder /build/rewriter /rewriter
COPY --from=plugins /build/* /plugins/

CMD ["./rewriter"]