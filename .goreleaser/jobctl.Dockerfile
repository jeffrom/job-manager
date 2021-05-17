FROM scratch
COPY jobctl /usr/local/bin/
ENTRYPOINT ["jobctl"]
