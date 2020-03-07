FROM golang:1.13.0-alpine as builder
RUN apk add git build-base
COPY ./go.mod ./go.sum /
COPY ./ /
WORKDIR /
RUN go mod download
RUN CGO_ENABLED=1  go build -o /app/detection-api

FROM alpine as runner
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*
COPY --from=builder /app /app
COPY --from=builder /resources /app/resources
COPY --from=builder /migrations /app/migrations
WORKDIR /app
ENTRYPOINT /app/detection-api
EXPOSE 3000
