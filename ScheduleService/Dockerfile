FROM golang:1.11.2 as dev


ADD . /go/src/github.com/DemoMicroservice/ScheduleService

WORKDIR /go/src/github.com/DemoMicroservice/ScheduleService

FROM alpine:3.9.2 as prod
WORKDIR /root/
COPY --from=dev /go/src/github.com/DemoMicroservice/ScheduleService .
EXPOSE 8080
ENTRYPOINT ["./schedule"]