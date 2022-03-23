FROM golang:alpine

RUN adduser \
    --uid 1000 \
    --no-create-home \
    --disabled-password \
    --shell /bin/sh \
    noroot

RUN mkdir /build
COPY *.go go.* /build/
WORKDIR /build
RUN go build -o gitlab-client .

USER noroot

ENTRYPOINT ["./gitlab-client"]

