# Build stage
FROM golang:1.26.3-alpine3.23 AS builder 
WORKDIR /app
COPY . .
RUN go build -o main main.go
# RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1 \
#     && mv $(go env GOPATH)/bin/migrate /app/migrate

# Run stage
FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]