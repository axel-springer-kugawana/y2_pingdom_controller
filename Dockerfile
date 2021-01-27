FROM golang:latest AS builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/pingdom_controller
ADD . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/pingdom_controller ./
ENTRYPOINT ["./pingdom_controller"]