FROM golang:1.17-alpine AS build-env
RUN apk --no-cache add build-base git mercurial gcc bash
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go mod download
COPY . .
RUN make build
RUN make build.migrate
RUN make build.consumer
RUN swag init -g internal/api/route.go

FROM alpine
WORKDIR /app
COPY --from=build-env /app/bin/app /app/
COPY --from=build-env /app/bin/migrate /app/
COPY --from=build-env /app/bin/scanworker /app/
COPY --from=build-env /app/assets /app/assets
ENV RULE_FILE_PATH=/app/assets/rules.json
RUN mkdir -p /tmp/git-watchdog
ENTRYPOINT ./app
