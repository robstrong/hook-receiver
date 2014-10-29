FROM scratch
MAINTAINER rstrong

EXPOSE 80
ENTRYPOINT ["/hook-receiver"]
ADD bin/hook-receiver-linux-amd64 hook-receiver
