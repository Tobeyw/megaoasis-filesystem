FROM golang:1.17

ENV GO111MODULE="on"

ENV GOPROXY="https://goproxy.io"

ARG RT

RUN echo $RT

ENV RUNTIME=$RT

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  go build -o main ./app/main.go

EXPOSE 8888

CMD ["./main"]
