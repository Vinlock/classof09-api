FROM alpine:3.2

WORKDIR /app

RUN apk add --update ca-certificates # Certificates for SSL
COPY dist/ ./dist/

ENTRYPOINT /app/dist/app
