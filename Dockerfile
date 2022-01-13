FROM golang:1.14-alpine3.13 as bulider
ADD . /go-admin
WORKDIR /go-admin
ENV GOPROXY https://goproxy.cn,direct
RUN go get -u github.com/swaggo/swag/cmd/swag \
    && swag init \
    && CGO_ENABLED=0 GOOS=linux go build -o go-admin main.go \
    && chmod 755 go-admin

FROM alpine:3.13 as prod
WORKDIR /go-admin
RUN echo 'Asia/Shanghai' > /etc/timezone \
    && mkdir conf && mkdir -p runtime/logs && mkdir docs && mkdir ssl

COPY --from=bulider /go-admin/go-admin .
COPY --from=bulider /go-admin/conf/app.ini ./conf
COPY --from=bulider /go-admin/docs/swagger.json ./docs
COPY --from=bulider /go-admin/docs/swagger.yaml ./docs
COPY --from=bulider /go-admin/ssl/server.crt ./ssl
COPY --from=bulider /go-admin/ssl/server.key ./ssl
EXPOSE 80
EXPOSE 443
ENV LANG C.UTF-8
ENV TZ='Asia/Shanghai'
#ENTRYPOINT ["./go-admin"]
CMD ["./go-admin"]
