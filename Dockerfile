# 1st stage build go binary
FROM golang:1.24.1-alpine AS build-go
WORKDIR /usr/app-be

RUN set -ex && \
    apk add --no-cache gcc musl-dev

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o ./app_bin64
RUN mkdir /usr/app

RUN mv ./app_bin64 /usr/app/app_bin64
RUN mv ./configs /usr/app/configs
RUN mv ./migrations /usr/app/migrations
RUN mv ./resources /usr/app/resources

# 2nd stage
FROM alpine:latest

RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
WORKDIR /usr/app

RUN adduser -D app

COPY --from=build-go /usr/app /usr/app

RUN chown -R app:app /usr/app

USER app

CMD ["./app_bin64"]

