# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build .

ENV ENV=PRD
ENV VERSION=/app/VERSION
ENV PORT=3000
ENV GIN_MODE=release
# DB_DSN env is set in heroku for production. 
# DB_DSN env is set in docker-production.yml for local production build testing. 

EXPOSE 3000

CMD [ "/app/server" ]