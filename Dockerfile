FROM ubuntu:latest
LABEL authors="default"

ENTRYPOINT ["top", "-b"]