FROM scratch

COPY ./build/logexplorer-linux-amd64 /usr/bin/logexplorer

ENTRYPOINT ["/usr/bin/logexplorer"]
