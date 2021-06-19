FROM golang:1.16.5-alpine3.13

RUN mkdir /app && mkdir /src

ADD . /src

WORKDIR /src

RUN go build -o /app/main && rm -rf /src

WORKDIR /app

CMD ["/app/main"]