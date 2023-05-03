FROM scratch

LABEL org.opencontainers.image.source https://github.com/berlingoqc/logviewer

COPY ./build/logviewer-linux-amd64 /usr/bin/logviewer

ENTRYPOINT ["/usr/bin/logviewer"]
