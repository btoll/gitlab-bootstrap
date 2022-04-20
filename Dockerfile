FROM golang:alpine

# $ docker run \
#     --rm
#     --init
#     -e GITLAB_API_PRIVATE_TOKEN="$GITLAB_API_PRIVATE_TOKEN"
#     -v $(pwd)/examples:/build/examples
#     btoll/gitlab-client:latest -file examples/gitlab.json

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

