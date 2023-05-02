FROM scratch

LABEL org.opencontainers.image.source https://github.com/berlingoqc/logviewer

COPY ./build/logexplorer-linux-amd64 /usr/bin/logexplorer

ENTRYPOINT ["/usr/bin/logviewer"]
