FROM alpine

RUN apk add --no-cache curl jq coreutils
COPY graylog-configurer /

CMD ["/graylog-configurer"]
