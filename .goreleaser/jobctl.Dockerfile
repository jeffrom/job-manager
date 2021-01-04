FROM alpine:3.12.3
COPY jobctl /usr/local/bin/
ENTRYPOINT ["jobctl"]
