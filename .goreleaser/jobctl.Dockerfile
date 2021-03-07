FROM alpine:3.13.2
COPY jobctl /usr/local/bin/
ENTRYPOINT ["jobctl"]
