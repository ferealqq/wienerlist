FROM golang:1.18-alpine as build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o server

FROM scratch as runner 
COPY --from=build /app/server /opt/app/server
COPY --from=build /app/VERSION /opt/app/VERSION

ENV ENV=PRD
ENV VERSION=/opt/app/VERSION
ENV PORT=3000
ENV GIN_MODE=release
# DB_DSN env is set in heroku for production. 
# DB_DSN env is set in docker-production.yml for local production build testing. 

EXPOSE 3000

CMD [ "/opt/app/server" ]