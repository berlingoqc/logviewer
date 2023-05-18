FROM golang

WORKDIR /app

COPY main.go .
COPY go.* ./
ADD pkg ./pkg
ADD cmd ./cmd
RUN go build -o logviewer

FROM busybox:glibc

LABEL org.opencontainers.image.source https://github.com/berlingoqc/logviewer

COPY --from=0 /app/logviewer /logviewer
ENTRYPOINT ["/logviewer"]