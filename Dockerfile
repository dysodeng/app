# builder
FROM registry.cn-chengdu.aliyuncs.com/dysodeng/golang:1.22.2 AS builder

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && echo 'Asia/Shanghai' >/etc/timezone

WORKDIR /app

ADD ./go.mod /app
ADD ./go.sum /app

RUN export GOPROXY=https://goproxy.cn && go mod download

ADD . /app

RUN CGO_ENABLED=0 go build -o app

FROM registry.cn-chengdu.aliyuncs.com/dysodeng/alpine:3.19 AS runner

COPY --from=builder /app/app /app/app
COPY --from=builder /app/var /app/var
COPY --from=builder /app/configs /app/configs
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /etc/timezone /etc/timezone

WORKDIR /app

RUN chmod -R a+w /app/var

EXPOSE 8080
EXPOSE 5000
EXPOSE 4000
EXPOSE 3000

CMD ["/app/app"]
