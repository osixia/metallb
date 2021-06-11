FROM alpine:latest

RUN apk --update-cache add ca-certificates

ADD speaker /speaker
ENTRYPOINT ["/speaker"]
