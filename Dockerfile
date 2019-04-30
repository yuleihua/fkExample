FROM golang:alpine AS dev-env

WORKDIR /usr/local/go/src/airman.com/fkExample
COPY . /usr/local/go/src/airman.com/fkExample

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

RUN go get ./...

RUN go build -o dist/fkExample &&\
    cp -f dist/fkExample /usr/local/bin/ &&\
    cp -f dist/fkExample.json /usr/local/etc/ &&\

RUN ls -l && ls -l dist

CMD ["/usr/local/bin/fkExample", "-c", "/usr/local/etc/fkExample.json" ]