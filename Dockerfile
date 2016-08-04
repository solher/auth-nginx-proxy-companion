FROM golang:1.6-alpine

ENV APP_NAME="auth-nginx-proxy-companion"
ENV SRC_PATH="/go/src/github.com/solher/auth-nginx-proxy-companion"

RUN apk add --update git \
&& mkdir -p $SRC_PATH
COPY . $SRC_PATH
WORKDIR $SRC_PATH

RUN go get -u ./... \
&& go build -v \
&& cp $APP_NAME /usr/bin \
&& apk del git \
&& rm -rf /go/*

COPY config.default.yml /app/config.yml
COPY swagger.json /app/
VOLUME /var/lib/auth-nginx-proxy-companion

WORKDIR /

ENV CONFIG=/app/config.yml
ENV SWAGGER_LOCATION=/app/swagger.json
ENV DB_LOCATION=/var/lib/auth-nginx-proxy-companion/
ENV GC_LOCATION=/var/lib/auth-nginx-proxy-companion/
EXPOSE 3000
ENTRYPOINT ["auth-nginx-proxy-companion"]