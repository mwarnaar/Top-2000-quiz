FROM --platform=linux/amd64 golang:1.23.3 AS build
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/main

FROM --platform=linux/amd64 debian:bookworm-slim
COPY --from=build /out/main /app/
COPY --from=build /src/templates /app/templates
COPY --from=build /src/public /app/public
COPY --from=build /src/go.mod /app/
COPY --from=build /src/go.sum /app/
COPY --from=build /src/top2000.sqlite3 /app/top2000.sqlite3

WORKDIR /app
ENV GIN_MODE=release
EXPOSE 8080
CMD ["/app/main"]
