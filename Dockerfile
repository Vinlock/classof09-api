FROM alpine:3.2

WORKDIR /app

RUN apk add --update ca-certificates # Certificates for SSL
COPY dist/ ./dist/
COPY templates/ ./templates/

ENTRYPOINT /app/dist/discordbot
