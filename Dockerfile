FROM golang:1.22.1-alpine3.19 as builder

RUN apk add bash

RUN apk add --no-cache openssh-client ansible git

RUN mkdir /workspace
WORKDIR /workspace

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build cmd/macros-backend/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /workspace .

ARG DEFAULT_PORT=3030
ENV PORT $DEFAULT_PORT

EXPOSE $PORT

CMD ["./main"]
