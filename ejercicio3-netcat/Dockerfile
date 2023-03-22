FROM alpine:latest
RUN apk update && apk add netcat-openbsd
COPY ejercicio3-netcat/netcat.sh /
RUN chmod +x netcat.sh
ENTRYPOINT ["/bin/sh"]