FROM golang:1.17.3-buster AS builder

WORKDIR /app
ADD go.mod .
ADD go.sum .

RUN go build -v github.com/gliderlabs/ssh && \
    go build -v go.uber.org/zap

ADD . .
RUN go build -v mutagen-ssh/cmd/mutagen-ssh

FROM golang:1.17.2-buster AS builder-mutagen

WORKDIR /app
ADD build-mutagen.sh .
RUN bash ./build-mutagen.sh

FROM debian:buster

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder-mutagen  /app/mutagen-agent-v0.12.0-beta2 /usr/bin/mutagen-agent-v0.12.0-beta2
COPY --from=builder-mutagen  /app/mutagen-agent-v0.12.0-beta6 /usr/bin/mutagen-agent-v0.12.0-beta6
COPY --from=builder-mutagen  /app/mutagen-agent-v0.12.0-beta7 /usr/bin/mutagen-agent-v0.12.0-beta7
COPY --from=builder-mutagen  /app/mutagen-agent-v0.13.0-beta2 /usr/bin/mutagen-agent-v0.13.0-beta2

COPY --from=builder /app/mutagen-ssh /usr/bin/mutagen-ssh
