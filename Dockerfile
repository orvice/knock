FROM golang:1.9 as builder

## Create a directory and Add Code
RUN mkdir -p /go/src/github.com/orvice/knock
WORKDIR /go/src/github.com/orvice/knock
ADD .  /go/src/github.com/orvice/knock

# Download and install any required third party dependencies into the container.
RUN go-wrapper download
# RUN go-wrapper install
RUN CGO_ENABLED=0 go build

# EXPOSE 8300

# Now tell Docker what command to run when the container starts
# CMD ["go-wrapper", "run"]

FROM alpine

COPY --from=builder /go/src/github.com/orvice/knock/knock .

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
# Change TimeZone
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
# Clean APK cache
RUN rm -rf /var/cache/apk/*

ENTRYPOINT [ "./knock" ]