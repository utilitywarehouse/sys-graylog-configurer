FROM alpine

COPY run.sh /

RUN apk add --no-cache curl jq coreutils \
    && chmod +x /run.sh

CMD ["/run.sh"]
