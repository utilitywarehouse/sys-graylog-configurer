FROM alpine

COPY graylog-configurer /

RUN apk add --no-cache curl jq coreutils

CMD ["/graylog-configurer"]
