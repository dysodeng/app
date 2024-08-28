# builder
FROM registry.huaxisy.com/library/golang:1.22.2 AS Builder

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && echo 'Asia/Shanghai' >/etc/timezone

WORKDIR /app

ADD ./go.mod /app
ADD ./go.sum /app

RUN export GOPROXY=https://goproxy.cn && go mod download

ADD . /app

RUN CGO_ENABLED=0 go build -o app

FROM registry.huaxisy.com/library/alpine:3.19 AS Runner

COPY --from=Builder /app/app /app/app
COPY --from=Builder /app/var /app/var
COPY --from=Builder /app/configs /app/configs
COPY --from=Builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=Builder /etc/localtime /etc/localtime
COPY --from=Builder /etc/timezone /etc/timezone

WORKDIR /app

RUN chmod -R a+w /app/var

EXPOSE 8080
EXPOSE 5000

CMD /app/app
