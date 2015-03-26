FROM scratch
COPY GoWiki /
VOLUME ["/db/"]
EXPOSE 3000
ENTRYPOINT ["/GoWiki"]
