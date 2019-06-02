FROM golang:1.12.5 AS build-env

ENV GO111MODULE on

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine

RUN apk add ca-certificates
RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=build-env /go/src/app .

CMD ["./aerospike-viewer"]
