FROM haproxy:3.4-alpine3.24
USER root
RUN apk update && apk add --no-cache consul-template
COPY ./start.sh /start.sh
ENTRYPOINT ["/start.sh"]