ARG GIT_COMMIT
ARG VERSION
ARG PROJECT

FROM golang:latest AS build

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

ARG VERSION
ENV VERSION=$VERSION

ARG PROJECT
ENV PROJECT=$PROJECT

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '\
    -w -extldflags "-static"\
    -X ${PROJECT}/app/config.Version=${VERSION} -X ${PROJECT}/app/config.Commit=${GIT_COMMIT}\
' -o ./run ./cmd/shorturl


FROM scratch

WORKDIR /app

COPY --from=build /app/run .
COPY --from=build /app/web /app/web
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow

EXPOSE 8080

CMD ["./run"]
