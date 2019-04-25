FROM golang:1.11.2 as dev


ADD . /go/src/github.com/DemoMicroservice/AuthService

WORKDIR /go/src/github.com/DemoMicroservice/AuthService

FROM alpine:3.9.2 as prod
WORKDIR /root/
COPY --from=dev /go/src/github.com/DemoMicroservice/AuthService .
EXPOSE 8080
ENTRYPOINT ["./auth"]