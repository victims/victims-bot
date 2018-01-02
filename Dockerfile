FROM busybox:latest
MAINTAINER Steve Milner

COPY ./victims-bot /victims-bot

EXPOSE 9999
ENTRYPOINT ["/victims-bot"]
CMD ["-secret=$VICTIMS_GITHUB_SECRET"]
