FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o weight-stats-be .

FROM alpine:3.21
COPY --from=build /app/weight-stats-be /weight-stats-be
EXPOSE 8083
CMD ["/weight-stats-be"]
