FROM centos:latest
MAINTAINER Steve Milner

COPY ./victims-bot /victims-bot

EXPOSE 9999
ENTRYPOINT ["/victims-bot"]
