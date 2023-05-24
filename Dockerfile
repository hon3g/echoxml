# syntax=docker/dockerfile:1

FROM golang@sha256:3fccedea46315261e4b6205bcffe91ece1e2aea60c23aab0f033f35461849b42

WORKDIR /app

COPY go.mod ./

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /echoxml

EXPOSE 8080

CMD ["/echoxml"]