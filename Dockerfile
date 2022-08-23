# Building the binary of the App
FROM golang:1.19 AS build

WORKDIR /go/src/neptune

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM alpine:latest

WORKDIR /app

RUN mkdir ./static
COPY ./ ./

COPY --from=build /go/src/neptune/app .

EXPOSE 443

CMD ["./app"]